package gspackoper

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	//"gstunnellib"

	//"gstunnellib"

	//"gstunnellib"
	"sync"

	"gstunnel/gstunnellib/gsbase"
	. "gstunnel/gstunnellib/gsrand"

	"google.golang.org/protobuf/proto"
)

const version string = gsbase.Version

var debug_tag bool

var p func(...interface{}) (int, error)

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

const (
	POBegin        uint32 = 0
	POGenOper      uint32 = 1
	POChangeCryKey uint32 = 2
	POVersion      uint32 = 3
	POEnd          uint32 = 4
)

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

func (po *packOper) GetSha256_old() [32]byte {
	return GetSha256_32bytes(append(Intu8ToBytes(uint8(po.OperType)), append(po.OperData, append(po.Data, po.Rand...)...)...))
}

func (po *packOper) GetSha256() []byte {
	v1 := po.GetSha256_buf()
	//return hex.EncodeToString(v1[:])
	return v1[:]
}

//GetSha256_buf速度比GetSha256_old快15%-20%。
func (po *packOper) GetSha256_buf() []byte {

	var buf bytes.Buffer
	buf.Write(Intu8ToBytes(uint8(po.OperType)))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(po.Rand)

	sha := GetSha256_32bytes(buf.Bytes())

	return sha[:]

}

//为了更好的安全性，默认情况下没有使用GetSha256_pool函数。
//GetSha256_pool的速度比GetSha256_buf的速度快10%-15%。

func (po *packOper) GetSha256_pool() []byte {

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.Write(Intu8ToBytes(uint8(po.OperType)))
	buf.Write(po.OperData)
	buf.Write(po.Data)
	buf.Write(po.Rand)

	h1 := GetSha256_32bytes(buf.Bytes())

	return h1[:]

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

func createPackOperGen(data []byte) *packOper {

	rd := GetRDCBytes(8)

	pd := packOper{PackOperPro: PackOperPro{
		OperType: POGenOper,
		Data:     data,
		Rand:     rd,
	}}

	pd.HashHex = pd.GetSha256()
	return &pd
}

func createPackOperChangeKey() *packOper {

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

func createPackOperVersion() *packOper {

	rd := GetRDCBytes(8)

	pd := packOper{PackOperPro: PackOperPro{
		OperType: POVersion,
		OperData: []byte(version),
		Rand:     rd,
	}}

	pd.HashHex = pd.GetSha256()

	return &pd
}

func JsonPacking_OperVersion() []byte {

	pd := createPackOperVersion()

	re, _ := proto.Marshal(&pd.PackOperPro)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

func JsonPacking_OperChangeKey() []byte {

	pd := createPackOperChangeKey()

	re, _ := proto.Marshal(&pd.PackOperPro)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

func jsonPacking_OperGen(data []byte) []byte {

	pd := createPackOperGen(data)

	re, _ := proto.Marshal(&pd.PackOperPro)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
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

func jsonUnPack_OperGen(data []byte) ([]byte, error) {

	p1 := UnPack_Oper(data)

	return p1.Data[:], p1.IsOk()
}

func JsonPacking(data []byte) []byte {
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
