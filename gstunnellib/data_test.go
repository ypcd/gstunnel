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

	//"math"
	. "gstunnellib/gsrand"
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

	a1 := createAes([]byte(GetrandString(32)))

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

	a1 := createAes([]byte(GetrandString(32)))
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

	a1 := createAesPack(GetrandString(32))

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

	a1 := createAesPack(GetrandString(32))
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

func Test_binConv(t *testing.T) {
	buf := make([]byte, 0)
	p1 := &buf
	t.Log(p1)
	bbuf := bytes.NewBuffer(buf)
	err := binary.Write(bbuf, binary.LittleEndian, int32(123))
	t.Log(err)
	t.Log(&buf, buf)
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

func Test_nullFunc(t *testing.T) {
	Nullprint(1, "a")
	Nullprintf("%s.\n", "test")
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

func Test_pack_type_size(t *testing.T) {
	ap1 := createAesPack(GetrandString(32))
	t.Log(len(ap1.Packing([]byte{})))
	t.Log(len(ap1.ChangeCryKey()))
	t.Log(len(ap1.IsTheVersionConsistent()))
}
