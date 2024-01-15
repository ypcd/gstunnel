package gsrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

// byte len
const g_oneLen int = 370
const g_rsaKeyLen int = 512      //4096 bit
const g_rsaKeyLenbits int = 4096 //4096 bit

//var g_logger
/*
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
*/

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	gserror.CheckError_panic(err)
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	gserror.CheckError_panic(err)

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// PrivateKeyFromBytes bytes to private key
func PrivateKeyFromBytes(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			gserror.CheckError_panic(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	gserror.CheckError_panic(err)
	return key
}

// PublicKeyFromBytes bytes to public key
func PublicKeyFromBytes(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			gserror.CheckError_panic(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	gserror.CheckError_panic(err)
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Panic("not ok")
	}
	return key
}

// The maximum byte length of msg is 382, 4096 key.
// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	gserror.CheckError_panic(err)
	return ciphertext
}

// The byte length of ciphertext is 512.
// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	gserror.CheckError_panic(err)
	return plaintext
}

type rsaEncryptData struct {
	Data [][]byte
}

// EncryptWithPublicKey encrypts data with public key
func encryptWithPublicKeyAnyLen_old(msg []byte, pub *rsa.PublicKey) *rsaEncryptData {

	if pub.Size() != g_rsaKeyLen {
		panic("r.pri.Size() != g_rsaKeyLen")
	}

	var datals [][]byte

	count := len(msg) / g_oneLen
	for i := 1; i <= count; i++ {
		datals = append(datals, EncryptWithPublicKey(msg[g_oneLen*(i-1):g_oneLen*i], pub))
	}
	if len(msg)%g_oneLen != 0 {
		datals = append(datals, EncryptWithPublicKey(msg[g_oneLen*count:], pub))
	}

	return &rsaEncryptData{datals}
}

// DecryptWithPrivateKey decrypts data with private key
func decryptWithPrivateKeyAnyLen_old(ciphertext *rsaEncryptData, priv *rsa.PrivateKey) []byte {

	if priv.Size() != g_rsaKeyLen {
		panic("r.pri.Size() != g_rsaKeyLen")
	}

	data := make([]byte, 0, g_oneLen*len(ciphertext.Data))
	ctext := ciphertext.Data

	for _, v := range ctext {
		data = append(data, DecryptWithPrivateKey(v, priv)...)
	}

	return data
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKeyAnyLen(msg []byte, pub *rsa.PublicKey) []byte {

	//datals := make([][]byte, 0, len(msg)*2)

	var datals []byte

	count := len(msg) / g_oneLen
	for i := 1; i <= count; i++ {
		datals = append(datals, EncryptWithPublicKey(msg[g_oneLen*(i-1):g_oneLen*i], pub)...)
	}
	if len(msg)%g_oneLen != 0 {
		datals = append(datals, EncryptWithPublicKey(msg[g_oneLen*count:], pub)...)
	}

	return datals
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKeyAnyLen(ciphertext []byte, priv *rsa.PrivateKey) []byte {

	count := len(ciphertext) / g_rsaKeyLen
	data := make([]byte, 0, g_oneLen*count)

	for i := 1; i <= count; i++ {
		data = append(data, DecryptWithPrivateKey(ciphertext[g_rsaKeyLen*(i-1):g_rsaKeyLen*i], priv)...)
	}

	return data
}

func PrivateKeyToBase64(priv *rsa.PrivateKey) []byte {
	src := PrivateKeyToBytes(priv)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func PublicKeyToBase64(pub *rsa.PublicKey) []byte {
	src := PublicKeyToBytes(pub)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func PrivateKeyFromBase64(src []byte) *rsa.PrivateKey {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(dst, src)

	return PrivateKeyFromBytes(dst)
}

func PublicKeyFromBase64(src []byte) *rsa.PublicKey {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(dst, src)

	return PublicKeyFromBytes(dst)
}
