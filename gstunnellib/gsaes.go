package gstunnellib

import (
	"crypto/aes"
	"crypto/cipher"
)

type gsaes struct {
	cpr cipher.Block
	gcm cipher.AEAD
}

func createAes(key []byte) gsaes {
	var v1 gsaes
	var err error

	v1.cpr, err = aes.NewCipher(key)
	checkError_panic(err)

	v1.gcm, err = cipher.NewGCMWithNonceSize(v1.cpr, len(commonIV))
	checkError_panic(err)

	//v2 := v1.gcm.NonceSize()
	//_ = v2
	return v1
}

func (a *gsaes) encrypter(data []byte) []byte {
	return a.gcm.Seal(nil, commonIV, data, nil)
}

func (a *gsaes) decrypter(data []byte) []byte {
	decry, err := a.gcm.Open(nil, commonIV, data, nil)
	checkError_panic(err)
	return decry
}
