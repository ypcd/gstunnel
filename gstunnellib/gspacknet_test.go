package gstunnellib

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_gspacknet_pack_unpack(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	rawdata := []byte("123456")
	encrydata := pn.Packing(rawdata)
	decrydata, err := pn.Unpack(encrydata)
	CheckError_test(err, t)
	if !bytes.Equal(rawdata, decrydata) {
		t.Fatal("error.")
	}
}

func Test_gspacknet_pack_unpack_m(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)

	for i := 0; i < 10; i++ {
		rawdata := gsrand.GetRDBytes(50 * 1024)
		encrydata := pn.Packing(rawdata)
		decrydata, err := pn.Unpack(encrydata)
		CheckError_test(err, t)
		if !bytes.Equal(rawdata, decrydata) {
			t.Fatal("error.")
		}
	}
}

// 1
func Test_gspacknet_WriteEncryData_GetDecryData1(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	rawdata := gsrand.GetRDBytes(int(gsrand.GetRDCInt_max(30 * 1024)))
	encrydata := pn.Packing(rawdata)

	for i := 0; i < 10; i++ {
		fmt.Println("Hash:", gshash.GetSha256Hex(encrydata))
		pn.WriteEncryData(encrydata)
		data, err := pn.GetDecryData()
		CheckError_test(err, t)
		if !bytes.Equal(rawdata, data) {
			t.Fatal("error.")
		}
	}
}

// <1
func Test_gspacknet_WriteEncryData_GetDecryData1_2(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	rawdata := gsrand.GetRDBytes(int(100 + gsrand.GetRDCInt_max(30*1024)))
	encrydata := pn.Packing(rawdata)

	pn.WriteEncryData(encrydata[:100])
	data, err := pn.GetDecryData()
	CheckError_panic(err)
	if data != nil {
		log.Fatal("Error.")
	}
	pn.WriteEncryData(encrydata[100:])
	data, err = pn.GetDecryData()
	CheckError_test(err, t)
	if !bytes.Equal(rawdata, data) {
		t.Fatal("error.")
	}
}

// >1
func Test_gspacknet_WriteEncryData_GetDecryData1_3(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	rawdata := gsrand.GetRDBytes(int(gsrand.GetRDCInt_max(30 * 1024)))
	encrydata := pn.Packing(rawdata)
	rawdata2 := gsrand.GetRDBytes(int(1000 + gsrand.GetRDCInt_max(30*1024)))
	encrydata2 := pn.Packing(rawdata2)

	pn.WriteEncryData(encrydata)
	pn.WriteEncryData(encrydata2[:1000])
	data, err := pn.GetDecryData()
	CheckError_panic(err)
	if !bytes.Equal(rawdata, data) {
		t.Fatal("error.")
	}
	pn.WriteEncryData(encrydata2[1000:])
	data, err = pn.GetDecryData()
	CheckError_test(err, t)
	if !bytes.Equal(rawdata2, data) {
		t.Fatal("error.")
	}
}

// ==2
func Test_gspacknet_WriteEncryData_GetDecryData1_4(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	rawdata := gsrand.GetRDBytes(int(gsrand.GetRDCInt_max(30 * 1024)))
	encrydata := pn.Packing(rawdata)
	rawdata2 := gsrand.GetRDBytes(int(1000 + gsrand.GetRDCInt_max(30*1024)))
	encrydata2 := pn.Packing(rawdata2)

	pn.WriteEncryData(encrydata)
	pn.WriteEncryData(encrydata2)
	data, err := pn.GetDecryData()
	CheckError_panic(err)
	rawdataList := append(rawdata, rawdata2...)
	if !bytes.Equal(rawdataList, data) {
		t.Fatal("error.")
	}
}

// ==1000
func Test_gspacknet_WriteEncryData_GetDecryData1_5(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)

	rawdataList := []byte{}
	for i := 0; i < 1000; i++ {
		rawdata := gsrand.GetRDBytes(int(gsrand.GetRDCInt_max(30 * 1024)))
		encrydata := pn.Packing(rawdata)
		pn.WriteEncryData(encrydata)
		rawdataList = append(rawdataList, rawdata...)
	}

	data, err := pn.GetDecryData()
	CheckError_panic(err)
	if !bytes.Equal(rawdataList, data) {
		t.Fatal("error.")
	}
}

func Test_gspacknet_WriteEncryData_GetDecryData2(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)

	for i := 0; i < 10000*1; i++ {
		//fmt.Println("Hash:", gshash.GetSha256Hex(encrydata))
		rawdata := gsrand.GetRDBytes(int(gsrand.GetRDCInt_max(65535 - 10000)))
		encrydata := pn.Packing(rawdata)
		pn.WriteEncryData(encrydata)
		data, err := pn.GetDecryData()
		CheckError_test(err, t)
		if !bytes.Equal(rawdata, data) {
			t.Fatal("error.")
		}
	}
}

func nouse_Test_gspacknet_WriteEncryData_GetDecryData2_2(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	fd, err := os.Create("in.data")
	CheckErrorEx_panic(err)
	defer fd.Close()

	for i := 0; i < 10000*1; i++ {
		//fmt.Println("Hash:", gshash.GetSha256Hex(encrydata))
		rawdata := gsrand.GetRDBytes(int(gsrand.GetRDCInt_max(65535 - 10000)))
		encrydata := pn.Packing(rawdata)
		fd.Write([]byte(fmt.Sprintf("%d\n", len(encrydata))))
		pn.WriteEncryData(encrydata)
		data, err := pn.GetDecryData()
		CheckError_test(err, t)
		if !bytes.Equal(rawdata, data) {
			t.Fatal("error.")
		}
	}
}

func Test_gspacknet_gspack(t *testing.T) {
	gspack := NewGsPack(g_key_Default)
	gspacknet := NewGsPackNet(g_key_Default)

	rawdata := []byte("123456")
	encrydata := gspacknet.Packing(rawdata)
	decrydata, err := gspack.Unpack(encrydata)
	CheckError_test(err, t)
	if !bytes.Equal(rawdata, decrydata) {
		t.Fatal("Error.")
	}
}

func Test_gspacknet_to_gspacknetimp(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	rawdata := []byte("123456")
	//rawdata = append(rawdata, 0)

	pn.WriteEncryData(rawdata)
	gpn, ok := pn.(*gsPackNetImp)
	if !ok {
		t.Fatal("Error.")
	}
	data, err := gpn.GetDecryData()
	CheckError_test(err, t)
	_ = data
}

func Test_gspacknet_GetDecryData_data_len0(t *testing.T) {
	pn := NewGsPackNet(g_key_Default)
	rawdata := []byte{}
	//rawdata = append(rawdata, 0)

	pn.WriteEncryData(rawdata)
	gpn, ok := pn.(*gsPackNetImp)
	if !ok {
		t.Fatal("Error.")
	}
	data, err := gpn.GetDecryData()
	CheckError_test(err, t)
	_ = data
}
