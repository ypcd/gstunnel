/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package gstunnellib

import (
	"bytes"
	"compress/flate"
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
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"sync"
	"sync/atomic"
)

const Version string = "V3.8"

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

type GsConfig_1 struct {
	Listen             string
	Server             string
	Key                string
	Debug              bool
	Tmr_display_time   int
	Tmr_changekey_time int
	Mt_model           bool
}

type GsConfig struct {
	Listen             string
	Servers            []string
	Key                string
	Debug              bool
	Tmr_display_time   int
	Tmr_changekey_time int
	Mt_model           bool
}

func (gs *GsConfig) GetServer() string {
	return gs.Servers[GetRDCInt_max(int64(len(gs.Servers)))]
}
func (gs *GsConfig) GetServers() []string {
	return gs.Servers
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

	gsconfig.Debug = false
	gsconfig.Tmr_display_time = 5
	gsconfig.Tmr_changekey_time = 60
	gsconfig.Mt_model = true

	json.Unmarshal(buf, &gsconfig)
	/*
		if gsconfig.Tmr_display_time == 0 {
			gsconfig.Tmr_display_time = 5
		}
		if gsconfig.Tmr_changekey_time == 0 {
			gsconfig.Tmr_changekey_time = 60
		}
	*/
	if gsconfig.Servers == nil {
		Logger.Fatalln("gsconfig.Servers==nil")
	}
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
	POBegin        uint8 = 0
	POGenOper      uint8 = 1
	POChangeCryKey uint8 = 2
	POVersion      uint8 = 3
	POEnd          uint8 = 4
)

type packOper_1 struct {
	OperType uint8
	OperData []byte
	Data     []byte
	Rand     int64
	HashHex  string
	//HashHex [32]byte
}

//Data, Rand, HashHex, opertype, operData
type packOper struct {
	OperType uint8
	OperData []byte
	Data     []byte
	Rand     []byte
	HashHex  []byte
}

func GetSha256Hex(data []byte) string {
	s1 := sha256.Sum256(data)
	return hex.EncodeToString(
		s1[:],
	)
}

func GetSha256_32bytes(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func GetSha256_old_po1(po *packOper_1) string {
	return GetSha256Hex(append(Intu8ToBytes(po.OperType), append(po.OperData, append(po.Data, Int64ToBytes(po.Rand)...)...)...))
}

func (po *packOper) GetSha256_old() [32]byte {
	return GetSha256_32bytes(append(Intu8ToBytes(po.OperType), append(po.OperData, append(po.Data, po.Rand...)...)...))
}

func (po *packOper) GetSha256() []byte {
	v1 := po.GetSha256_buf()
	//return hex.EncodeToString(v1[:])
	return v1[:]
}

//GetSha256_buf速度比GetSha256_old快15%-20%。
func (po *packOper) GetSha256_buf() [32]byte {

	var buf bytes.Buffer
	buf.Write(Intu8ToBytes(po.OperType))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(po.Rand)

	return GetSha256_32bytes(buf.Bytes())

}

//为了更好的安全性，默认情况下没有使用GetSha256_pool函数。
//GetSha256_pool的速度比GetSha256_buf的速度快10%-15%。

func (po *packOper) GetSha256_pool() [32]byte {

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()
	//fmt.Println("CAP LEN:", buf.Cap(), buf.Len())

	buf.Write(Intu8ToBytes(po.OperType))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(po.Rand)

	h1 := GetSha256_32bytes(buf.Bytes())

	return h1

}

func (po *packOper) IsOk() error {
	p1 := po
	if p1.OperType < POBegin || p1.OperType > POEnd {
		return errors.New("PackOper OperType is error.")
	}

	h1 := p1.GetSha256()

	if bytes.Equal(p1.HashHex, h1) {
		return nil
	} else {
		return errors.New("The packoper hash is inconsistent.")
	}
}

func (po *packOper) IsChangeCryKey() bool {
	return po.OperType == POChangeCryKey
}

func (po *packOper) IsPOVersion() bool {
	return po.OperType == POVersion
}

func (po *packOper) IsPOGen() bool {
	return po.OperType == POGenOper
}

func GetRDCInt_max(max int64) int64 {
	rd, _ := randc.Int(randc.Reader,
		big.NewInt(max))
	return rd.Int64()
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
	return (UnPack_Oper(Data).OperType == POChangeCryKey)
}

func IsPOVersion(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POVersion)
}

func IsPOGen(Data []byte) bool {
	return (UnPack_Oper(Data).OperType == POGenOper)
}

func GetKey(Data []byte) []byte {
	return (UnPack_Oper(Data).OperData)
}

func Intu8ToBytes(v uint8) []byte {
	bbuf := new(bytes.Buffer)
	err := binary.Write(bbuf, binary.LittleEndian, v)
	if err != nil {
		panic(err)
	}
	return bbuf.Bytes()
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

func CreatePackOperGen_po1(data []byte) *packOper_1 {

	rd := GetRDCInt64()

	pd := packOper_1{
		OperType: POGenOper,
		//OperData: []byte(""),
		Data: data,
		Rand: rd,
		//HashHex: nil,
	}

	pd.HashHex = GetSha256_old_po1(&pd)
	return &pd
}

func CreatePackOperGen(data []byte) *packOper {

	rd := GetRDCBytes(8)

	pd := packOper{
		OperType: POGenOper,
		//OperData: []byte(""),
		Data: data,
		Rand: rd,
		//HashHex: nil,
	}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func CreatePackOperChangeKey() *packOper {

	//rd := rand_s.Int63()
	rd := GetRDCBytes(8)

	key := GetRDBytes(32)

	pd := packOper{
		OperType: POChangeCryKey,
		OperData: []byte(key),
		Rand:     rd,
		//HashHex:  "",
	}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func CreatePackOperVersion() *packOper {

	rd := GetRDCBytes(8)

	pd := packOper{
		OperType: POVersion,
		OperData: []byte(Version),
		Rand:     rd,
		//HashHex:  "",
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

func UnPack_Oper(data []byte) *packOper {
	var msg packOper
	ix, re := Find0(data)
	if !re {
		return &packOper{}
	}
	data2 := data[:ix]

	json.Unmarshal(data2, &msg)
	if debug_tag {
		p(msg)
		p("json: ", string(data2))
	}

	return &msg
}

func jsonUnPack_OperGen(data []byte) ([]byte, error) {

	p1 := UnPack_Oper(data)

	return p1.Data[:], p1.IsOk()
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
	a        Aes
	fewriter *flate.Writer
	fereader io.ReadCloser
}

func (ap *Aespack) compress2(data []byte) []byte   { return data }
func (ap *Aespack) uncompress2(data []byte) []byte { return data }

func (ap *Aespack) compress(data []byte) []byte {
	var b bytes.Buffer

	if ap.fewriter == nil {
		zw, err := flate.NewWriter(&b, 1)
		ap.fewriter = zw
		if err != nil {
			Logger.Fatalln(err)
		}
	} else {
		ap.fewriter.Reset(&b)
	}

	zw := ap.fewriter

	if _, err := io.Copy(zw, bytes.NewReader(data)); err != nil {
		Logger.Fatalln(err)
	}
	if err := zw.Close(); err != nil {
		Logger.Fatalln(err)
	}

	return b.Bytes()
}

func (ap *Aespack) uncompress(data []byte) []byte {
	var b bytes.Buffer

	if ap.fereader == nil {
		zr := flate.NewReader(bytes.NewReader(data))
		ap.fereader = zr
	} else {
		zr := ap.fereader
		if err := zr.(flate.Resetter).Reset(bytes.NewReader(data), nil); err != nil {
			Logger.Fatalln(err)
		}
	}
	zr := ap.fereader

	if _, err := io.Copy(&b, zr); err != nil {
		Logger.Fatalln(err)
	}
	if err := zr.Close(); err != nil {
		Logger.Fatalln(err)
	}

	return b.Bytes()
}

func (ap *Aespack) Packing(data []byte) []byte {
	jdata := jsonPacking(data)
	cdata := ap.compress(jdata)
	crydata := ap.a.encrypter(cdata)
	edata := base64.StdEncoding.EncodeToString(crydata)
	return append([]byte(edata), 0)
}
func (ap *Aespack) Unpack(data []byte) ([]byte, error) {
	d1, _ := base64.StdEncoding.DecodeString(string(data))
	jdata := ap.a.decrypter(d1)
	undata := ap.uncompress(jdata)
	jdata = undata
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
		//po1 := po
		/*
			err := po1.IsOk()
			if err != nil {
				return []byte{}, err
			}
		*/
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
	if po1.IsPOGen() {
		return po1.Data, err
	}

	return pdata, err
}

func (ap *Aespack) setKey(key []byte) error {
	ap.a = CreateAes(string(key))
	return nil
}

func (ap *Aespack) ChangeCryKey() []byte {

	jdata := jsonPacking_OperChangeKey()
	cdata := ap.compress(jdata)
	crydata := ap.a.encrypter(cdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	ap.setKey(GetKey(jdata))

	return append([]byte(edata), 0)
}

func (ap *Aespack) IsTheVersionConsistent() []byte {

	jdata := jsonPacking_OperVersion()
	cdata := ap.compress(jdata)
	crydata := ap.a.encrypter(cdata)
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
	ap1.fewriter = nil
	ap1.fereader = nil

	return ap1
}

type Gorou_status struct {
	status int32
}

const (
	gorou_s_begin = 0
	gorou_s_ok    = 1
	gorou_s_close = 2
	gorou_s_end   = 3
)

func (g *Gorou_status) IsOk() bool {
	v1 := atomic.LoadInt32(&g.status)
	return v1 == gorou_s_ok
}
func (g *Gorou_status) SetOk()    { atomic.SwapInt32(&g.status, gorou_s_ok) }
func (g *Gorou_status) SetClose() { atomic.SwapInt32(&g.status, gorou_s_close) }

func CreateGorouStatus() *Gorou_status {
	g1 := new(Gorou_status)
	g1.status = gorou_s_ok
	return g1
}
