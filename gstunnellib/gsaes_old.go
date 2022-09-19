package gstunnellib

import (
	"crypto/aes"
	"crypto/cipher"
)

type gsaes_old struct {
	cpr        cipher.Block
	cfbE, cfbD cipher.Stream
}

func createAes_old(key []byte) gsaes_old {
	var v1 gsaes_old

	v1.cpr, _ = aes.NewCipher(key)

	v1.cfbE = cipher.NewCFBEncrypter(v1.cpr, commonIV)
	v1.cfbD = cipher.NewCFBDecrypter(v1.cpr, commonIV)
	return v1
}

func (a *gsaes_old) encrypter(data []byte) []byte {
	dst := make([]byte, len(data))
	a.cfbE.XORKeyStream(dst, data)
	return dst
}

func (a *gsaes_old) decrypter(data []byte) []byte {
	dst := make([]byte, len(data))
	a.cfbD.XORKeyStream(dst, data)
	return dst
}
