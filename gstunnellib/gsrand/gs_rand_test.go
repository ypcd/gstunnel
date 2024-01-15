package gsrand

import (
	randc "crypto/rand"
	"fmt"
	"log"
	"testing"
)

func Test_GetRDF64(t *testing.T) {

	f1 := GetRDF64()

	//var rd_s rand.Source = rand.NewSource(time.Now().Unix())

	for i := 0; i < 8; i++ {
		//t.Log(rd_s.Int63())
		//t.Log(float64(rd_s.Int63()) / math.Pow(2, 64))
		//t.Log(GetRDF64())
		t.Log(GetRDInt8())
	}
	t.Log(GetRDBytes(32))
	if len(GetRDBytes(32)) == 32 && f1 < 1 {
		t.Log("OK.")
	} else {
		t.Error("error.")
	}
}

func Test_GetRDF64_2(t *testing.T) {

	var total float64
	loopNum := 10000 * 100
	for i := 0; i < loopNum; i++ {
		total += GetRDF64()
	}
	log.Println(total, total/float64(loopNum))

	if (total / float64(loopNum)) <= 0.498 {
		panic("error")
	}
}

/*
func Test_GetRDF64_3(t *testing.T) {

	//var total float64
	loopNum := 10000
	for i := 0; i < loopNum; i++ {
		//fmt.Println(GetRDF64())
		fmt.Println(getRDF64_2())
	}

}

func Test_GetRDF64_3_2(t *testing.T) {

		var total float64
		loopNum := 10000 * 100
		for i := 0; i < loopNum; i++ {
			total += getRDF64_2()
		}
		log.Println(total, total/float64(loopNum))

		if (total / float64(loopNum)) <= 0.49 {
			panic("error")
		}
	}
*/
func Test_IntsToBytes(t *testing.T) {
	t.Log(Int32ToBytes(GetRDInt32()))
	t.Log(Int64ToBytes(GetRDInt64()))

}

func Test_GetRDBytes(t *testing.T) {
	bs := GetRDBytes(1000000)
	sum := uint64(0)
	for i := range bs {
		sum += uint64(bs[i])
	}

	t.Log(sum, sum/uint64(len(bs)))
	if 123 < sum/uint64(len(bs)) && sum/uint64(len(bs)) < 130 {

	} else {
		t.Error()
	}

}

func Test_GetRDInt8(t *testing.T) {
	for i := 0; i < 1000; i++ {
		re := GetRDInt8()
		if re < -128 || re > 127 {
			t.Log(re)
			t.Fatal("error.")
		}
	}
}

func Test_GetRDInt16(t *testing.T) {
	for i := 0; i < 100000*6; i++ {
		if GetRDInt16() < -32768 || int(GetRDInt16()) > 32767 {
			t.Fatal("error.")
		}
	}
}

func Test_GetRDInt_max(t *testing.T) {
	for i := 0; i < 10; i++ {
		rd := GetRDInt_max(6)
		t.Log(rd)
	}
}

func Test_getrandString(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Log(GetrandString(100))
	}
}

func Test_getrandString2(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(GetrandStringPlus(32))
	}
}

func Test_getbytes(t *testing.T) {
	rd := make([]byte, 16)
	_, err := randc.Reader.Read(rd)
	if err != nil {
		panic(err)
	}
	fmt.Println(rd)

}

func Test_GetRD_netPort(t *testing.T) {
	for i := 0; i < 10000; i++ {
		fmt.Println(GetRD_netPortNumber())
	}
}
