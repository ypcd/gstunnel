package gstunnellib

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_gspacknet_pack_unpack(t *testing.T) {
	pn := NewGsPackNet(key_defult)
	rawdata := []byte("123456")
	encrydata := pn.Packing(rawdata)
	decrydata, err := pn.Unpack(encrydata)
	CheckError_test(err, t)
	if !bytes.Equal(rawdata, decrydata) {
		t.Fatal("error.")
	}
}

func Test_gspacknet_pack_unpack_m(t *testing.T) {
	pn := NewGsPackNet(key_defult)
	rawdata := []byte("123456")
	encrydata := pn.Packing(rawdata)

	for i := 0; i < 10; i++ {
		decrydata, err := pn.Unpack(encrydata)
		CheckError_test(err, t)
		if !bytes.Equal(rawdata, decrydata) {
			t.Fatal("error.")
		}
	}
}

func Test_gspacknet_WriteEncryData_GetDecryData(t *testing.T) {
	pn := NewGsPackNet(key_defult)
	rawdata := gsrand.GetRDBytes(int(gsrand.GetRDCInt_max(1024 * 1024)))
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

func Test_gspacknet_gspack(t *testing.T) {
	gspack := NewGsPack(key_defult)
	gspacknet := NewGsPackNet(key_defult)

	rawdata := []byte("123456")
	encrydata := gspacknet.Packing(rawdata)
	decrydata, err := gspack.Unpack(encrydata)
	CheckError_test(err, t)
	_ = decrydata
}
