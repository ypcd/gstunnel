package gstunnellib

import (
	"bytes"
	"encoding/base64"
	"errors"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	. "github.com/ypcd/gstunnel/v6/gstunnellib/gspackoper"
)

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

	enlen := base64.StdEncoding.EncodedLen(len(crydata)) + 1
	endata := make([]byte, enlen)
	base64.StdEncoding.Encode(endata, crydata)
	endata[enlen-1] = 0

	return endata
}
func (ap *aespack) Unpack(data []byte) ([]byte, error) {
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

func (ap *aespack) setKey(key []byte) error {
	ap.a = createAes(key)
	return nil
}

func (ap *aespack) ChangeCryKey() []byte {

	jdata := JsonPacking_OperChangeKey()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	ap.setKey(GetEncryKey(jdata))

	return append([]byte(edata), 0)
}

func (ap *aespack) IsTheVersionConsistent() []byte {

	jdata := JsonPacking_OperVersion()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	return append([]byte(edata), 0)
}

func newAesPack(key string) *aespack {
	if len(key) != 32 {
		checkError_exit(errors.New("Error: The key is not 32 bytes."))
	}
	ap1 := aespack{}
	ap1.a = createAes([]byte(key))

	return &ap1
}

func NewGsPack(key string) GsPack {
	ap1 := newAesPack(key)
	return ap1
}
