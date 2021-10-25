package gsrand

import (
	"bytes"
	randc "crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
)

var big255 *big.Int = big.NewInt(255)

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
		big255)
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
		v1[i] = strpool[GetRDCInt_max(int64(len(strpool)))]
	}
	return string(v1)
}
