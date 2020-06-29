/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package gstunnellib

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	//"math"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var testi1 = 1
var pf = fmt.Printf
var wlist sync.WaitGroup

func init() {
	p = Nullprint
	p = fmt.Println
	debug_tag = true
}

func getrand() int {
	rd := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	return rd.Int()
}

func Test_JPackandun(t *testing.T) {
	p("....................")

	v1 := []byte{1, 2, 3}
	j1 := jsonPacking(v1)
	p(j1)
	p(string(j1))
	dp1 := jsonUnpack(j1)
	p(dp1)
	p(string(v1) == string(dp1))
	if string(v1) == string(dp1) {
		t.Log("ok.")
	} else {
		t.Error()
	}
}

func Test_AesStr(t *testing.T) {

	v11 := "1234567890"
	v11 = v11
	v1 := []byte(v11)
	p(string(v1))
	e1 := encrypter(v1)
	p(e1)
	d1 := decrypter(e1)
	p(string(d1))
	if string(v1) == string(d1) {
		t.Log("ok.")
	} else {
		t.Error()
	}
}

func filesha1(fp string) string {
	f, _ := os.Open(fp)
	buf := make([]byte, 1024*1024)
	h := sha1.New()
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			return fmt.Sprintf("%X", h.Sum(nil))
		} else {
			buf2 := buf[:n]
			h.Write(buf2)
		}
	}
}

func getsha1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%X", h.Sum(nil))
}

func Test_Filesha1(t *testing.T) {
	fp2 := `testaes.data`
	p(filesha1(fp2))
	if filesha1(fp2) == "3C8243734CF43DD7BB2332BA05B58CCACFA4377C" {
		t.Log("ok.")
	} else {
		t.Error()
	}
}

func Test_Aes2(t *testing.T) {

	fp2 := `testaes.data`
	fp3 := "./aestmp.data"
	fp4 := "./aes.data"
	f, _ := os.Open(fp2)
	outf, _ := os.Create(fp3)

	buf := make([]byte, 1024*1024*8)

	for {
		n, _ := f.Read(buf)
		buf = buf[:n]
		if 0 == n {
			f.Close()
			outf.Close()
			break
		} else {
			outf.Write(encrypter(buf))
		}
	}

	p("...................................")
	buf = make([]byte, 1024*1024*8)
	f2, _ := os.Open(fp3)

	outf2, _ := os.Create(fp4)
	p("...")
	for {
		n, _ := f2.Read(buf)
		buf = buf[:n]
		if 0 == n {
			p(f2.Close())
			p(outf2.Close())
			break
		} else {
			ebuf := decrypter(buf)
			outf2.Write(ebuf)
		}
	}
	if filesha1(fp2) == filesha1(fp4) {
		t.Log("ok.")
		p(os.RemoveAll(fp3))
		p(os.RemoveAll(fp4))
	} else {
		t.Error()
	}

}

func Test_Packandun(t *testing.T) {

	p("....................")

	v1 := []byte{1, 2, 3}
	j1 := Packing(v1)
	p(j1)
	p(string(j1))
	dp1 := Unpack(j1)
	p(dp1)
	p(string(v1) == string(dp1))
	if string(v1) == string(dp1) {
		t.Log("ok.")
	} else {
		t.Error()
	}
}

var unpbuf []byte

func Test_find0(t *testing.T) {
	data := make([]byte, 0)
	buf := make([]byte, 100)

	for i := 0; i < 1000; i++ {
		blen := binary.PutVarint(buf, rand_s.Int63())
		data = append(data, buf[:blen]...)
	}
	data[rand.Intn(len(data))] = 0
	ix, _ := Find0(data)
	ix2 := bytes.IndexByte(data, 0)
	fmt.Println(ix, ix2)
	if ix == ix2 {
		t.Log("ok.")
	} else {
		t.Error("error.")
	}
}

func Unpackbuf(data []byte) []byte {
	unpbuf = append(unpbuf, data...)
	var outbuf, buf []byte

	for ix, re := Find0(unpbuf); re; ix, re = Find0(unpbuf) {
		if !re {
			return outbuf
		}

		buf = unpbuf[:ix+1]
		unpbuf = unpbuf[ix+1:]
		outbuf = append(outbuf, Unpack(buf)...)
	}
	return outbuf
}

type ftype func()

func mtF(frun ftype) {
	for i := 0; i < 10*5; i++ {
		go frun()
		wlist.Add(1)
	}
	p("go")
	wlist.Wait()
}

func Test_Aest(t *testing.T) {

	fp2 := `testaes2.data`
	f, _ := os.Open(fp2)

	buf := make([]byte, 1024*128)
	var fbuf []byte
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		fbuf = append(fbuf, buf[:n]...)
	}

	a1 := CreateAes("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")

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

	fp2 := `testaes4.data`
	f, _ := os.Open(fp2)

	buf := make([]byte, 1024*128)
	var fbuf []byte
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		fbuf = append(fbuf, buf[:n]...)
	}

	a1 := CreateAes("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")
	tmp := a1.encrypter(fbuf)
	outbuf := a1.decrypter(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		p("ok.", getrand())

	} else {
		p("Error")
	}
	wlist.Done()
}

func Test_MtAest(t *testing.T) {
	mtF(aest)
}

func Benchmark_Aes(b *testing.B) {
	b.StopTimer()
	fp2 := `testaes4.data`
	f, _ := os.Open(fp2)

	buf := make([]byte, 1024*128)
	var fbuf []byte
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		fbuf = append(fbuf, buf[:n]...)
	}
	b.StartTimer()
	var tmp, outbuf []byte
	for i := 0; i < b.N; i++ {
		tmp = encrypter(fbuf)
		outbuf = decrypter(tmp)
	}
	b.StopTimer()
	tmp = outbuf[:]

}

func Benchmark_Aest(b *testing.B) {
	b.StopTimer()

	fp2 := `testaes4.data`
	f, _ := os.Open(fp2)

	buf := make([]byte, 1024*128)
	var fbuf []byte
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		fbuf = append(fbuf, buf[:n]...)
	}
	b.StartTimer()
	a1 := CreateAes("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")
	var tmp, outbuf []byte
	for i := 0; i < b.N; i++ {
		tmp = a1.encrypter(fbuf)
		outbuf = a1.decrypter(tmp)
	}
	b.StopTimer()
	tmp = outbuf[:]
}

func Test_Aestpack(t *testing.T) {

	fp2 := `testaes4.data`
	f, _ := os.Open(fp2)

	buf := make([]byte, 1024*128)
	var fbuf []byte
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		fbuf = append(fbuf, buf[:n]...)
	}

	a1 := CreateAesPack("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")

	tmp := a1.Packing(fbuf)
	p(tmp[len(tmp)-3 : len(tmp)])
	outbuf := a1.Unpack(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		t.Log("ok.")
		t.Log(getsha1(fbuf))
	} else {
		t.Error()
	}

}

func aestpack() {

	fp2 := `testaes4.data`
	f, _ := os.Open(fp2)

	buf := make([]byte, 1024*128)
	var fbuf []byte
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		fbuf = append(fbuf, buf[:n]...)
	}

	a1 := CreateAesPack("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")
	tmp := a1.Packing(fbuf)
	outbuf := a1.Unpack(tmp)

	if getsha1(fbuf) == getsha1(outbuf) {
		p("ok.", getrand())
	} else {
		//t.Error()
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
		p("d1 == d2, ok.")
	}
	if p1.Rand == p2.Rand {
		p("p1.Rand == p2.Rand, ok.")
	}
	fmt.Println(d2)
}

func Test_jsonPackUnPack_Rand(t *testing.T) {
	d1 := []byte("123Abc")
	p1 := PackRand{Pack{d1}, rand.Int63()}
	j1 := jsonPacking_Rand(d1)
	p(p1)
	p(j1)

	p2 := PackRand{}
	d2 := jsonUnPack_Rand(j1)
	p(p2)
	//d2 := p2.Data
	if string(d1) == string(d2) {
		p("d1 == d2, ok.")
	} else {
		t.Error()
	}
	/*
	   if p1.Rand == p2.Rand{
	       p("p1.Rand == p2.Rand, ok.")
	   } else{ t.Error()}
	*/
	fmt.Println(d2)
}

func Test_GetSha256Hex(t *testing.T) {
	data := []byte{}

	buf := make([]byte, 100)

	for i := range [100]int{} {
		_ = i
		blen := binary.PutVarint(buf, rand_s.Int63())
		data = append(data, buf[:blen]...)
	}

	p(GetSha256Hex(data))

}

func Test_json_RandHash(t *testing.T) {
	data := []byte{}

	buf := make([]byte, 100)

	for i := range [100]int{} {
		_ = i
		blen := binary.PutVarint(buf, rand_s.Int63())
		data = append(data, buf[:blen]...)
	}

	j1 := jsonPacking_RandHash(data)
	d2 := jsonUnPack_RandHash(j1)

	if GetSha256Hex(data) == GetSha256Hex(d2) {
		t.Log("OK.")
	} else {
		t.Error("!=")
	}

}

func Test_binConv(t *testing.T) {
	buf := make([]byte, 0)
	p1 := &buf
	p(p1)
	bbuf := bytes.NewBuffer(buf)
	err := binary.Write(bbuf, binary.LittleEndian, int32(123))
	p(err)
	p(&buf, buf)
}

func Test_IsChangeCryKey(t *testing.T) {
	pgen := CreatePackOperGen([]byte{})
	pckey := CreatePackOperChangeKey()
	j1, _ := json.Marshal(pgen)
	j2, _ := json.Marshal(pckey)
	if IsChangeCryKey(append(j1, 0)) == false && IsChangeCryKey(append(j2, 0)) == true {
		t.Log("OK.")
	} else {
		t.Error("error.")
	}
}

func Test_GetRDF64(t *testing.T) {

	f1 := GetRDF64()

	//var rd_s rand.Source = rand.NewSource(time.Now().Unix())

	for i := 0; i < 8; i++ {
		//p(rd_s.Int63())
		//p(float64(rd_s.Int63()) / math.Pow(2, 64))
		//p(GetRDF64())
		p(GetRDInt8())
	}
	p(GetRDBytes(32))
	if len(GetRDBytes(32)) == 32 && f1 < 1 {
		t.Log("OK.")
	} else {
		t.Error("error.")
	}
}

func Test_IntsToBytes(t *testing.T) {
	p(Int32ToBytes(GetRDInt32()))
	p(Int64ToBytes(GetRDInt64()))
	t.Log("ok.")
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

	if GetSha256Hex(data) == GetSha256Hex(jsonUnPack_OperGen(b1)) {
		t.Log("Ok.")
	} else {
		t.Error("Error.")
	}
	_ = p1
}

func Test_Aespack_changeCryKey(t *testing.T) {
	key := GetRDBytes(32)
	ap1 := CreateAesPack(string(key))
	cp1 := ap1.a
	key2 := GetRDBytes(32)
	ap1.changeCryKey(key2)
	cp2 := ap1.a

	p(key)
	p(key2)
	p(cp1, cp2)
	if cp1 != cp2 {
		t.Log("Ok.")
	} else {
		t.Error("Error.")
	}
	_ = key
}

func Test_Aespack_changeCryKey2(t *testing.T) {
	key := GetRDBytes(32)
	ap1 := CreateAesPack(string(key))
	cp1 := ap1.a

	ap1.ChangeCryKey()
	cp2 := ap1.a

	p(key)
	//	p(key2)
	p(cp1, cp2)
	if cp1 != cp2 {
		t.Log("Ok.")
	} else {
		t.Error("Error.")
	}
	_ = key
}

func Test_JPackandUnpack_oper(t *testing.T) {
	p("....................")

	v1 := GetRDBytes(953512)
	j1 := jsonPacking(v1)
	p(j1)
	p(string(j1))
	dp1 := jsonUnpack(j1)
	p(dp1)
	p(GetSha256Hex(v1) == GetSha256Hex(dp1))
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

	ap1 := CreateAesPack(string(key1))
	ap2 := CreateAesPack(string(key2))

	ed1 := ap1.a.encrypter(d1)

	dd1, dd2 := ap1.a.decrypter(ed1), ap2.a.decrypter(ed1)

	p(dd1)
	p("------------------------------------------------------------------------------------------------------------------------------")
	p(dd2)
}

func Test_wbuf(t *testing.T) {
	wbuf := []byte{1, 2, 3, 4, 5, 6}
	p(wbuf)
	wbuf = wbuf[len(wbuf):]
	p(wbuf)
}
