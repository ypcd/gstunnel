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
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"

	//"sync"
	"time"
)

var p func(...interface{}) (int, error)

var begin_Dinfo int = 0
var rand_s rand.Source = rand.NewSource(time.Now().Unix())

var rd_s rand.Source = rand.NewSource(time.Now().Unix())

var debug_tag bool = false

func init() {
	p = fmt.Println
	p = Nullprint
}

func Nullprint(v ...interface{}) (int, error)                       { return 1, errors.New("") }
func Nullprintf(format string, a ...interface{}) (n int, err error) { return 1, errors.New("") }

type v struct {
	V1 int
}

type ls struct {
	Vls []v
}

type Pack struct {
	Data []byte
}

//Data, Rand
type PackRand struct {
	Pack
	Rand int64
}

//Data, Rand, HashHex
type PackRandHash struct {
	PackRand
	HashHex string
}

const (
	GenOper      int32 = 1
	ChangeCryKey int32 = 2
)

////Data, Rand, HashHex, opertype, operData
type PackOper struct {
	OperType int32
	OperData []byte
	Data     []byte
	Rand     int64
	HashHex  string
}

func GetRDF64() float64 {

	return float64(rd_s.Int63()+1) / math.Pow(2, 64)
}

func GetRDInt8() int8 {
	return int8(GetRDF64() * 256)
}

func GetRDInt16() int16 {
	return int16(GetRDF64() * math.Pow(2, 16))
}

func GetRDInt32() int32 {
	return int32(GetRDF64() * math.Pow(2, 32))
}

func GetRDInt64() int64 {
	return int64(GetRDF64() * math.Pow(2, 64))
}

func GetRDBytes(byteLen int) []byte {
	data := make([]byte, byteLen)
	for i := 0; i < byteLen; i++ {
		data[i] = byte(GetRDInt8())
	}
	return data
}

func IsChangeCryKey(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == ChangeCryKey)
}

func GetKey(Data []byte) []byte {
	return (UnPack_Oper(Data).OperData)
}

func Int32ToBytes(i32 int32) []byte {
	buf := make([]byte, 0)
	bbuf := bytes.NewBuffer(buf)
	err := binary.Write(bbuf, binary.LittleEndian, i32)
	if err != nil {
		panic(err)
	}
	return bbuf.Bytes()
}

func Int64ToBytes(data int64) []byte {
	buf := make([]byte, 0)
	bbuf := bytes.NewBuffer(buf)
	err := binary.Write(bbuf, binary.LittleEndian, data)
	if err != nil {
		panic(err)
	}
	return bbuf.Bytes()
}

func CreatePackOperGen(data []byte) PackOper {

	rd := rand_s.Int63()

	pd := PackOper{
		OperType: GenOper,
		//OperData:	[]byte,
		Data: data, Rand: rd,
		HashHex: GetSha256Hex(append(data, Int64ToBytes(rd)...)),
	}
	return pd
}

func CreatePackOperChangeKey() PackOper {

	rd := rand_s.Int63()
	key := GetRDBytes(32)

	pd := PackOper{
		OperType: ChangeCryKey,
		OperData: []byte(key),
		Rand:     rd,
		HashHex:  GetSha256Hex(append(Int32ToBytes(GenOper), append([]byte(key), Int64ToBytes(rd)...)...)),
	}
	return pd
}

func jsonPacking_OperChangeKey() []byte {

	pd := CreatePackOperChangeKey()

	re, _ := json.Marshal(pd)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
}

func jsonPacking_OperGen(data []byte) []byte {

	pd := CreatePackOperGen(data)

	re, _ := json.Marshal(pd)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
}

func UnPack_Oper(data []byte) PackOper {
	var msg PackOper
	ix, re := Find0(data)
	if !re {
		return PackOper{}
	}
	data2 := data[:ix]

	json.Unmarshal(data2, &msg)
	if debug_tag {
		p(msg)
		p("json: ", string(data2))
	}
	return msg
}

func jsonUnPack_OperGen(data []byte) []byte {

	p1 := UnPack_Oper(data)

	return p1.Data[:]
}

func GetSha256Hex(data []byte) string {
	s1 := sha256.Sum256(data)
	return hex.EncodeToString(
		s1[:],
	)
}

func find0_1(v1 []byte) (int, bool) {
	for i := 0; i < len(v1); i++ {
		if v1[i] == 0 {
			return i, true
		}
	}
	return -1, false
}

func Find0(v1 []byte) (int, bool) {
	i := bytes.IndexByte(v1, 0)
	if i == -1 {
		return -1, false
	} else {
		return i, true
	}
}

func jsonPacking(data []byte) []byte {
	return jsonPacking_OperGen(data)
}

func jsonUnpack(data []byte) []byte {
	if begin_Dinfo == 0 {
		fmt.Println("Pack_Rand use.")
		begin_Dinfo = 1
	}
	return jsonUnPack_OperGen(data)
}

func jsonUnpack_1(data []byte) []byte {
	var msg Pack
	ix, re := Find0(data)
	if !re {
		return []byte{}
	}
	data = data[:ix]

	json.Unmarshal(data, &msg)
	return msg.Data[:]
}

func jsonPacking_1(data []byte) []byte {
	pd := Pack{Data: data}
	re, _ := json.Marshal(pd)

	return append(re, 0)
}

func jsonUnPack_Rand(data []byte) []byte {
	var msg PackRand
	ix, re := Find0(data)
	if !re {
		return []byte{}
	}
	data = data[:ix]

	json.Unmarshal(data, &msg)
	return msg.Data[:]
}

func jsonPacking_Rand(data []byte) []byte {
	pd := PackRand{Pack: Pack{Data: data}, Rand: rand_s.Int63()}
	re, _ := json.Marshal(pd)

	return append(re, 0)
}

func jsonPacking_RandHash(data []byte) []byte {

	rd := rand_s.Int63()
	buf := make([]byte, 100)
	blen := binary.PutVarint(buf, rd)
	buf = buf[:blen]

	pd := PackRandHash{
		PackRand{Pack: Pack{Data: data}, Rand: rd},
		GetSha256Hex(append(data, buf...)),
	}

	re, _ := json.Marshal(pd)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
}

func jsonUnPack_RandHash(data []byte) []byte {
	var msg PackRandHash
	ix, re := Find0(data)
	if !re {
		return []byte{}
	}
	data = data[:ix]

	json.Unmarshal(data, &msg)
	if debug_tag {
		p(msg)
		p("json: ", string(data))
	}
	return msg.Data[:]
}

func encrypter(data []byte) []byte {
	commonIV := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	key := "5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ"

	cpr, _ := aes.NewCipher([]byte(key))
	cfbE := cipher.NewCFBEncrypter(cpr, commonIV)

	dst := make([]byte, len(data))

	cfbE.XORKeyStream(dst, data)
	return dst
}

func decrypter(data []byte) []byte {
	commonIV := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	key1 := "5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ"

	key := key1
	cpr, _ := aes.NewCipher([]byte(key))
	cfbD := cipher.NewCFBDecrypter(cpr, commonIV)

	dst := make([]byte, len(data))

	cfbD.XORKeyStream(dst, data)
	return dst
}

func Packing(data []byte) []byte {
	data = encrypter(data)
	return jsonPacking(data)
}
func Unpack(data []byte) []byte {
	data = jsonUnpack(data)
	return decrypter(data)
}

func CreatePack() Pack {
	return Pack{}
}

type Aes struct {
	//locken, lockde sync.Mutex
	//commonIV []byte
	//key string
	cpr        cipher.Block
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

func (a *Aes) encrypter(data []byte) []byte {
	dst := make([]byte, len(data))
	a.cfbE.XORKeyStream(dst, data)
	return dst
}

func (a *Aes) decrypter(data []byte) []byte {
	dst := make([]byte, len(data))
	a.cfbD.XORKeyStream(dst, data)
	return dst
}

type Aespack struct {
	a Aes
}

func (ap *Aespack) Packing(data []byte) []byte {
	jdata := jsonPacking(data)
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)
	return append([]byte(edata), 0)
}
func (ap *Aespack) Unpack(data []byte) []byte {

	d1, _ := base64.StdEncoding.DecodeString(string(data))
	jdata := ap.a.decrypter(d1)
	pdata := []byte{}
	if IsChangeCryKey(jdata) {
		ap.changeCryKey(GetKey(jdata))
		pdata = []byte{}
	} else {
		pdata = jsonUnpack(jdata)
	}
	return pdata
}

func (ap *Aespack) changeCryKey(key []byte) error {
	ap.a = CreateAes(string(key))
	return nil
}

func (ap *Aespack) ChangeCryKey() []byte {

	jdata := jsonPacking_OperChangeKey()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	ap.changeCryKey(GetKey(jdata))

	return append([]byte(edata), 0)
}

func CreateAesPack(key string) Aespack {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Error: The key is not 16, 24, or 32 bytes.")
		//panic(errors.New("Error: The key is not 16, 24, or 32 bytes."))
		os.Exit(10)
	}
	return Aespack{a: CreateAes(key)}
}
