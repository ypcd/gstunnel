/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
*/
package gstunnellib

import (
    "encoding/json"
    "errors"
    "crypto/aes"
    "crypto/cipher"
    "fmt"
    //"sync"
)

var p func(...interface{}) (int, error)

func init() {
    p = fmt.Println
    p = Nullprint
}

func Nullprint(v ...interface{}) (int, error) {return 1, errors.New("")}
func Nullprintf(format string, a ...interface{}) (n int, err error){return 1, errors.New("")}

type v struct {
    V1 int
}

type ls struct {
    Vls []v
}

type Pack struct {
    Data []byte
}

func find0(v1 []byte) (int, bool) {
    for i := 0; i<len(v1); i++ {
        if v1[i]==0 {
            return i, true
        }
    }
    return -1, false
}

func jsonUnpack(data []byte) ([]byte) {
    var msg Pack
    ix, re := find0(data)
    if !re {
        return []byte{}
    }
    data = data[:ix]

    json.Unmarshal( data, &msg)
    return msg.Data[:]
}

func jsonPacking(data []byte) ([]byte) {
    pd := Pack{Data:data}
    re, _ := json.Marshal( pd)
    
    return append(re, 0)
}



func encrypter(data []byte) ([]byte) {
    commonIV := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
    key := "5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ"

    cpr, _ := aes.NewCipher([]byte(key))
    cfbE := cipher.NewCFBEncrypter(cpr, commonIV)

    dst := make([]byte, len(data))

    cfbE.XORKeyStream(dst, data)
    return dst
}

func decrypter(data []byte) ([]byte) {
    commonIV := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
    key1 := "5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ"

    key := key1
    cpr, _ := aes.NewCipher([]byte(key))
    cfbD := cipher.NewCFBDecrypter(cpr, commonIV)

    dst := make([]byte, len(data))

    cfbD.XORKeyStream(dst, data)
    return dst
}

func Packing(data []byte) ([]byte) {
    data = encrypter(data)
    return jsonPacking(data)
}
func Unpack(data []byte) ([]byte) {
    data = jsonUnpack(data)
    return decrypter(data)
}


func CreatePack() (Pack) {
    return Pack{}
}

type Aes struct {
    //locken, lockde sync.Mutex
    //commonIV []byte
    //key string
    cpr cipher.Block
    cfbE, cfbD cipher.Stream
}

func CreateAes(key string) Aes {
    var v1 Aes
    commonIV := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

    v1.cpr, _ = aes.NewCipher([]byte(key))
    v1.cfbE = cipher.NewCFBEncrypter(v1.cpr, commonIV)
    v1.cfbD = cipher.NewCFBDecrypter(v1.cpr, commonIV)
    return v1
}


func (a *Aes) encrypter(data []byte) ([]byte) {   
    dst := make([]byte, len(data))
    a.cfbE.XORKeyStream(dst, data)
    return dst
}

func (a *Aes) decrypter(data []byte) ([]byte) {   
    dst := make([]byte, len(data))
    a.cfbD.XORKeyStream(dst, data)
    return dst
}

type Aespack struct {
    a Aes
}

func (ap *Aespack) Packing(data []byte) ([]byte) {
    data = ap.a.encrypter(data)
    return jsonPacking(data)
}
func (ap *Aespack) Unpack(data []byte) ([]byte) {
    data = jsonUnpack(data)
    return ap.a.decrypter(data)
}


func CreateAesPack(key string) (Aespack) {
    return Aespack{a:CreateAes(key)}
}

