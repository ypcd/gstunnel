package gsrsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
)

type RSA struct {
	Pri *rsa.PrivateKey
	Pub *rsa.PublicKey

	lock_pri, lock_pub sync.Mutex
}

func NewRSAObj(Pri *rsa.PrivateKey) *RSA {
	r := new(RSA)
	r.lock_pri.Lock()
	defer r.lock_pri.Unlock()
	r.lock_pub.Lock()
	defer r.lock_pub.Unlock()
	pri2 := PrivateKeyFromBytes(PrivateKeyToBytes(Pri))
	r.Pri = pri2
	r.Pub = &pri2.PublicKey
	return r
}

func NewRSAObjPub(Pub *rsa.PublicKey) *RSA {
	r := new(RSA)
	//r.lock_pri.Lock()
	//defer r.lock_pri.Unlock()
	r.lock_pub.Lock()
	defer r.lock_pub.Unlock()
	//PublicKeyFromBytes(PublicKeyToBytes(Pub))
	r.Pub = PublicKeyFromBytes(PublicKeyToBytes(Pub))
	return r
}

func NewRSAObjFromBytes(keydata []byte) *RSA {
	Pri := PrivateKeyFromBytes(keydata)
	return &RSA{Pri: Pri, Pub: &Pri.PublicKey}
}

func NewRSAObjFromBase64(priKeyBase64 []byte) *RSA {
	Pri := PrivateKeyFromBase64(priKeyBase64)
	return &RSA{Pri: Pri, Pub: &Pri.PublicKey}
}

func NewRSAObjFromPubKeyBase64(pubKeyBase64 []byte) *RSA {
	return &RSA{Pub: PublicKeyFromBase64(pubKeyBase64)}
}

func NewGenRSAObj(bits int) *RSA {
	Pri, Pub := GenerateKeyPair(bits)
	return &RSA{Pri: Pri, Pub: Pub}
}

func (r *RSA) SetPubkey(Pub *rsa.PublicKey) {
	r.Pub = Pub
}

// PrivateKeyToBytes private key to bytes
func (r *RSA) PrivateKeyToBytes() []byte {
	if r.Pri == nil {
		panic("value is nil")
	}
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(r.Pri),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func (r *RSA) PublicKeyToBytes() []byte {
	//r.lock_pub.Lock()
	//defer r.lock_pub.Unlock()
	if r.Pub == nil {
		panic("value is nil")
	}
	pubASN1, err := x509.MarshalPKIXPublicKey(r.Pub)
	gserror.CheckError_panic(err)

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// EncryptWithPublicKey encrypts data with public key
func (r *RSA) EncryptWithPublicKeyAnyLen(msg []byte) []byte {
	if r.Pub == nil {
		panic("value is nil")
	}
	if r.Pub.Size() != g_rsaKeyLen {
		panic("r.Pri.Size() != g_rsaKeyLen")
	}

	var datals []byte

	count := len(msg) / g_oneLen
	for i := 1; i <= count; i++ {
		datals = append(datals, EncryptWithPublicKey(msg[g_oneLen*(i-1):g_oneLen*i], r.Pub)...)
	}
	if len(msg)%g_oneLen != 0 {
		datals = append(datals, EncryptWithPublicKey(msg[g_oneLen*count:], r.Pub)...)
	}

	return datals
}

// DecryptWithPrivateKey decrypts data with private key
func (r *RSA) DecryptWithPrivateKeyAnyLen(ciphertext []byte) []byte {

	if r.Pri == nil {
		panic("value is nil")
	}
	if r.Pri.Size() != g_rsaKeyLen {
		panic("r.Pri.Size() != g_rsaKeyLen")
	}

	count := len(ciphertext) / g_rsaKeyLen
	data := make([]byte, 0, g_oneLen*count)

	for i := 1; i <= count; i++ {
		data = append(data, DecryptWithPrivateKey(ciphertext[g_rsaKeyLen*(i-1):g_rsaKeyLen*i], r.Pri)...)
	}

	return data
}

func (r *RSA) PrivateKeyToBase64() []byte {
	if r.Pri == nil {
		panic("value is nil")
	}
	src := PrivateKeyToBytes(r.Pri)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func (r *RSA) PublicKeyToBase64() []byte {
	if r.Pub == nil {
		panic("value is nil")
	}
	src := PublicKeyToBytes(r.Pub)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

// Deep copy, new obj.
func (r *RSA) NewRSA() *RSA {
	r.lock_pri.Lock()
	defer r.lock_pri.Unlock()
	r.lock_pub.Lock()
	defer r.lock_pub.Unlock()
	if r.Pri != nil {
		return NewRSAObjFromBytes(PrivateKeyToBytes(r.Pri))
	}
	if r.Pub != nil {
		return NewRSAObjPub(PublicKeyFromBytes(PublicKeyToBytes(r.Pub)))
	}
	panic("pri and pub is nil")
}

// Deep copy, new obj.
func (r *RSA) NewRSAPub() *RSA {
	//r.lock_pri.Lock()
	//defer r.lock_pri.Unlock()
	r.lock_pub.Lock()
	defer r.lock_pub.Unlock()
	if r.Pub != nil {
		return NewRSAObjPub(PublicKeyFromBytes(PublicKeyToBytes(r.Pub)))
	}
	panic("pub is nil")
}

//func (r *RSA) GetPrivateKey() *rsa.PrivateKey {}
