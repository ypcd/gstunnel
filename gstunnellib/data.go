/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package gstunnellib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"

	"gstunnel/gstunnellib/gsbase"
	. "gstunnel/gstunnellib/gspackoper"
)

const Version string = gsbase.Version

var p func(...interface{}) (int, error)

var begin_Dinfo int = 0

var debug_tag bool

var commonIV = []byte{171, 158, 1, 73, 31, 98, 64, 85, 209, 217, 131, 150, 104, 219, 33, 220}

var logger *log.Logger

const Info_protobuf bool = true

func Nullprint(v ...interface{}) (int, error)                       { return 1, nil }
func Nullprintf(format string, a ...interface{}) (n int, err error) { return 1, nil }

func init() {
	debug_tag = false
	p = Nullprint

	logger = CreateFileLogger("gstunnellib.data.log")
	//debug_tag = true
	//p = fmt.Println
}

func Find0(v1 []byte) (int, bool) {
	i := bytes.IndexByte(v1, 0)
	if i == -1 {
		return -1, false
	} else {
		return i, true
	}
}

type gsaes struct {
	cpr        cipher.Block
	cfbE, cfbD cipher.Stream
}

func createAes(key []byte) gsaes {
	var v1 gsaes

	v1.cpr, _ = aes.NewCipher(key)
	v1.cfbE = cipher.NewCFBEncrypter(v1.cpr, commonIV)
	v1.cfbD = cipher.NewCFBDecrypter(v1.cpr, commonIV)
	return v1
}

func (a *gsaes) encrypter(data []byte) []byte {
	dst := make([]byte, len(data))
	a.cfbE.XORKeyStream(dst, data)
	return dst
}

func (a *gsaes) decrypter(data []byte) []byte {
	dst := make([]byte, len(data))
	a.cfbD.XORKeyStream(dst, data)
	return dst
}

type GsPack interface {
	Packing(data []byte) []byte
	Unpack(data []byte) ([]byte, error)

	ChangeCryKey() []byte
	IsTheVersionConsistent() []byte
}

type aespack struct {
	a gsaes
}

func (ap *aespack) Packing(data []byte) []byte {
	jdata := JsonPacking(data)
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)
	return append([]byte(edata), 0)
}
func (ap *aespack) Unpack(data []byte) ([]byte, error) {
	d1, _ := base64.StdEncoding.DecodeString(string(data))
	jdata := ap.a.decrypter(d1)
	var pdata []byte
	var err error

	po1 := UnPack_Oper(jdata)

	err = po1.IsOk()
	if err != nil {
		return []byte{}, err
	}

	if po1.IsChangeCryKey() {
		err = ap.setKey(GetKey(jdata))
		pdata = []byte{}
		return pdata, err
	}

	if po1.IsPOVersion() {
		if po1.OperType == POVersion {
			if string(po1.OperData) == Version {
				pdata = []byte{}
				err = nil
			} else {
				pdata = []byte{}
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

	ap.setKey(GetKey(jdata))

	return append([]byte(edata), 0)
}

func (ap *aespack) IsTheVersionConsistent() []byte {

	jdata := JsonPacking_OperVersion()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	return append([]byte(edata), 0)
}

func createAesPack(key string) *aespack {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		logger.Fatalln("Error: The key is not 16, 24, or 32 bytes.")
	}
	ap1 := aespack{}
	ap1.a = createAes([]byte(key))

	return &ap1
}

func NewGsPack(key string) GsPack {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		logger.Fatalln("Error: The key is not 16, 24, or 32 bytes.")
	}
	ap1 := aespack{}
	ap1.a = createAes([]byte(key))

	return &ap1
}
