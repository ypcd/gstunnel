package gsrand

import (
	"bytes"
	randc "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"math"
	"math/big"
)

var g_big255 *big.Int = big.NewInt(255)
var g_bigmaxint64 *big.Int = big.NewInt(math.MaxInt64)

func GetRDInt_max(max int64) int64 {
	rd, _ := randc.Int(randc.Reader,
		big.NewInt(max))
	return rd.Int64()
}

func GetRDInt64() int64 {
	rd, _ := randc.Int(randc.Reader,
		g_bigmaxint64)
	return rd.Int64()
}

func GetRDInt8() int8 {
	rd, _ := randc.Int(randc.Reader,
		g_big255)
	return int8(rd.Int64())
}

func GetRDbyte() byte {
	rd, _ := randc.Int(randc.Reader,
		g_big255)
	return byte(rd.Int64())
}

func getRDBytes_old(byteLen int) []byte {
	data := make([]byte, byteLen)
	for i := 0; i < byteLen; i++ {
		data[i] = GetRDbyte()
	}
	return data
}

func GetRDBytes(byteLen int) []byte {
	data := make([]byte, byteLen)
	rlen, err := randc.Reader.Read(data)
	if err != nil {
		panic(err)
	}
	if rlen != byteLen {
		panic("rlen != byteLen")
	}
	return data
}

func GetRDF64() float64 {
	return float64(GetRDInt64()) / float64(9223372036854775807.0)
}

/*
	func getRDF64_2() float64 {
		rd, _ := randc.Int(randc.Reader,
			g_bigmaxint64)
		f, _ := rd.f.Float64()
		return f / float64(9223372036854775807.0)
	}
*/
func GetRDInt16() int16 {
	return int16(GetRDF64() * math.MaxUint16)
}

func GetRDInt32() int32 {
	return int32(GetRDF64() * math.MaxUint32)
}

func GetRDUint16() uint16 {
	return uint16(GetRDF64() * math.MaxUint16)
}

func GetRD_netPortNumber() uint16 {
	port := GetRDUint16()
	if port < 10000 {
		return 10000 + port
	}
	return port
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

func GetrandString(Len int) string {
	var strpool = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	v1 := make([]byte, Len)
	for i := 0; i < Len; i++ {
		v1[i] = strpool[GetRDInt_max(int64(len(strpool)))]
	}
	return string(v1)
}

func GetrandStringPlus(Len int) string {
	var strpool = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~!@#$%^&*()-_=+[{]}|;:',<.>/?")
	v1 := make([]byte, Len)
	for i := 0; i < Len; i++ {
		v1[i] = strpool[GetRDInt_max(int64(len(strpool)))]
	}
	return string(v1)
}

func GetrandBytesArray(Len int) [3][]byte {
	rdlist := [3][]byte{}
	for i := 0; i < 3; i++ {
		rdlist[i] = []byte(GetrandStringPlus(Len))
	}

	return rdlist
}

func GetRDKeyBase64(len int) string {
	return base64.StdEncoding.EncodeToString(GetRDBytes(len))
}
