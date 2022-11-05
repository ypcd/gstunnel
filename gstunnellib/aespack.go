package gstunnellib

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	. "github.com/ypcd/gstunnel/v6/gstunnellib/gspackoper"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

var key_hash_salt = []byte{
	208, 48, 220, 100, 85, 10, 162, 39, 88, 13, 220, 88, 121, 70, 220, 120, 40, 139, 124, 162, 85, 25, 242, 90, 8, 112, 69, 93, 17, 54, 176, 22, 165, 191, 226, 231, 35, 208, 99, 37, 78, 228, 138, 214, 46, 84, 243, 195, 242, 13, 151, 5, 123, 28, 71, 101, 41, 20, 194, 228, 30, 220, 72, 216, 24, 255, 106, 140, 64, 44, 132, 249, 146, 79, 52, 138, 3, 15, 186, 3, 148, 148, 155, 137, 61, 58, 68, 199, 102, 104, 69, 186, 210, 231, 196, 246, 181, 135, 175, 208, 109, 176, 196, 252, 159, 36, 249, 90, 125, 24, 136, 172, 207, 69, 199, 206, 12, 21, 177, 10, 121, 98, 43, 39, 163, 107, 246, 226,
}

func Find0(v1 []byte) (int, bool) {
	i := bytes.IndexByte(v1, 0)
	if i == -1 {
		return -1, false
	} else {
		return i, true
	}
}

type GsPack interface {
	Packing([]byte) []byte
	Unpack([]byte) ([]byte, error)

	ChangeCryKey() []byte
	IsTheVersionConsistent() []byte
}

type aespack struct {
	a gsaes
}

func (ap *aespack) Packing(data []byte) []byte {

	jdata := JsonPacking(data)
	crydata := ap.a.encrypter(jdata)
	/*
		enlen := base64.StdEncoding.EncodedLen(len(crydata)) + 1
		endata := make([]byte, enlen)
		base64.StdEncoding.Encode(endata, crydata)
		endata[enlen-1] = 0
	*/
	IsOk_PackDataLen_panic(crydata)

	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata
}
func (ap *aespack) Unpack_old(data []byte) ([]byte, error) {
	var jdata []byte
	if Deep_debug {
		logger.Println("aespack unpack data hash:", gshash.GetSha256Hex(data))

		ix := bytes.IndexByte(data, 0)
		if ix != -1 {
			if ix == (len(data) - 1) {
				data = data[:ix]
			} else {
				checkError_panic(errors.New("Data is error."))
			}
		}

		dedata := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		re, err := base64.StdEncoding.Decode(dedata, data)
		checkError_panic(err)
		dedata = dedata[:re]
		logger.Println("aespack unpack base64 decode hash:", gshash.GetSha256Hex(dedata))

		jdata := ap.a.decrypter(dedata)
		logger.Println("aespack unpack aes decrypt hash:", gshash.GetSha256Hex(jdata))
	} else {

		ix := bytes.IndexByte(data, 0)
		if ix != -1 {
			if ix == (len(data) - 1) {
				data = data[:ix]
			} else {
				checkError_panic(errors.New("Data is error."))
			}
		}

		dedata := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		re, err := base64.StdEncoding.Decode(dedata, data)
		checkError_panic(err)
		dedata = dedata[:re]

		jdata = ap.a.decrypter(dedata)

	}
	var pdata []byte
	var err error

	po1 := UnPack_Oper(jdata)

	err = po1.IsOk()
	if err != nil {
		return nil, err
	}

	if po1.IsChangeCryKey() {
		err = ap.setKey(GetEncryKey(jdata))
		logger.Println("gstunnel is ChangeCryKey.")
		pdata = nil
		return nil, err
	}

	if po1.IsPOVersion() {
		if po1.OperType == POVersion {
			if string(po1.OperData) == Version {
				logger.Println("gstunnel POVersion is ok.")
				pdata = nil
				err = nil
			} else {
				pdata = nil
				err = errors.New("Version is error.  " + "error version: " + string(po1.OperData))
			}
		}
		return pdata, err
	}
	if po1.IsPOGen() {
		return po1.Data, err
	}

	return pdata, err
}

func (ap *aespack) Unpack(data []byte) ([]byte, error) {
	var jdata []byte
	if Deep_debug {
		logger.Println("aespack unpack data hash:", gshash.GetSha256Hex(data))

		jdata := ap.a.decrypter(data[2:])
		logger.Println("aespack unpack aes decrypt hash:", gshash.GetSha256Hex(jdata))
	} else {
		jdata = ap.a.decrypter(data[2:])
	}
	var pdata []byte
	var err error

	po1 := UnPack_Oper(jdata)

	err = po1.IsOk()
	if err != nil {
		return nil, err
	}

	if po1.IsChangeCryKey() {
		err = ap.setKey(GetEncryKey(jdata))
		logger.Println("gstunnel is ChangeCryKey.")
		pdata = nil
		return nil, err
	}

	if po1.IsPOVersion() {
		if po1.OperType == POVersion {
			if string(po1.OperData) == Version {
				logger.Println("gstunnel POVersion is ok.")
				pdata = nil
				err = nil
			} else {
				pdata = nil
				err = errors.New("Version is error.  " + "error version: " + string(po1.OperData))
			}
		}
		return pdata, err
	}
	if po1.IsPOGen() {
		return po1.Data, err
	}

	return pdata, err
}

func (ap *aespack) setKey(key []byte) error {
	ap.a = createAes(key)
	return nil
}

func (ap *aespack) ChangeCryKey() []byte {

	jdata := JsonPacking_OperChangeKey()
	crydata := ap.a.encrypter(jdata)
	//edata := base64.StdEncoding.EncodeToString(crydata)

	ap.setKey(GetEncryKey(jdata))

	IsOk_PackDataLen_panic(crydata)
	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata
}

func (ap *aespack) IsTheVersionConsistent() []byte {

	jdata := JsonPacking_OperVersion()
	crydata := ap.a.encrypter(jdata)
	//edata := base64.StdEncoding.EncodeToString(crydata)

	IsOk_PackDataLen_panic(crydata)
	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata
}

func newAesPack(key string) *aespack {
	if len(key) != 32 {
		checkError_exit(errors.New("Error: The key is not 32 bytes."))
	}
	keyhash := GetSha256_32bytes(append([]byte(key), key_hash_salt...))
	ap1 := aespack{}
	ap1.a = createAes(keyhash[:])

	return &ap1
}

func NewGsPack(key string) GsPack {
	return newAesPack(key)
}

func GetRDKeyString32() string {
	return gsrand.GetrandStringPlus(32)
}

func IsOk_PackDataLen_panic(data []byte) {
	if len(data) >= (1 << 16) {
		panic("The pack data len >= 65536. It is Error.")
	}
}
