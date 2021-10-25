package gspackoper

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"timerm"

	//. "gstunnellib"
	. "gstunnellib/gsrand"
	"math/rand"
	"strings"
	"testing"
	"unsafe"
)

func Test_JPackandun(t *testing.T) {
	t.Log("....................")

	v1 := []byte{1, 2, 3}
	j1 := JsonPacking(v1)
	t.Log(j1)
	t.Log(string(j1))
	dp1, _ := jsonUnpack(j1)
	t.Log(dp1)
	t.Log(string(v1) == string(dp1))
	if string(v1) == string(dp1) {
		t.Log("ok.")
	} else {
		t.Error()
	}
}
func Test_PackRand(t *testing.T) {
	d1 := []byte("123Abc")
	p1 := PackRand{Pack{d1}, rand.Int63()}
	j1, _ := json.Marshal(p1)

	p2 := PackRand{}
	json.Unmarshal(j1, &p2)
	d2 := p2.Data
	if string(d1) == string(d2) {
		t.Log("d1 == d2, ok.")
	}
	if p1.Rand == p2.Rand {
		t.Log("p1.Rand == p2.Rand, ok.")
	}
	t.Log(d2)
}

func Test_GetSha256Hex(t *testing.T) {
	data := []byte{}

	buf := make([]byte, 100)

	for i := range [100]int{} {
		_ = i
		blen := binary.PutVarint(buf, GetRDInt64())
		data = append(data, buf[:blen]...)
	}

	t.Log(GetSha256Hex(data))

}
func Test_JsonPacking_OperChangeKey(t *testing.T) {
	b1 := JsonPacking_OperChangeKey()
	p1 := UnPack_Oper(b1)
	_ = p1
	t.Log("Ok.")
}

func Test_JsonPacking_OperGen(t *testing.T) {

	data := GetRDBytes(10024)

	b1 := jsonPacking_OperGen(data)
	p1 := UnPack_Oper(b1)
	_ = p1

	unb2, _ := jsonUnPack_OperGen(b1)
	h1 := GetSha256Hex(unb2)
	if GetSha256Hex(data) == h1 {
		t.Log("Ok.")
	} else {
		t.Error("Error.")
	}
	_ = p1

}

func Test_JPackandUnpack_oper(t *testing.T) {
	t.Log("....................")

	v1 := GetRDBytes(953512)
	j1 := JsonPacking(v1)
	t.Log(j1)
	t.Log(string(j1))
	dp1, _ := jsonUnpack(j1)
	t.Log(dp1)
	t.Log(GetSha256Hex(v1) == GetSha256Hex(dp1))
	if GetSha256Hex(v1) == GetSha256Hex(dp1) {
		t.Log("ok.")
	} else {
		t.Error()
	}
	_ = v1
}

func Test_bytesjoin(t *testing.T) {

	//const totalLoop int = 10000 * 100

	po1 := createPackOperChangeKey()

	po1.Data = GetRDBytes(1024)

	b1 := po1.GetSha256_old()
	b2 := po1.GetSha256()
	b3 := po1.GetSha256_buf()
	b4 := po1.GetSha256_pool()

	if bytes.Compare(b1[:], b2) != 0 {
		t.Error("error.")
	}
	if bytes.Compare(b2, b3) != 0 {
		t.Fatal()
	}
	if bytes.Compare(b3, b4) != 0 {
		t.Fatal()
	}

	if strings.Compare(hex.EncodeToString(b1[:]), hex.EncodeToString(b2[:])) != 0 {
		t.Error("error.")
	}

}

func Test_po_size(t *testing.T) {

	//const totalLoop int = 10000 * 100

	po2 := createPackOperChangeKey()

	_ = po2

	pd := createPackOperGen([]byte{})

	re, _ := json.Marshal(pd)

	_ = re

	pd2 := createPackOperGen_po1([]byte{})

	re2, _ := json.Marshal(pd2)

	t.Log(string(re), len(re))
	t.Log(string(re2), len(re2))
	t.Log("po size:", unsafe.Sizeof(po2))
	t.Log("po size:", unsafe.Sizeof(po2))

}

func Test_compress_un(t *testing.T) {
	ap1 := NewCompresser()
	data1 := GetRDCBytes(1024 * 6)

	const data2 = `<?xml version="1.0"?>
<book>
	<meta name="title" content="The Go Programming Language"/>
	<meta name="authors" content="Alan Donovan and Brian Kernighan"/>
	<meta name="published" content="2015-10-26"/>
	<meta name="isbn" content="978-0134190440"/>
	<data>...</data>
</book>
`

	pd := createPackOperGen([]byte(data1))

	re, _ := json.Marshal(pd)

	_ = re

	data := []byte(re)

	rt := timerm.CreateRecoTime()
	t.Log(rt.Run())
	cdata := ap1.compress(data)
	t.Log(rt.Run())
	t.Log("compress:", float32(len(cdata))/float32(len(data)))

	undata := ap1.uncompress(cdata)
	if bytes.Equal(data, undata) {
		t.Log("ok.")
	} else {
		t.Log("Error.")
	}

}

func Test_po_is(t *testing.T) {
	poData := jsonPacking_OperGen(nil)
	if IsChangeCryKey(poData) {
		t.Fatal("error.")
	}
	if !IsPOGen(poData) {
		t.Fatal()
	}
	if IsPOVersion(poData) {
		t.Fatal()
	}

	pog := createPackOperGen(nil)
	if pog.IsChangeCryKey() {
		t.Fatal("error.")
	}
	if !pog.IsPOGen() {
		t.Fatal()
	}
	if pog.IsPOVersion() {
		t.Fatal()
	}

}

func Test_po1(t *testing.T) {
	data := []byte("123 ABcd test.")

	po := createPackOperGen_po1(nil)
	if po.IsOk() != nil {
		t.Fatal()
	}
	pdata := jsonPacking_OperGen_po1(data)
	undata, _ := jsonUnpack_po1(pdata)
	if !bytes.Equal(data, undata) {
		t.Fatal()
	}
}

func Test_json__(t *testing.T) {
	pvdata := JsonPacking_OperVersion()
	if !IsPOVersion(pvdata) {
		t.Fatal()
	}

	pov := createPackOperVersion()
	if pov.IsOk() != nil {
		t.Fatal()
	}

	pdata := JsonPacking_OperChangeKey()
	if !IsChangeCryKey(pdata) {
		t.Fatal()
	}
}