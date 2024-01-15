package gstunnellib

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gspackoper"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

/*
var key_hash_salt = []byte{
	208, 48, 220, 100, 85, 10, 162, 39, 88, 13, 220, 88, 121, 70, 220, 120, 40, 139, 124, 162, 85, 25, 242, 90, 8, 112, 69, 93, 17, 54, 176, 22, 165, 191, 226, 231, 35, 208, 99, 37, 78, 228, 138, 214, 46, 84, 243, 195, 242, 13, 151, 5, 123, 28, 71, 101, 41, 20, 194, 228, 30, 220, 72, 216, 24, 255, 106, 140, 64, 44, 132, 249, 146, 79, 52, 138, 3, 15, 186, 3, 148, 148, 155, 137, 61, 58, 68, 199, 102, 104, 69, 186, 210, 231, 196, 246, 181, 135, 175, 208, 109, 176, 196, 252, 159, 36, 249, 90, 125, 24, 136, 172, 207, 69, 199, 206, 12, 21, 177, 10, 121, 98, 43, 39, 163, 107, 246, 226,
}
*/

func Find0(v1 []byte) (int, bool) {
	i := bytes.IndexByte(v1, 0)
	if i == -1 {
		return -1, false
	} else {
		return i, true
	}
}

type IGSPack interface {
	Packing([]byte) []byte
	Unpack([]byte) ([]byte, error)

	ChangeCryKey() []byte
	PackVersion() []byte
}

type aespack struct {
	a gsaes
}

//type aespackPlusRSA struct

func (ap *aespack) Packing(data []byte) []byte {

	jdata := gspackoper.NewPackPOGen(data)
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

func (ap *aespack) Unpack(data []byte) ([]byte, error) {
	var jdata []byte
	if G_Deep_debug {
		G_logger.Println("aespack unpack data hash:", gshash.GetSha256Hex(data))

		jdata := ap.a.decrypter(data[2:])
		G_logger.Println("aespack unpack aes decrypt hash:", gshash.GetSha256Hex(jdata))
	} else {
		jdata = ap.a.decrypter(data[2:])
	}
	//var pdata []byte
	var err error

	po1 := gspackoper.UnPack_Oper(jdata)

	err = po1.IsOk()
	if err != nil {
		return nil, err
	}

	if po1.IsPOGen() {
		return po1.Data, err
	}
	if po1.IsChangeCryKey() {
		err = ap.setKey(po1.GetEncryKey())
		G_logger.Println("gstunnel is ChangeCryKey.")
		//pdata = nil
		return nil, err
	}

	if po1.IsPOVersion() {
		if string(po1.OperData) == G_Version {
			G_logger.Println("gstunnel POVersion is ok.")
			//pdata = nil
			//err = nil
			return nil, nil
		} else {
			//pdata = nil
			err = errors.New("G_Version is error. " + "error version: " + string(po1.OperData))
			panic(err)
		}
		//return pdata, err
	}

	panic("error")
	//return pdata, err
}

func (ap *aespack) setKey(key []byte) error {
	ap.a = createAes(key)
	return nil
}

func (ap *aespack) ChangeCryKey() []byte {

	jdata, key := gspackoper.JsonPacking_OperChangeKey()
	crydata := ap.a.encrypter(jdata)
	//edata := base64.StdEncoding.EncodeToString(crydata)

	ap.setKey(key)

	IsOk_PackDataLen_panic(crydata)

	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata
}

func (ap *aespack) PackVersion() []byte {

	jdata := gspackoper.NewPackPOVersion()
	crydata := ap.a.encrypter(jdata)
	//edata := base64.StdEncoding.EncodeToString(crydata)

	IsOk_PackDataLen_panic(crydata)

	packdata := []byte{}
	packdata = binary.BigEndian.AppendUint16(packdata, uint16(len(crydata)))
	packdata = append(packdata, crydata...)

	return packdata
}

func newAesPack(key string) *aespack {
	if len([]byte(key)) != gsbase.G_AesKeyLen {
		checkError_exit(
			fmt.Errorf("error: The key is not %d bytes", gsbase.G_AesKeyLen))
	}

	return &aespack{createAes([]byte(key))}
}

func NewGsPack(key string) IGSPack {
	return newAesPack(key)
}

/*
	func GetRDKeyString32() string {
		return gsrand.GetrandStringPlus(32)
	}
*/
func GetRDKeyString96() string {
	return gsrand.GetrandStringPlus(gsbase.G_AesKeyLen)
}

/*
func GenBase64Key(key string) string {
	return base64.StdEncoding.EncodeToString([]byte(key))
}
*/

func IsOk_PackDataLen_panic(data []byte) {
	if len(data) >= (1 << 16) {
		panic("The pack data len >= 65536. It is Error.")
	}
}
