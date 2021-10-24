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
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"google.golang.org/protobuf/proto"
)

const Version string = "V4.5.6"

//var Version string = "V4.5.3"

var p func(...interface{}) (int, error)

var begin_Dinfo int = 0

var debug_tag bool

var commonIV = []byte{171, 158, 1, 73, 31, 98, 64, 85, 209, 217, 131, 150, 104, 219, 33, 220}

/*
type Gst_logger struct {
	*log.logger
}

func (plog *Gst_logger) Println(v ...interface{}) {
	plog.logger.Println(v...)
	_, _ = fmt.Println(v...)
}

var logger *Gst_logger
*/

var logger *log.Logger

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

var Info_protobuf bool = true

func Nullprint(v ...interface{}) (int, error)                       { return 1, nil }
func Nullprintf(format string, a ...interface{}) (n int, err error) { return 1, nil }

func init() {
	debug_tag = false
	p = Nullprint

	logger = CreateFileLogger("gstunnellib.data.log")
	//debug_tag = true
	//p = fmt.Println
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		logger.Fatalln("Fatal error:", err.Error())
	}
}

func CreateFileLogger(FileName string) *log.Logger {
	lf, err := os.Create(FileName)
	CheckError(err)

	logger = log.New(lf, "", log.Lshortfile|log.LstdFlags|log.Lmsgprefix)
	return logger
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
	POBegin        uint32 = 0
	POGenOper      uint32 = 1
	POChangeCryKey uint32 = 2
	POVersion      uint32 = 3
	POEnd          uint32 = 4
)

type packOper_1 struct {
	OperType uint32
	OperData []byte
	Data     []byte
	Rand     []byte //[8]byte
	//HashHex  string
	HashHex []byte //[32]byte
}

func (po *packOper_1) IsOk() error {
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

func (po *packOper_1) GetSha256_buf() [32]byte {

	var buf bytes.Buffer
	buf.Write(Intu8ToBytes(uint8(po.OperType)))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write([]byte(po.Rand))

	return GetSha256_32bytes(buf.Bytes())

}

func (po *packOper_1) GetSha256() []byte {
	v1 := po.GetSha256_buf()
	//return hex.EncodeToString(v1[:])
	return v1[:]
}

/*
//Data, Rand, HashHex, opertype, operData
type packOper struct {
	OperType uint8
	OperData []byte
	Data     []byte
	Rand     []byte
	HashHex  []byte
}
*/

type packOper struct {
	PackOperPro
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

/*
func GetSha256_old_po1(po *packOper_1) string {
	return GetSha256Hex(append(Intu8ToBytes(po.OperType), append(po.OperData, append(po.Data, Int64ToBytes(po.Rand)...)...)...))
}
*/
func (po *packOper) GetSha256_old() [32]byte {
	return GetSha256_32bytes(append(Intu8ToBytes(uint8(po.OperType)), append(po.OperData, append(po.Data, po.Rand...)...)...))
}

func (po *packOper) GetSha256() []byte {
	v1 := po.GetSha256_buf()
	//return hex.EncodeToString(v1[:])
	return v1[:]
}

//GetSha256_buf速度比GetSha256_old快15%-20%。
func (po *packOper) GetSha256_buf() [32]byte {

	var buf bytes.Buffer
	buf.Write(Intu8ToBytes(uint8(po.OperType)))
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

	buf.Write(Intu8ToBytes(uint8(po.OperType)))
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

func CreatePackOperGen_po1(data []byte) *packOper_1 {

	rd := GetRDCBytes(8)

	pd := packOper_1{
		OperType: POGenOper,
		//OperData: []byte(""),
		Data: data,
		Rand: rd,
		//HashHex: nil,
	}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func CreatePackOperGen(data []byte) *packOper {

	rd := GetRDCBytes(8)
	/*
		pop := PackOperPro{
			OperType: POGenOper,
			//OperData: []byte(""),
			Data: data,
			Rand: rd,
		}
	*/

	pd := packOper{PackOperPro: PackOperPro{
		OperType: POGenOper,
		//OperData: []byte(""),
		Data: data,
		Rand: rd,
	}}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func CreatePackOperChangeKey() *packOper {

	rd := GetRDCBytes(8)

	key := GetRDBytes(32)

	pd := packOper{PackOperPro: PackOperPro{
		OperType: POChangeCryKey,
		OperData: []byte(key),
		Rand:     rd,
	}}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func CreatePackOperVersion() *packOper {

	rd := GetRDCBytes(8)

	pd := packOper{PackOperPro: PackOperPro{
		OperType: POVersion,
		OperData: []byte(Version),
		Rand:     rd,
	}}

	pd.HashHex = pd.GetSha256()

	return &pd
}

func jsonPacking_OperVersion() []byte {

	pd := CreatePackOperVersion()

	re, _ := proto.Marshal(&pd.PackOperPro)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
}

func jsonPacking_OperChangeKey() []byte {

	pd := CreatePackOperChangeKey()

	re, _ := proto.Marshal(&pd.PackOperPro)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
}

func jsonPacking_OperGen(data []byte) []byte {

	pd := CreatePackOperGen(data)

	re, _ := proto.Marshal(&pd.PackOperPro)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
}

func jsonPacking_OperGen_po1(data []byte) []byte {

	pd := CreatePackOperGen_po1(data)

	re, _ := json.Marshal(pd)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return append(re, 0)
}

func UnPack_Oper(data []byte) *packOper {
	var msg packOper
	//var msg2 PackOperPro
	//ix, re := Find0(data)
	/*if !re {
		return &packOper{}
	}
	*/
	data2 := data

	proto.Unmarshal(data2, &msg.PackOperPro)

	//msg.PackOperPro = msg2
	if debug_tag {
		p(msg)
		p("json: ", string(data2))
	}

	return &msg
}

func UnPack_Oper_po1(data []byte) *packOper_1 {
	var msg packOper_1
	//var msg2 PackOperPro
	//ix, re := Find0(data)
	/*if !re {
		return &packOper{}
	}
	*/
	data2 := data

	json.Unmarshal(data2, &msg)

	//msg.PackOperPro = msg2
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

func jsonUnPack_OperGen_po1(data []byte) ([]byte, error) {

	p1 := UnPack_Oper_po1(data)

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
	/*
		if begin_Dinfo == 0 {
			fmt.Println("Pack_Rand use.")
			begin_Dinfo = 1
		}
	*/
	return jsonUnPack_OperGen(data)
}

func jsonUnpack_po1(data []byte) ([]byte, error) {
	/*
		if begin_Dinfo == 0 {
			fmt.Println("Pack_Rand use.")
			begin_Dinfo = 1
		}
	*/
	return jsonUnPack_OperGen(data)
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
	a        gsaes
	fewriter *flate.Writer
	fereader io.ReadCloser
}

func (ap *aespack) compress2(data []byte) []byte   { return data }
func (ap *aespack) uncompress2(data []byte) []byte { return data }

func (ap *aespack) compress(data []byte) []byte {
	var b bytes.Buffer

	if ap.fewriter == nil {
		zw, err := flate.NewWriter(&b, 1)
		ap.fewriter = zw
		if err != nil {
			logger.Fatalln(err)
		}
	} else {
		ap.fewriter.Reset(&b)
	}

	zw := ap.fewriter

	if _, err := io.Copy(zw, bytes.NewReader(data)); err != nil {
		logger.Fatalln(err)
	}
	if err := zw.Close(); err != nil {
		logger.Fatalln(err)
	}

	return b.Bytes()
}

func (ap *aespack) uncompress(data []byte) []byte {
	var b bytes.Buffer

	if ap.fereader == nil {
		zr := flate.NewReader(bytes.NewReader(data))
		ap.fereader = zr
	} else {
		zr := ap.fereader
		if err := zr.(flate.Resetter).Reset(bytes.NewReader(data), nil); err != nil {
			logger.Fatalln(err)
		}
	}
	zr := ap.fereader

	if _, err := io.Copy(&b, zr); err != nil {
		logger.Fatalln(err)
	}
	if err := zr.Close(); err != nil {
		logger.Fatalln(err)
	}

	return b.Bytes()
}

func (ap *aespack) Packing(data []byte) []byte {
	jdata := jsonPacking(data)
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

	jdata := jsonPacking_OperChangeKey()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	ap.setKey(GetKey(jdata))

	return append([]byte(edata), 0)
}

func (ap *aespack) IsTheVersionConsistent() []byte {

	jdata := jsonPacking_OperVersion()
	crydata := ap.a.encrypter(jdata)
	edata := base64.StdEncoding.EncodeToString(crydata)

	return append([]byte(edata), 0)
}

func createAesPack(key string) *aespack {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Error: The key is not 16, 24, or 32 bytes.")
		logger.Fatalln("Error: The key is not 16, 24, or 32 bytes.")
	}
	ap1 := aespack{}
	ap1.a = createAes([]byte(key))
	ap1.fewriter = nil
	ap1.fereader = nil

	return &ap1
}

func NewGsPack(key string) GsPack {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Error: The key is not 16, 24, or 32 bytes.")
		logger.Fatalln("Error: The key is not 16, 24, or 32 bytes.")
	}
	ap1 := aespack{}
	ap1.a = createAes([]byte(key))
	ap1.fewriter = nil
	ap1.fereader = nil

	return &ap1
}
