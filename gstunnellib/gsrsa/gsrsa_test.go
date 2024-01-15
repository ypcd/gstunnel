package gsrsa

import (
	"bytes"
	randc "crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

func Test_Rsa(t *testing.T) {
	text := []byte("abcdef123456")

	pri, pub := GenerateKeyPair(4096)

	ciphertext := EncryptWithPublicKey(text, pub)

	dedata := DecryptWithPrivateKey(ciphertext, pri)

	log.Println("data:", dedata)
	log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func noUseTest_Rsa_RDdata_maxLen(t *testing.T) {
	pri, pub := GenerateKeyPair(4096)

	for i := 370; i < 10000; i++ {
		text := make([]byte, i)
		_, err := randc.Read(text)
		gserror.CheckError_panic(err)

		//pri, pub := GenerateKeyPair(4096)

		ciphertext := EncryptWithPublicKey(text, pub)

		dedata := DecryptWithPrivateKey(ciphertext, pri)

		//log.Println("data:", dedata)
		//log.Println("data:", string(dedata))

		if !bytes.Equal(text, dedata) {
			panic("!bytes.Equal(text, dedata)")
		}
	}
}

func in_Rsa_RDdata_test(pri *rsa.PrivateKey, pub *rsa.PublicKey) {
	text := make([]byte, 300)
	_, err := randc.Read(text)
	gserror.CheckError_panic(err)

	//pri, pub := GenerateKeyPair(4096)

	ciphertext := EncryptWithPublicKey(text, pub)

	dedata := DecryptWithPrivateKey(ciphertext, pri)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func Test_Rsa_RDdata(t *testing.T) {
	pri, pub := GenerateKeyPair(4096)
	in_Rsa_RDdata_test(pri, pub)
}

func Test_key_bytes(t *testing.T) {
	pri, pub := GenerateKeyPair(4096)

	pribytes := PrivateKeyToBytes(pri)

	pubbytes := PublicKeyToBytes(pub)

	pri2 := PrivateKeyFromBytes(pribytes)
	pub2 := PublicKeyFromBytes(pubbytes)

	in_Rsa_RDdata_test(pri, pub)
	in_Rsa_RDdata_test(pri2, pub2)

	in_Rsa_RDdata_test(pri, pub2)
	in_Rsa_RDdata_test(pri2, pub)

}

func noUseTest_print_key_bytes_base64(t *testing.T) {
	pri, pub := GenerateKeyPair(4096)

	pribytes := PrivateKeyToBytes(pri)

	pubbytes := PublicKeyToBytes(pub)

	fmt.Println(string(pribytes))
	fmt.Println(string(pubbytes))

	priendata := PrivateKeyToBase64(pri)

	pubendata := PublicKeyToBase64(pub)

	fmt.Println(string(priendata))
	fmt.Println(string(pubendata))

}
func Test_key_base64(t *testing.T) {
	text := make([]byte, 300)
	_, err := randc.Read(text)
	gserror.CheckError_panic(err)

	pri, pub := GenerateKeyPair(4096)
	_ = pub
	priendata := PrivateKeyToBase64(pri)
	//println("rdata:", string(priendata))

	pri2 := PrivateKeyFromBase64(priendata)

	pubendata := PublicKeyToBase64(pub)

	pub2 := PublicKeyFromBase64(pubendata)

	in_Rsa_RDdata_test(pri, pub)
	in_Rsa_RDdata_test(pri2, pub2)

	in_Rsa_RDdata_test(pri, pub2)
	in_Rsa_RDdata_test(pri2, pub)
}
func noTest_key_json(t *testing.T) {
	//json.Marshal()
}

func in_Rsa_RDdata_test_anyLen(nbytes int, pri *rsa.PrivateKey, pub *rsa.PublicKey) {

	text := make([]byte, nbytes)
	_, err := randc.Read(text)
	gserror.CheckError_panic(err)

	//pri, pub := GenerateKeyPair(4096)

	ciphertext := EncryptWithPublicKeyAnyLen(text, pub)

	dedata := DecryptWithPrivateKeyAnyLen(ciphertext, pri)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func Test_Rsa_RDdata_test_anyLen(t *testing.T) {
	pri, pub := GenerateKeyPair(4096)

	for i := 1000; i < 1020; i++ {
		in_Rsa_RDdata_test_anyLen(i, pri, pub)
	}
}

func Test_Rsa_RDdata_test_anyLen_json(t *testing.T) {

	pri, pub := GenerateKeyPair(4096)

	text := make([]byte, 5000)
	_, err := randc.Read(text)
	gserror.CheckError_panic(err)

	text = PublicKeyToBytes(pub)

	ciphertext := encryptWithPublicKeyAnyLen_old(text, pub)

	jendata, err := json.Marshal(ciphertext)
	gserror.CheckError_panic(err)

	rsadata2 := &rsaEncryptData{}
	err = json.Unmarshal(jendata, rsadata2)
	gserror.CheckError_panic(err)

	dedata := decryptWithPrivateKeyAnyLen_old(rsadata2, pri)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func Test_Rsa_RDdata_test_anyLen_json2(t *testing.T) {

	pri, pub := GenerateKeyPair(4096)

	text := PublicKeyToBytes(pub)

	ciphertext := EncryptWithPublicKeyAnyLen(text, pub)

	dedata := DecryptWithPrivateKeyAnyLen(ciphertext, pri)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func inTest_genrsa_time(bits int) {
	t1 := time.Now()

	loopNum := 2
	for i := 0; i < loopNum; i++ {
		key1, key2 := GenerateKeyPair(bits)
		_, _ = key1, key2
	}

	t2 := time.Now()

	fmt.Printf("bits:%d total:%fs avg:%fs\n", bits, t2.Sub(t1).Seconds(), t2.Sub(t1).Seconds()/float64(loopNum))
}

/*
bits:512 total:0.006582 avg:0.003291
bits:1024 total:0.027462 avg:0.013731
bits:2048 total:0.171080 avg:0.085540
bits:4096 total:5.923486 avg:2.961743
bits:8192 total:49.676865 avg:24.838432

bits:512 total:0.009366s avg:0.004683s
bits:1024 total:0.040924s avg:0.020462s
bits:2048 total:0.375675s avg:0.187837s
bits:4096 total:2.904750s avg:1.452375s
*/
func noTest_genrsa_time(t *testing.T) {
	for i := 1.0; i <= 4; i++ {
		inTest_genrsa_time(int(256 * math.Pow(2, i)))
	}
}
