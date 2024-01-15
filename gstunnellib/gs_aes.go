package gstunnellib

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
)

type aesItem struct {
	cpr cipher.Block
	gcm cipher.AEAD
}

func createAesItem(key []byte) *aesItem {
	var v1 aesItem
	var err error

	v1.cpr, err = aes.NewCipher(key)
	checkError_panic(err)

	v1.gcm, err = cipher.NewGCMWithNonceSize(v1.cpr, len(g_commonIV))
	checkError_panic(err)

	return &v1
}

func (a *aesItem) encrypter(data []byte) []byte {
	return a.gcm.Seal(nil, g_commonIV, data, nil)
}

func (a *aesItem) decrypter(data []byte) []byte {
	decry, err := a.gcm.Open(nil, g_commonIV, data, nil)
	checkError_panic(err)
	return decry
}

type gsaes struct {
	aes *aesItem
}

func createAes(key []byte) gsaes {
	if len(key) != gsbase.G_AesKeyLen {
		checkError_exit(
			fmt.Errorf("error: the key is not %d bytes", gsbase.G_AesKeyLen))
	}

	var v1 gsaes

	v1.aes = createAesItem(key)

	return v1
}

func (a *gsaes) encrypter(data []byte) []byte {

	data = a.aes.encrypter(data)

	return data
}

func (a *gsaes) decrypter(data []byte) []byte {
	data = a.aes.decrypter(data)

	return data
}
