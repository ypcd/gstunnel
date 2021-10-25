package gspackoper

import (
	"bytes"
	"encoding/json"
	"errors"

	. "gstunnellib/gsrand"
)

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
func GetSha256_old_po1(po *packOper_1) string {
	return GetSha256Hex(append(Intu8ToBytes(po.OperType), append(po.OperData, append(po.Data, Int64ToBytes(po.Rand)...)...)...))
}
*/

func createPackOperGen_po1(data []byte) *packOper_1 {

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

func jsonPacking_OperGen_po1(data []byte) []byte {

	pd := createPackOperGen_po1(data)

	re, _ := json.Marshal(pd)

	if debug_tag {
		p("pd: ", pd)
		p("json: ", string(re))
	}

	return re
}

func unPack_Oper_po1(data []byte) *packOper_1 {
	var msg packOper_1
	//var msg2 PackOperPro
	//ix, re := Find0(data)
	/*if !re {
		return &packOper{}
	}
	*/
	data2 := data

	err := json.Unmarshal(data2, &msg)
	_ = err

	//msg.PackOperPro = msg2
	if debug_tag {
		p(msg)
		p("json: ", string(data2))
	}

	return &msg
}

func jsonUnPack_OperGen_po1(data []byte) ([]byte, error) {

	p1 := unPack_Oper_po1(data)

	return p1.Data[:], p1.IsOk()
}

func jsonUnpack_po1(data []byte) ([]byte, error) {
	/*
		if begin_Dinfo == 0 {
			fmt.Println("Pack_Rand use.")
			begin_Dinfo = 1
		}
	*/
	return jsonUnPack_OperGen_po1(data)
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
