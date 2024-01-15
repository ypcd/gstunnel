package gstunnellib

import (
	"bytes"
	"crypto/rsa"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func Test_GsConfig(t *testing.T) {
	gs := CreateGsconfig("config.test.json")
	t.Log(gs)
}

func Test_GsConfig_getservers(t *testing.T) {
	gs := CreateGsconfig("config3.test.json")
	t.Log(gs)
	re := gs.GetServers()
	t.Log(re)
	for i := 0; i < 10; i++ {
		t.Log(gs.GetServer_rand())
	}
}

func in_gsconfig_Rsa_RDdata_test(pri *rsa.PrivateKey, pub *rsa.PublicKey) {
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

func Test_gsconfig_rsa_serverPrivate_public(t *testing.T) {
	conf1 := CreateGsconfig("config-rsa.test.json")
	pri := conf1.GetRSAServerPrivate()
	pub := conf1.GetRSAServerPublic()
	in_gsconfig_Rsa_RDdata_test(pri, pub)
}

func Test_gsconfig_rsa_RSAserver(t *testing.T) {
	conf1 := CreateGsconfig("config-rsa.test.json")
	r1 := conf1.GetRSAServer()
	in_gsconfig_Rsa_RDdata_test(r1.Pri, r1.Pub)
}
