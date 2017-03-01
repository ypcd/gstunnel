/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
*/
package gstunnellib

import (
    "crypto/sha1"
    "fmt"
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
}

func getrand() int {
    rd := rand.New(
        rand.NewSource(time.Now().UnixNano()))
    return rd.Int()
}

func Test_JPackandun(t *testing.T) {
    p("....................")

    v1 := []byte{1,2,3}
    j1 := jsonPacking(v1)
    p( j1)
    p(string(j1))
    dp1 := jsonUnpack(j1)
    p(dp1)
    p( string(v1)==string(dp1))
    if string(v1)==string(dp1) {
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
        if n==0 {
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

    v1 := []byte{1,2,3}
    j1 := Packing(v1)
    p( j1)
    p(string(j1))
    dp1 := Unpack(j1)
    p(dp1)
    p( string(v1)==string(dp1))
    if string(v1)==string(dp1) {
        t.Log("ok.")
    } else {
        t.Error()
    }
}

var unpbuf []byte
func Unpackbuf(data []byte) []byte {
    unpbuf = append(unpbuf, data...)
    var outbuf, buf []byte
    
    for ix, re := find0(unpbuf); re; ix, re = find0(unpbuf) {
    if !re { return outbuf }
    
    buf = unpbuf[:ix+1]
    unpbuf = unpbuf[ix+1:]
    outbuf = append(outbuf, Unpack(buf)...)
    }
    return outbuf
}


type ftype func()
func mtF(frun ftype) {
    for i:=0; i<10*5; i++ {
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
        if n==0 {break}
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
        if n==0 {break}
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
        if n==0 {break}
        fbuf = append(fbuf, buf[:n]...)
    }
    b.StartTimer()
    var tmp, outbuf []byte
    for i:=0; i < b.N; i++ {
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
        if n==0 {break}
        fbuf = append(fbuf, buf[:n]...)
    }
    b.StartTimer()
    a1 := CreateAes("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")
    var tmp, outbuf []byte
    for i:=0; i < b.N; i++ {
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
        if n==0 {break}
        fbuf = append(fbuf, buf[:n]...)
    }

    a1 := CreateAesPack("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")

    tmp := a1.Packing(fbuf)
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
        if n==0 {break}
        fbuf = append(fbuf, buf[:n]...)
    }

    a1 := CreateAesPack("5Wl)hPO9~UF_IecIN$e#uW!xc%7Yo$iQ")
    tmp := a1.Packing(fbuf)
    outbuf := a1.Unpack(tmp)

    if getsha1(fbuf) == getsha1(outbuf) {
        p("ok.", getrand())
    } else {
        p("Error")
    }
    wlist.Done()
}

func Test_MtAestpack(t *testing.T) {
    mtF(aestpack)
}

