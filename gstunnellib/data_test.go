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
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"timerm"
	"unsafe"

	//"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var testi1 = 1

var wlist sync.WaitGroup

func init() {
	p = Nullprint
	//p = t.Log
	debug_tag = true
}

func getrand() int {
	rd := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	return rd.Int()
}

func Test_JPackandun(t *testing.T) {
	t.Log("....................")

	v1 := []byte{1, 2, 3}
	j1 := jsonPacking(v1)
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

func getsha1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%X", h.Sum(nil))
}

var unpbuf []byte

func Test_find0(t *testing.T) {
	data := make([]byte, 0)
	buf := make([]byte, 100)

	for i := 0; i < 1000; i++ {
		blen := binary.PutVarint(buf, GetRDInt64())
		data = append(data, buf[:blen]...)
	}
	data[rand.Intn(len(data))] = 0
	ix, _ := Find0(data)
	ix2 := bytes.IndexByte(data, 0)
	t.Log(ix, ix2)
	if ix == ix2 {
		t.Log("ok.")
	} else {
		t.Error("error.")
	}

}

type ftype func()

func mtF(frun ftype) {
	for i := 0; i < 16; i++ {
		go frun()
		wlist.Add(1)
	}
	p("go")
	wlist.Wait()
}

func Test_Aest(t *testing.T) {

	fbuf := GetRDCBytes(1024 * 1024)

	a1 := createAes([]byte(getrandString(32)))

	tmp := a1.encrypter(fbuf)
	outbuf := a1.decrypter(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		t.Log("ok.")
		t.Log(getsha1(fbuf))
	} else {
		t.Error()
	}

}

func aest() {

	fbuf := GetRDCBytes(1024)

	a1 := createAes([]byte(getrandString(32)))
	tmp := a1.encrypter(fbuf)
	outbuf := a1.decrypter(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		p("ok.", getrand())

	} else {
		log.Fatal("Error")
	}
	wlist.Done()
}

func Test_MtAest(t *testing.T) {
	mtF(aest)
}

func Test_Aestpack(t *testing.T) {

	fbuf := GetRDCBytes(1024 * 1024)

	a1 := createAesPack(getrandString(32))

	tmp := a1.Packing(fbuf)
	t.Log(tmp[len(tmp)-3:])
	outbuf, _ := a1.Unpack(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		t.Log("ok.")
		t.Log(getsha1(fbuf))
	} else {
		t.Error()
	}

}

func aestpack() {

	fbuf := GetRDCBytes(1024 * 1024)

	a1 := createAesPack(getrandString(32))
	tmp := a1.Packing(fbuf)
	outbuf, _ := a1.Unpack(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		p("ok.", getrand())
	} else {
		log.Fatal("error")
	}
	wlist.Done()
}

func Test_MtAestpack(t *testing.T) {
	mtF(aestpack)
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

func Test_binConv(t *testing.T) {
	buf := make([]byte, 0)
	p1 := &buf
	t.Log(p1)
	bbuf := bytes.NewBuffer(buf)
	err := binary.Write(bbuf, binary.LittleEndian, int32(123))
	t.Log(err)
	t.Log(&buf, buf)
}

func Test_jsonPacking_OperChangeKey(t *testing.T) {
	b1 := jsonPacking_OperChangeKey()
	p1 := UnPack_Oper(b1)
	_ = p1
	t.Log("Ok.")
}

func Test_jsonPacking_OperGen(t *testing.T) {

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

func Test_Aespack_changeCryKey(t *testing.T) {
	key := GetRDBytes(32)
	ap1 := createAesPack(string(key))
	cp1 := ap1.a
	key2 := GetRDBytes(32)
	ap1.setKey(key2)
	cp2 := ap1.a

	t.Log(key)
	t.Log(key2)
	t.Log(cp1, cp2)
	if cp1 != cp2 {
		t.Log("Ok.")
	} else {
		t.Error("Error.")
	}
	_ = key
}

func Test_Aespack_changeCryKey2(t *testing.T) {
	key := GetRDBytes(32)
	ap1 := createAesPack(string(key))
	cp1 := ap1.a

	ap1.ChangeCryKey()
	cp2 := ap1.a

	t.Log(key)
	//	t.Log(key2)
	t.Log(cp1, cp2)
	if cp1 != cp2 {
		t.Log("Ok.")
	} else {
		t.Error("Error.")
	}
	_ = key
}

func Test_JPackandUnpack_oper(t *testing.T) {
	t.Log("....................")

	v1 := GetRDBytes(953512)
	j1 := jsonPacking(v1)
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

func Test_AesEncryDecry(t *testing.T) {
	d1 := GetRDBytes(15834)
	key1 := GetRDBytes(32)
	key2 := GetRDBytes(32)

	ap1 := createAesPack(string(key1))
	ap2 := createAesPack(string(key2))

	ed1 := ap1.a.encrypter(d1)
	ed2 := ap2.a.encrypter(d1)

	dd1, dd2 := ap1.a.decrypter(ed1), ap2.a.decrypter(ed2)

	t.Log(dd1)
	t.Log("------------------------------------------------------------------------------------------------------------------------------")
	t.Log(dd2)

	if bytes.Equal(dd1, d1) && bytes.Equal(d1, dd2) {
		t.Log("ok.")
	} else {
		t.Error("Error.")
	}
}

func Test_wbuf(t *testing.T) {
	wbuf := []byte{1, 2, 3, 4, 5, 6}
	t.Log(wbuf)
	wbuf = wbuf[len(wbuf):]
	t.Log(wbuf)
}

func Test_poversion(t *testing.T) {
	//d1 := GetRDBytes(15834)
	key1 := GetRDBytes(32)
	//key2 := GetRDBytes(32)

	ap1 := createAesPack(string(key1))
	//ap2 := createAesPack(string(key2))

	pb1 := ap1.IsTheVersionConsistent()
	unb1, err := ap1.Unpack(pb1)
	if err == nil {
		t.Log("ok.")
	} else {
		t.Error("Error.")
	}
	_, _ = unb1, err

}

func Test_bytesjoin(t *testing.T) {

	//const totalLoop int = 10000 * 100

	po1 := CreatePackOperChangeKey()

	po1.Data = GetRDBytes(1024)

	/*
		if po1.GetSha256_old() != po1.GetSha256() {
			t.Error("error.")
		}

		b1 := po1.GetSha256_old()
		b2 := po1.GetSha256()

		if hex.EncodeToString(b1[:]) != hex.EncodeToString(b2[:]) {
			t.Error("error.")
		}
	*/
}

func Test_po_size(t *testing.T) {

	//const totalLoop int = 10000 * 100

	key1 := GetRDBytes(32)
	//key2 := GetRDBytes(32)

	ap1 := createAesPack(string(key1))

	po1 := ap1.Packing([]byte{})

	po2 := CreatePackOperChangeKey()

	_ = po2

	pd := CreatePackOperGen([]byte{})

	re, _ := json.Marshal(pd)

	_ = re

	pd2 := CreatePackOperGen_po1([]byte{})

	re2, _ := json.Marshal(pd2)

	t.Log(string(re), len(re))
	t.Log(string(re2), len(re2))
	t.Log("po size:", len(po1), unsafe.Sizeof(po2))
	t.Log("po size:", len(po1), unsafe.Sizeof(po2))

}

type s_str struct {
	Data string
}

type s_bytes struct {
	Data  []byte
	Data2 [2]byte
}

func Test_hex_data(t *testing.T) {
	v1 := GetRDCInt64()
	hx1 := hex.EncodeToString(Int64ToBytes(v1))
	be1 := base64.StdEncoding.EncodeToString(Int64ToBytes(v1))
	v2, _ := hex.DecodeString(hx1)
	t.Log(v1, hx1, be1, bytes.Equal(Int64ToBytes(v1), v2))

	vv1 := GetRDCBytes(2048)
	hx2 := hex.EncodeToString(vv1)
	vv2, _ := hex.DecodeString(hx2)
	t.Log(len(vv1), len(hx2), bytes.Equal(vv1, vv2))

	sv1 := s_bytes{Data: vv1}
	sv2 := s_str{hx2}

	re, _ := json.Marshal(&sv1)
	re2, _ := json.Marshal(&sv2)

	t.Log(len(re), len(re2)) //, string(re), string(re2))
}

func Test_json_data(t *testing.T) {

	type s_str2 struct {
		Data []string
	}

	sv2 := s_str2{[]string{"1", "2"}}

	//re, _ := json.Marshal(&sv1)
	re := ""
	re2, _ := json.Marshal(&sv2)

	t.Log(len(re), len(re2))
	t.Log((re), string(re2))
}

func Test_compress(t *testing.T) {
	const data = `<?xml version="1.0"?>
<book>
	<meta name="title" content="The Go Programming Language"/>
	<meta name="authors" content="Alan Donovan and Brian Kernighan"/>
	<meta name="published" content="2015-10-26"/>
	<meta name="isbn" content="978-0134190440"/>
	<data>...</data>
</book>123456
`

	var b bytes.Buffer

	rt := timerm.CreateRecoTime()
	// Compress the data using the specially crafted dictionary.
	t.Log(rt.Run())
	zw, err := flate.NewWriter(&b, 1)

	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(zw, strings.NewReader(data)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	t.Log(rt.Run())
	t.Log(b.Bytes())
	t.Log(b.String())
	t.Log(len(data), b.Len(), float32(b.Len())/float32(len(data)))

	zr := flate.NewReader(bytes.NewReader(b.Bytes()))
	bwrbuf := make([]byte, 0)
	bwr := bytes.NewBuffer(bwrbuf)

	if _, err := io.Copy(bwr, zr); err != nil {
		t.Fatal(err)
	}
	if !(strings.EqualFold(data, bwr.String())) {
		t.Fatal("error.")
	}
	t.Log(bwr.String())
	if err := zr.Close(); err != nil {
		t.Fatal(err)
	}
}

func Test_compress_un(t *testing.T) {
	ap1 := createAesPack(getrandString(32))
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

	pd := CreatePackOperGen([]byte(data1))

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

func Test_compress_un2(t *testing.T) {
	ap1 := createAesPack(getrandString(32))
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

	pd := CreatePackOperGen([]byte(data1))

	re, _ := json.Marshal(pd)

	_ = re

	data := []byte(re)

	rt := timerm.CreateRecoTime()
	t.Log(rt.Run())
	cdata := ap1.Packing(data)
	t.Log(rt.Run())
	t.Log("compress:", float32(len(cdata))/float32(len(data)))

	undata, _ := ap1.Unpack(cdata)
	if bytes.Equal(data, undata) {
		t.Log("ok.")
	} else {
		t.Log("Error.")
	}

}

func Test_pack_type_size(t *testing.T) {
	ap1 := createAesPack(getrandString(32))
	t.Log(len(ap1.Packing([]byte{})))
	t.Log(len(ap1.ChangeCryKey()))
	t.Log(len(ap1.IsTheVersionConsistent()))
}

func Test_GsPack(t *testing.T) {

	fbuf := GetRDCBytes(1024 * 1024)

	a1 := NewGsPack("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")

	tmp := a1.Packing(fbuf)
	t.Log(tmp[len(tmp)-3 : len(tmp)])
	outbuf, _ := a1.Unpack(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		t.Log("ok.")
		t.Log(getsha1(fbuf))
	} else {
		t.Error()
	}

}
