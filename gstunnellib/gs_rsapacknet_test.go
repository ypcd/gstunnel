package gstunnellib

import (
	"bytes"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func getRSAPack(v IGSRSAPackNet) *gsRSAPackImp {
	r, ok := v.(*gsRSAPackNetImp)
	if !ok {
		panic("error")
	}
	rp, ok := r.apack.(*gsRSAPackImp)
	if !ok {
		panic("error")
	}
	return rp
}

func Test_rsapacknet_ClientPublicKeyPack(t *testing.T) {
	key1 := string(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	spri, spub := gsrsa.GenerateKeyPair(4096)
	//cpri, cpub := gsrsa.GenerateKeyPair(4096)

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	serverPack := NewGSRSAPackNetImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	clientPack := NewGSRSAPackNetImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	pubkpack := clientPack.ClientPublicKeyPack()
	serverPack.GetDecryDataFormBytes(pubkpack)

	rpServer := getRSAPack(serverPack)
	rpClient := getRSAPack(clientPack)

	if !bytes.Equal(rpServer.clientkey.PublicKeyToBytes(), rpClient.clientkey.PublicKeyToBytes()) {
		panic("error")
	}
}

func Test_rsapacknet_pack_unpack_gstserver(t *testing.T) {

	fbuf := gsrand.GetRDBytes(50000)

	spri, _ := gsrsa.GenerateKeyPair(4096)

	key1 := string(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	a1 := NewGSRSAPackNetImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	tmp := a1.Packing(fbuf)

	outbuf, _ := a1.GetDecryDataFormBytes(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

}

/*
func Test_rsapacknet_pack_unpack_gstclient(t *testing.T) {

	fbuf := gsrand.GetRDBytes(50000)

	_, spub := gsrsa.GenerateKeyPair(4096)

	a1 := NewGSRSAPackNetImp(string(gsrand.GetRDBytes(gsbase.G_AesKeyLen)), spub)
	tmp := a1.Packing(fbuf)

	outbuf, _ := a1.GetDecryDataFormBytes(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

}
*/
func Test_rsapacknet_pack_unpack_gstserver_and_gstclient(t *testing.T) {

	fbuf := gsrand.GetRDBytes(50000)

	spri, spub := gsrsa.GenerateKeyPair(4096)

	key1 := string(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewGSRSAPackNetImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewGSRSAPackNetImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	tmp := a1.Packing(fbuf)

	outbuf, _ := a2.GetDecryDataFormBytes(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

	fbuf = gsrand.GetRDBytes(50000)
	tmp = a2.Packing(fbuf)

	outbuf, _ = a1.GetDecryDataFormBytes(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

}

func Test_rsapacknet_IsTheVersionConsistent(t *testing.T) {
	key1 := string(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	spri, spub := gsrsa.GenerateKeyPair(4096)

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewGSRSAPackNetImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewGSRSAPackNetImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	tmp := a1.PackVersion()
	_, _ = a2.GetDecryDataFormBytes(tmp)

	tmp = a2.PackVersion()
	_, _ = a1.GetDecryDataFormBytes(tmp)

}

func Test_rsapacknet_ChangeCryKeyFromGSTServer(t *testing.T) {
	keydata := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	key1 := string(keydata)
	rawdata := gsrand.GetRDBytes(1024)

	spri, spub := gsrsa.GenerateKeyPair(4096)

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewGSRSAPackNetImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewGSRSAPackNetImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	rpServer := getRSAPack(a1)
	rpClient := getRSAPack(a2)

	if !bytes.Equal(rpServer.a.encrypter(rawdata), rpClient.a.encrypter(rawdata)) {
		panic("error.")
	}

	pubkpack := a2.ClientPublicKeyPack()
	a1.GetDecryDataFormBytes(pubkpack)

	for i := 0; i < 50; i++ {
		tmp, _ := a1.ChangeCryKeyFromGSTServer()
		_, _ = a2.GetDecryDataFormBytes(tmp)

		rawdata = gsrand.GetRDBytes(1024)

		rpServer := getRSAPack(a1)
		rpClient := getRSAPack(a2)

		if !bytes.Equal(rpServer.a.encrypter(rawdata), rpClient.a.encrypter(rawdata)) {
			panic("error.")
		}
	}
}

func Test_rsapacknet_ChangeCryKeyFromGSTClient(t *testing.T) {
	keydata := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	key1 := string(keydata)
	rawdata := gsrand.GetRDBytes(1024)

	spri, spub := gsrsa.GenerateKeyPair(4096)

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewGSRSAPackNetImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewGSRSAPackNetImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	rpServer := getRSAPack(a1)
	rpClient := getRSAPack(a2)

	if !bytes.Equal(rpServer.a.encrypter(rawdata), rpClient.a.encrypter(rawdata)) {
		panic("error.")
	}

	pubkpack := a2.ClientPublicKeyPack()
	a1.GetDecryDataFormBytes(pubkpack)

	for i := 0; i < 50; i++ {

		tmp, _ := a2.ChangeCryKeyFromGSTClient()
		_, _ = a1.GetDecryDataFormBytes(tmp)

		rawdata = gsrand.GetRDBytes(1024)

		rpServer := getRSAPack(a1)
		rpClient := getRSAPack(a2)

		if !bytes.Equal(rpServer.a.encrypter(rawdata), rpClient.a.encrypter(rawdata)) {
			panic("error.")
		}

		tmp, _ = a1.ChangeCryKeyFromGSTServer()
		_, _ = a2.GetDecryDataFormBytes(tmp)

		rawdata = gsrand.GetRDBytes(1024)

		rpServer = getRSAPack(a1)
		rpClient = getRSAPack(a2)

		if !bytes.Equal(rpServer.a.encrypter(rawdata), rpClient.a.encrypter(rawdata)) {
			panic("error.")
		}
	}
}
