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
	randc "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"sync"
)

const Version string = "V3.2"

var p func(...interface{}) (int, error)

var begin_Dinfo int = 0

var debug_tag bool

var commonIV = []byte{171, 158, 1, 73, 31, 98, 64, 85, 209, 217, 131, 150, 104, 219, 33, 220}

var Logger *log.Logger

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

func Nullprint(v ...interface{}) (int, error)                       { return 1, nil }
func Nullprintf(format string, a ...interface{}) (n int, err error) { return 1, nil }

func init() {
	debug_tag = false
	p = Nullprint

	//debug_tag = true
	//p = fmt.Println
}
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		Logger.Fatalln(err.Error())
		//os.Exit(-11)
	}
}

func CreateFileLogger(FileName string) *log.Logger {
	lf, err := os.Create(FileName)
	CheckError(err)
	Logger = log.New(lf, "Error: ", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return Logger
}

type GsConfig struct {
	Listen string
	Server string
	Key    string
	Debug  bool
}

func CreateGsconfig(confn string) *GsConfig {
	f, err := os.Open(confn)
	CheckError(err)

	defer func() {
		f.Close()
	}()

	buf, err := ioutil.ReadAll(f)
	CheckError(err)

	//fmt.Println(string(buf))
	var gsconfig GsConfig
	json.Unmarshal(buf, &gsconfig)
	return &gsconfig
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
	POBegin      int32 = 0
	GenOper      int32 = 1
	ChangeCryKey int32 = 2
	POVersion    int32 = 3
	POEnd        int32 = 4
)

//Data, Rand, HashHex, opertype, operData
type packOper struct {
	OperType int32
	OperData []byte
	Data     []byte
	Rand     int64
	HashHex  string
}

func (po *packOper) GetSha256_old() string {
	return GetSha256Hex(append(Int32ToBytes(po.OperType), append(po.OperData, append(po.Data, Int64ToBytes(po.Rand)...)...)...))
}

func (po *packOper) GetSha256() string {
	return po.GetSha256_buf()
}

//GetSha256_buf速度比GetSha256_old快15%-20%。
func (po *packOper) GetSha256_buf() string {

	var buf bytes.Buffer
	buf.Write(Int32ToBytes(po.OperType))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(Int64ToBytes(po.Rand))

	return GetSha256Hex(buf.Bytes())

}

//为了更好的安全性，默认情况下没有使用GetSha256_pool函数。
//GetSha256_pool的速度比GetSha256_buf的速度快10%-15%。

func (po *packOper) GetSha256_pool() string {

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()
	//fmt.Println("CAP LEN:", buf.Cap(), buf.Len())

	buf.Write(Int32ToBytes(po.OperType))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(Int64ToBytes(po.Rand))

	h1 := GetSha256Hex(buf.Bytes())

	return h1

}

func (this *packOper) IsOk() error {
	p1 := this
	if p1.OperType < POBegin || p1.OperType > POEnd {
		return errors.New("PackOper OperType is error.")
	}

	h1 := p1.GetSha256()

	if p1.HashHex == h1 {
		return nil
	} else {
		return errors.New("The packoper hash is inconsistent.")
	}
}

func GetRDCInt64() int64 {
	rd, _ := randc.Int(randc.Reader,
		big.NewInt(9223372036854775807))
	return rd.Int64()
}

func GetRDCInt8() int8 {
	rd, _ := randc.Int(randc.Reader,
		big.NewInt(255))
	return int8(rd.Int64())
}

func GetRDCbyte() byte {
	rd, _ := randc.Int(randc.Reader,
		big.NewInt(255))
	return byte(rd.Int64())
}

func GetRDCBytes(byteLen int) []byte {
	data := make([]byte, byteLen)
	for i := 0; i < byteLen; i++ {
		data[i] = GetRDCbyte()
	}
	return data
}

func GetRDF64() float64 {
	//var rd_s rand.Source = rand.NewSource(time.Now().Unix())
	//rnd := rand.New(rd_s)
	return float64(GetRDCInt64()) / 9223372036854775807.0
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
	return GetRDCBytes(byteLen)
}

func IsChangeCryKey(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == ChangeCryKey)
}

func IsPOVersion(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POVersion)
}

func GetKey(Data []byte) []byte {
	return (UnPack_Oper(Data).OperData)
}

func Int32ToBytes(i32 int32) []byte {
	bbuf := new(bytes.Buffer)
	err := binary.Write(bbuf, binary.LittleEndian, i32)
	if err != nil {
		panic(err)
	}
	return bbuf.Bytes()
}

func Int64ToBytes(data int64) []byte {
	bbuf := new(bytes.Buffer)
	err := binary.Write(bbuf, binary.LittleEndian, data)
	if err != nil {
		panic(err)
	}
	return bbuf.Bytes()
}

func CreatePackOperGen(data []byte) *packOper {

	rd := GetRDCInt64()

	pd := packOper{
		OperType: GenOper,
		//OperData: []byte(""),
		Data:    data,
		Rand:    rd,
		HashHex: "",
	}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func CreatePackOperChangeKey() *packOper {

	//rd := rand_s.Int63()
	rd := GetRDCInt64()

	key := GetRDBytes(32)

	pd := packOper{
		OperType: ChangeCryKey,
		OperData: []byte(key),
		Rand:     rd,
		HashHex:  "",
	}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func CreatePackOperVersion() *packOper {

	rd := GetRDCInt64()

	pd := packOper{
		OperType: POVersion,
		OperData: []byte(Version),
		Rand:     rd,
		HashHex:  "",
	}

	pd.HashHex = pd.GetSha256()

	return &pd
}

func jsonPacking_OperVersion() []byte {

	pd := CreatePackOperVersion()

	re, _ := json.Marshal(pd)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
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

func UnPack_Oper(data []byte) packOper {
	var msg packOper
	ix, re := Find0(data)
	if !re {
		return packOper{}
	}
	data2 := data[:ix]

	json.Unmarshal(data2, &msg)
	if debug_tag {
		p(msg)
		p("json: ", string(data2))
	}

	return msg
}

func jsonUnPack_OperGen(data []byte) ([]byte, error) {

	p1 := UnPack_Oper(data)

	return p1.Data[:], p1.IsOk()
}

func GetSha256Hex(data []byte) string {
	s1 := sha256.Sum256(data)
	return hex.EncodeToString(
		s1[:],
	)
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

func jsonUnpack(data []byte) ([]byte, error) {
	if begin_Dinfo == 0 {
		fmt.Println("Pack_Rand use.")
		begin_Dinfo = 1
	}
	return jsonUnPack_OperGen(data)
}

type Aes struct {
	cpr        cipher.Block
	cfbE, cfbD cipher.Stream
}

func CreateAes(key string) Aes {
	var v1 Aes

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
func (ap *Aespack) Unpack(data []byte) ([]byte, error) {
	d1, _ := base64.StdEncoding.DecodeString(string(data))
	jdata := ap.a.decrypter(d1)
	var pdata []byte
	var err error

	if IsChangeCryKey(jdata) {
		err = ap.setKey(GetKey(jdata))
		pdata = []byte{}
		return pdata, err
	}

	if IsPOVersion(jdata) {
		po1 := UnPack_Oper(jdata)
		err := po1.IsOk()
		if err != nil {
			return []byte{}, err
		}
		if po1.OperType == POVersion {
			if string(po1.OperData) == Version {
				pdata = []byte{}
				err = nil
			} else {
				pdata = []byte{}
				err = errors.New("Version is error.")
			}
		}
		return pdata, err
	}

	pdata, err = jsonUnpack(jdata)

	return pdata, err
}

func (ap *Aespack) setKey(key []byte) error {
	ap.a = CreateAes(string(key))
	return nil
}

func (ap *Aespack) ChangeCryKey() []byte {

	jdata := jsonPacking_OperChangeKey()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	ap.setKey(GetKey(jdata))

	return append([]byte(edata), 0)
}

func (ap *Aespack) IsTheVersionConsistent() []byte {

	jdata := jsonPacking_OperVersion()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	return append([]byte(edata), 0)
}

func CreateAesPack(key string) Aespack {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Error: The key is not 16, 24, or 32 bytes.")
		Logger.Fatalln("Error: The key is not 16, 24, or 32 bytes.")
		//panic(errors.New("Error: The key is not 16, 24, or 32 bytes."))
		//os.Exit(10)
	}
	ap1 := Aespack{}
	ap1.a = CreateAes(key)

	return ap1
}
