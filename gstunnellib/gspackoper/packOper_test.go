package gspackoper

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"google.golang.org/protobuf/proto"

	//. "gstunnellib"

	"strings"
	"testing"
	"unsafe"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
	//"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func Test_JPackandun(t *testing.T) {
	t.Log("....................")

	v1 := []byte{1, 2, 3}
	j1 := NewPackPOGen(v1)
	t.Log(j1)
	t.Log(string(j1))
	dp1, _ := jsonUnpack(j1)
	t.Log(dp1)
	t.Log(string(v1) == string(dp1))
	if string(v1) == string(dp1) {

	} else {
		t.Error()
	}
}

/*
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
*/
func Test_GetSha256Hex(t *testing.T) {
	data := []byte{}

	buf := make([]byte, 100)

	for i := range [100]int{} {
		_ = i
		blen := binary.PutVarint(buf, gsrand.GetRDInt64())
		data = append(data, buf[:blen]...)
	}

	t.Log(gshash.GetSha256Hex(data))

}

func Test_GetSha256Hex_all(t *testing.T) {
	rawdata := gsrand.GetRDBytes(1024 * 1024)

	po1 := createPackOperGen(rawdata)
	vsha256 := string(po1.GetSha256())
	if string(po1.GetSha256_buf()) != vsha256 || string(po1.GetSha256_nobuf()) != vsha256 {
		t.Fatalf("GetSha256 is error.")
	}
	if string(po1.GetSha256_old()) != vsha256 || string(po1.GetSha256_pool()) != vsha256 {
		t.Fatalf("GetSha256 is error.")
	}
}
func Test_JsonPacking_OperChangeKey(t *testing.T) {
	b1, _ := JsonPacking_OperChangeKey()
	p1 := UnPack_Oper(b1)
	_ = p1
	t.Log("Ok.")
}

func Test_JsonPacking_OperGen(t *testing.T) {

	data := gsrand.GetRDBytes(10024)

	b1 := NewPackPOGen(data)
	p1 := UnPack_Oper(b1)
	_ = p1

	unb2, _ := jsonUnPack_OperGen(b1)
	h1 := gshash.GetSha256Hex(unb2)
	if gshash.GetSha256Hex(data) == h1 {
		t.Log("Ok.")
	} else {
		t.Error("Error.")
	}
	_ = p1

}

func Test_JPackandUnpack_oper(t *testing.T) {
	t.Log("....................")

	v1 := gsrand.GetRDBytes(953512)
	j1 := NewPackPOGen(v1)
	t.Log(j1)
	t.Log(string(j1))
	dp1, _ := jsonUnpack(j1)
	t.Log(dp1)
	t.Log(gshash.GetSha256Hex(v1) == gshash.GetSha256Hex(dp1))
	if gshash.GetSha256Hex(v1) == gshash.GetSha256Hex(dp1) {

	} else {
		t.Error()
	}
	_ = v1
}

func Test_bytesjoin(t *testing.T) {

	//const totalLoop int = 10000 * 100

	po1, _ := createPackOperChangeKey()

	po1.Data = gsrand.GetRDBytes(1024)

	b1 := po1.GetSha256_old()
	b2 := po1.GetSha256()
	b3 := po1.GetSha256_buf()
	b4 := po1.GetSha256_pool()

	if !bytes.Equal(b1[:], b2) {
		t.Error("error.")
	}
	if !bytes.Equal(b2, b3) {
		t.Fatal()
	}
	if !bytes.Equal(b3, b4) {
		t.Fatal()
	}

	if strings.Compare(hex.EncodeToString(b1[:]), hex.EncodeToString(b2[:])) != 0 {
		t.Error("error.")
	}

}

func Test_po_size(t *testing.T) {

	//const totalLoop int = 10000 * 100

	po2, _ := createPackOperChangeKey()

	_ = po2

	pd := createPackOperGen([]byte{})

	re, _ := json.Marshal(pd)

	_ = re

	//	pd2 := createPackOperGen_po1([]byte{})

	//	re2, _ := json.Marshal(pd2)

	t.Log(string(re), len(re))
	//	t.Log(string(re2), len(re2))
	t.Log("po size:", unsafe.Sizeof(po2))
	//	t.Log("po size:", unsafe.Sizeof(po2))

}

func Test_po_is(t *testing.T) {
	poData := NewPackPOGen(nil)
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

/*
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
*/
func Test_json__(t *testing.T) {
	pvdata := NewPackPOVersion()
	if !IsPOVersion(pvdata) {
		t.Fatal()
	}

	pov := createPackOperVersion()
	if pov.IsOk() != nil {
		t.Fatal()
	}

	pdata, _ := JsonPacking_OperChangeKey()
	if !IsChangeCryKey(pdata) {
		t.Fatal()
	}
}

func Test_GetEncryKey(t *testing.T) {
	key := gsrand.GetRDBytes(32)

	keyPD := PackOper{packOperPro: packOperPro{
		OperType: POChangeCryKey,
		OperData: []byte(key),
		Rand:     gsrand.GetRDBytes(8),
	}}

	keyPD.HashHex = keyPD.GetSha256()

	keyData, err := proto.Marshal(&keyPD)
	if err != nil {
		t.Error(err.Error())
	}

	key1 := getEncryKey(keyData)
	if !bytes.Equal(key, key1) {
		t.Error("key != key1.")
	}
	if !bytes.Equal(key, keyPD.GetEncryKey()) {
		t.Error("key != key1.")
	}
}

func Test_sha256_file(t *testing.T) {

	fp1 := "./tmp.hash.txt"
	f, err := os.Create(fp1)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove(fp1)
		if err != nil {
			t.Fatal(err.Error())
		}
	}()

	f.Write([]byte("1234567890abcdefghijklmnopqrst"))
	f.Close()

	t1 := time.Now()
	f, err = os.Open(fp1)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()
	s2 := sha256.New()
	_, err = io.Copy(s2, f)
	if err != nil {
		t.Fatal(err.Error())
	}
	t2 := time.Since(t1)
	fmt.Println("time milSec:", t2.Milliseconds())

	hash1 := hex.EncodeToString(s2.Sum([]byte{}))
	hash2 := "7e09b21b83692de4d3146479188d9b703e630f531bc795dccaabc659915ad56c"

	fmt.Println(
		hash1,
	)
	if !strings.EqualFold(hash1, hash2) {
		t.Fatal("Error.")
	}
}

// old len:  55  80
func Test_randSize(t *testing.T) {
	for i := 0; i < 100; i++ {
		p1 := NewPackPOVersion()
		p2, _ := JsonPacking_OperChangeKey()
		p3 := NewPackPOGen(nil)
		fmt.Printf("len:  %d  %d  %d\n", len(p1), len(p2), len(p3))
	}

}

func Test_randSize_sort(t *testing.T) {
	list_p1 := []int{}
	list_p2 := []int{}
	list_p3 := []int{}

	for i := 0; i < 100; i++ {
		p1 := NewPackPOVersion()
		p2, _ := JsonPacking_OperChangeKey()
		p3 := NewPackPOGen(nil)
		fmt.Printf("len:  %d  %d  %d\n", len(p1), len(p2), len(p3))

		list_p1 = append(list_p1, len(p1))
		list_p2 = append(list_p2, len(p2))
		list_p3 = append(list_p3, len(p3))
	}

	sort.Ints(list_p1)
	sort.Ints(list_p2)
	sort.Ints(list_p3)

	fmt.Printf("len:  %v  \n", list_p1)
	fmt.Printf("len:  %v  \n", list_p2)
	fmt.Printf("len:  %v  \n", list_p3)
}

func Test_operType_bin(t *testing.T) {
	fmt.Printf("operType bin: %b  %b  %b  %b  %b\n",
		POBegin,
		POGenOper,
		POChangeCryKey,
		POVersion,
		POEnd)
}

func in_Rsa_RDdata_test(pri *rsa.PrivateKey, pub *rsa.PublicKey) {
	text := gsrand.GetRDBytes(300)

	//pri, pub := GenerateKeyPair(4096)

	ciphertext := gsrsa.EncryptWithPublicKey(text, pub)

	dedata := gsrsa.DecryptWithPrivateKey(ciphertext, pri)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func Test_clientPubKey(t *testing.T) {

	serverpri, serverpub := gsrsa.GenerateKeyPair(4096)

	clientpri, clientpub := gsrsa.GenerateKeyPair(4096)

	_ = clientpri

	clientpubkeybytes := gsrsa.PublicKeyToBytes(clientpub)

	ciphertext := gsrsa.EncryptWithPublicKeyAnyLen(clientpubkeybytes, serverpub)

	packdata := NewPackPOClientPubKey(ciphertext)

	pd := UnPack_Oper(packdata)

	clientpubkeydata := pd.GetClientPubKeyData()

	if pd.IsClientPubKey() {
		if !bytes.Equal(ciphertext, clientpubkeydata) {
			panic("error")
		}
	}

	clientpubkeydedata := gsrsa.DecryptWithPrivateKeyAnyLen(clientpubkeydata, serverpri)
	//in_Rsa_RDdata_test(clientpri, get)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(clientpubkeybytes, clientpubkeydedata) {
		panic("!bytes.Equal(text, dedata)")
	}

	clientpubkey2 := gsrsa.PublicKeyFromBytes(clientpubkeydedata)

	in_Rsa_RDdata_test(clientpri, clientpubkey2)
}
