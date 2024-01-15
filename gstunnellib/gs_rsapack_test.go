package gstunnellib

import (
	"bytes"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func Test_rsapack_ClientPublicKeyPack(t *testing.T) {
	key1 := string(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	spri, spub := gsrsa.GenerateKeyPair(4096)
	//cpri, cpub := gsrsa.GenerateKeyPair(4096)
	serverPri := gsrsa.NewRSAObj(spri)
	serverPub := gsrsa.NewRSAObjPub(spub)
	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	serverPack := NewRSAPackImpWithGSTServer(key1, serverPri)
	clientPack := NewRSAPackImpWithGSTClient(key1, serverPub, clientkey)

	pubkpack := clientPack.ClientPublicKeyPack()
	serverPack.Unpack(pubkpack)

	if !bytes.Equal(serverPack.clientkey.PublicKeyToBytes(), clientPack.clientkey.PublicKeyToBytes()) {
		panic("error")
	}
}

func Test_rsapack_pack_unpack_gstserver(t *testing.T) {

	fbuf := gsrand.GetRDBytes(50000)

	//	spri, _ := gsrsa.GenerateKeyPair(4096)
	rkey1 := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewRSAPackImpWithGSTServer(string(gsrand.GetRDBytes(gsbase.G_AesKeyLen)), rkey1)
	tmp := a1.Packing(fbuf)

	outbuf, _ := a1.Unpack(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

}

func inTest_rsapack_unpackex(rp *gsRSAPackImp) {

	rddata := gsrand.GetRDBytes(50000)

	//	spri, _ := gsrsa.GenerateKeyPair(4096)
	a1 := rp
	packdata := a1.Packing(rddata)
	a1.PackVersion()
	outbuf, err := a1.UnpackEx(packdata)
	checkError_panic(err)
	if !outbuf.IsPOGen() {
		panic("error")
	}
	if !bytes.Equal(rddata, outbuf.GenData) {
		panic("error")
	}

	//pack := gspackoper.NewPackPOVersion()
	packdata = a1.PackVersion()

	outbuf, err = a1.UnpackEx(packdata)
	checkErrorEx_panic(err)

	if !outbuf.IsPOVersion() {
		panic("error")
	}
	if outbuf.Version != gsbase.G_Version {
		panic("")
	}

	//clientpub := gsrsa.NewRSAObjPub(gsrsa.NewGenRSAObj(512).Pub)
	//pack := a1.ClientPublicKeyPack()
	packdata = a1.ClientPublicKeyPack()

	outbuf, err = a1.UnpackEx(packdata)
	checkErrorEx_panic(err)

	if !outbuf.IsClientPubKey() {
		panic("error")
	}
	if !outbuf.ClientKey.Pub.Equal(a1.clientkey.Pub) {
		panic("")
	}

	packdata = a1.PackPOHello()

	outbuf, err = a1.UnpackEx(packdata)
	checkErrorEx_panic(err)

	if !outbuf.IsPOHello() {
		panic("error")
	}

}

func Test_rsapack_unpackex(t *testing.T) {
	//	spri, _ := gsrsa.GenerateKeyPair(4096)
	rkey1 := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)
	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := newRSAPackImp(string(gsrand.GetRDBytes(gsbase.G_AesKeyLen)), rkey1, clientkey)

	inTest_rsapack_unpackex(a1)
}

/*
func Test_rsapack_pack_unpack_gstclient(t *testing.T) {

	fbuf := gsrand.GetRDBytes(50000)

	_, spub := gsrsa.GenerateKeyPair(4096)

	a1 := NewGSRSAPackImp(string(gsrand.GetRDBytes(gsbase.G_AesKeyLen)), spub)
	tmp := a1.Packing(fbuf)

	outbuf, _ := a1.Unpack(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

}*/

func Test_rsapack_pack_unpack_gstserver_and_gstclient(t *testing.T) {

	fbuf := gsrand.GetRDBytes(50000)

	spri, spub := gsrsa.GenerateKeyPair(4096)

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	key1 := string(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	a1 := NewRSAPackImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewRSAPackImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	tmp := a1.Packing(fbuf)

	outbuf, _ := a2.Unpack(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

	fbuf = gsrand.GetRDBytes(50000)

	tmp = a2.Packing(fbuf)

	outbuf, _ = a1.Unpack(tmp)

	if !bytes.Equal(fbuf, outbuf) {
		panic("error")
	}

}

func Test_rsapack_IsTheVersionConsistent(t *testing.T) {
	key1 := string(gsrand.GetRDBytes(gsbase.G_AesKeyLen))

	spri, spub := gsrsa.GenerateKeyPair(4096)

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewRSAPackImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewRSAPackImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	tmp := a1.PackVersion()
	_, _ = a2.Unpack(tmp)

	tmp = a2.PackVersion()
	_, _ = a1.Unpack(tmp)

}

func Test_rsapack_ChangeCryKeyFromGSTServer(t *testing.T) {
	keydata := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	key1 := string(keydata)
	rawdata := gsrand.GetRDBytes(1024)

	spri, spub := gsrsa.GenerateKeyPair(4096)

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewRSAPackImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewRSAPackImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	if !bytes.Equal(a1.a.encrypter(rawdata), a2.a.encrypter(rawdata)) {
		panic("error.")
	}

	pubkpack := a2.ClientPublicKeyPack()
	a1.Unpack(pubkpack)

	for i := 0; i < 50; i++ {
		tmp, _ := a1.ChangeCryKeyFromGSTServer()
		_, _ = a2.Unpack(tmp)

		rawdata = gsrand.GetRDBytes(1024)
		if !bytes.Equal(a1.a.encrypter(rawdata), a2.a.encrypter(rawdata)) {
			panic("error.")
		}
	}
}

func Test_rsapack_ChangeCryKeyFromGSTClient(t *testing.T) {
	keydata := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	key1 := string(keydata)
	rawdata := gsrand.GetRDBytes(1024)

	spri, spub := gsrsa.GenerateKeyPair(4096)
	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	a1 := NewRSAPackImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	a2 := NewRSAPackImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)

	if !bytes.Equal(a1.a.encrypter(rawdata), a2.a.encrypter(rawdata)) {
		panic("error.")
	}

	pubkpack := a2.ClientPublicKeyPack()
	a1.Unpack(pubkpack)

	for i := 0; i < 50; i++ {

		tmp, _ := a2.ChangeCryKeyFromGSTClient()
		_, _ = a1.Unpack(tmp)

		rawdata = gsrand.GetRDBytes(1024)
		if !bytes.Equal(a1.a.encrypter(rawdata), a2.a.encrypter(rawdata)) {
			panic("error.")
		}

		tmp, _ = a1.ChangeCryKeyFromGSTServer()
		_, _ = a2.Unpack(tmp)

		rawdata = gsrand.GetRDBytes(1024)
		if !bytes.Equal(a1.a.encrypter(rawdata), a2.a.encrypter(rawdata)) {
			panic("error.")
		}
	}
}

func Test_rsapack_IsExistsClientKey(t *testing.T) {
	spri, spub := gsrsa.GenerateKeyPair(4096)
	keydata := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	key1 := string(keydata)

	rp := NewRSAPackImpWithGSTServer(key1, gsrsa.NewRSAObj(spri))
	re := rp.IsExistsClientKey()
	if re {
		panic("error")
	}

	clientkey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	rp2 := NewRSAPackImpWithGSTClient(key1, gsrsa.NewRSAObjPub(spub), clientkey)
	re = rp2.IsExistsClientKey()
	if !re {
		panic("error")
	}
}
