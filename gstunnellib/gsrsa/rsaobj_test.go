package gsrsa

import (
	"bytes"
	randc "crypto/rand"
	"sync"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

func in_obj_Rsa_RDdata_test(r *RSA) {
	text := make([]byte, 300)
	_, err := randc.Read(text)
	gserror.CheckError_panic(err)

	//Pri, Pub := GenerateKeyPair(4096)

	ciphertext := r.EncryptWithPublicKeyAnyLen(text)

	dedata := r.DecryptWithPrivateKeyAnyLen(ciphertext)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func in_obj_Rsa_RDdata_test_anyLen(nbytes int, r *RSA) {

	text := make([]byte, nbytes)
	_, err := randc.Read(text)
	gserror.CheckError_panic(err)

	//Pri, Pub := GenerateKeyPair(4096)

	ciphertext := r.EncryptWithPublicKeyAnyLen(text)

	dedata := r.DecryptWithPrivateKeyAnyLen(ciphertext)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func Test_rsaobj1(t *testing.T) {
	rsa1 := NewGenRSAObj(4096)

	in_obj_Rsa_RDdata_test(rsa1)
	in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)
}

func Test_rsaobj2(t *testing.T) {
	Pri, _ := GenerateKeyPair(4096)
	rsa1 := NewRSAObj(Pri)

	in_obj_Rsa_RDdata_test(rsa1)
	in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)
}

func Test_rsaobj3(t *testing.T) {
	Pri, Pub := GenerateKeyPair(4096)
	r := NewRSAObjPub(Pub)

	//in_obj_Rsa_RDdata_test(rsa1)

	text := make([]byte, 300)
	_, err := randc.Read(text)
	gserror.CheckError_panic(err)

	//Pri, Pub := GenerateKeyPair(4096)

	ciphertext := r.EncryptWithPublicKeyAnyLen(text)

	dedata := DecryptWithPrivateKeyAnyLen(ciphertext, Pri)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}

	//in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)

	text = make([]byte, 1001)
	_, err = randc.Read(text)
	gserror.CheckError_panic(err)

	//Pri, Pub := GenerateKeyPair(4096)

	ciphertext = r.EncryptWithPublicKeyAnyLen(text)

	dedata = DecryptWithPrivateKeyAnyLen(ciphertext, Pri)

	//log.Println("data:", dedata)
	//log.Println("data:", string(dedata))

	if !bytes.Equal(text, dedata) {
		panic("!bytes.Equal(text, dedata)")
	}
}

func Test_rsaobj4(t *testing.T) {
	//Pri, _ := GenerateKeyPair(4096)
	rsa2 := NewGenRSAObj(4096)
	rsa1 := NewRSAObjFromBase64(rsa2.PrivateKeyToBase64())

	in_obj_Rsa_RDdata_test(rsa1)
	in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)
}

func Test_rsaobj5(t *testing.T) {
	//Pri, _ := GenerateKeyPair(4096)
	rsa2 := NewGenRSAObj(4096)
	rsa1 := NewRSAObjFromBytes(rsa2.PrivateKeyToBytes())

	in_obj_Rsa_RDdata_test(rsa1)
	in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)

	n := rsa1.Pri.Size()
	_ = n
	n2 := rsa1.Pub.Size()
	_ = n2
}

func Test_rsaobj6(t *testing.T) {
	//Pri, Pub := GenerateKeyPair(4096)
	rsa2 := NewGenRSAObj(4096)
	rsa1 := &RSA{
		Pri: rsa2.Pri,
		Pub: PublicKeyFromBase64(rsa2.PublicKeyToBase64()),
	}

	in_obj_Rsa_RDdata_test(rsa1)
	in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)
}

func Test_rsaobj7(t *testing.T) {
	//Pri, Pub := GenerateKeyPair(4096)
	rsa2 := NewGenRSAObj(4096)
	rsa1 := &RSA{
		Pri: rsa2.Pri,
		Pub: PublicKeyFromBytes(rsa2.PublicKeyToBytes()),
	}

	in_obj_Rsa_RDdata_test(rsa1)
	in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)
}

func noTest_rsaobj_no_4096(t *testing.T) {
	//Pri, Pub := GenerateKeyPair(4096)
	rsa2 := NewGenRSAObj(1024)
	rsa1 := &RSA{
		Pri: rsa2.Pri,
		Pub: PublicKeyFromBytes(rsa2.PublicKeyToBytes()),
	}

	in_obj_Rsa_RDdata_test(rsa1)
	in_obj_Rsa_RDdata_test_anyLen(1001, rsa1)
}

func Test_rsaobj1_mt(t *testing.T) {
	rsa1 := NewGenRSAObj(4096)

	wg := sync.WaitGroup{}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				in_obj_Rsa_RDdata_test(rsa1)
			}
		}()
	}
	wg.Wait()
}

func Test_b1(t *testing.T) {
	r1 := NewGenRSAObj(4096)
	r2 := r1
	_ = r2
	r3 := r1.NewRSA()
	_ = r3
}

func Test_rsa_race1(t *testing.T) {
	pri, pub := GenerateKeyPair(g_rsaKeyLenbits)

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v := NewRSAObjPub(pub)
			v2 := NewRSAObj(pri)
			_, _ = v, v2
		}()
	}
	wg.Wait()
}

func Test_rsa_race2(t *testing.T) {
	_, pub := GenerateKeyPair(g_rsaKeyLenbits)

	rsa1 := NewRSAObjPub(pub)
	//rsa2 := NewRSAObj(pri)

	loopNum := 10000 * 1
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < loopNum; i++ {
			v := rsa1.NewRSA()
			_ = v
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < loopNum; i++ {
			v := NewRSAObjPub(rsa1.Pub)
			_ = v
		}
	}()

	wg.Wait()
}

func Test_rsa_race3(t *testing.T) {
	_, pub := GenerateKeyPair(g_rsaKeyLenbits)

	bpub := PublicKeyToBytes(pub)
	//rsa1 := NewRSAObjPub(pub)
	//rsa2 := NewRSAObj(pri)
	loopNum := 10000 * 10
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < loopNum; i++ {
			//v := PublicKeyToBytes(pub)
			v := []byte{}
			v = append(v, bpub...)
			_ = v
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < loopNum; i++ {
			v := PublicKeyFromBytes(bpub)
			_ = v
		}
	}()

	wg.Wait()
}
