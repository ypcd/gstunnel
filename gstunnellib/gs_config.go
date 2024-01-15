package gstunnellib

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type GsConfig struct {
	Listen  string
	Servers []string
	Key     string

	RSAServerPrivate string
	RSAServerPublic  string

	Tmr_display_time   int
	Tmr_changekey_time int
	NetworkTimeout     int

	Debug    bool
	Mt_model bool
	WebUI    bool
}

func (c *GsConfig) GetServer_rand() string {
	return c.Servers[gsrand.GetRDInt_max(int64(len(c.Servers)))]
}
func (c *GsConfig) GetServers() []string {
	return c.Servers
}

func (c *GsConfig) GetRSAServerPrivate() *rsa.PrivateKey {
	return gsrsa.PrivateKeyFromBase64([]byte(c.RSAServerPrivate))
}

func (c *GsConfig) GetRSAServerPublic() *rsa.PublicKey {
	return gsrsa.PublicKeyFromBase64([]byte(c.RSAServerPublic))
}

// Deep copy, new obj.
func (c *GsConfig) GetRSAServer() *gsrsa.RSA {
	if c.RSAServerPrivate != "" {
		return gsrsa.NewRSAObjFromBase64([]byte(c.RSAServerPrivate))
	} else if c.RSAServerPublic != "" {
		return gsrsa.NewRSAObjFromPubKeyBase64([]byte(c.RSAServerPublic))
	}
	panic("error")
}

func CreateGsconfig(confg string) *GsConfig {
	f, err := os.Open(confg)
	checkError_exit(err)
	defer f.Close()

	buf, err := io.ReadAll(f)
	checkError_exit(err)

	//fmt.Println(string(buf))
	var gsconfig GsConfig

	gsconfig.Tmr_display_time = 6
	gsconfig.Tmr_changekey_time = 60
	gsconfig.NetworkTimeout = 60
	gsconfig.Debug = false
	gsconfig.Mt_model = true
	gsconfig.WebUI = false

	err = json.Unmarshal(buf, &gsconfig)
	checkError(err)

	if gsconfig.Servers == nil || gsconfig.Key == "" || gsconfig.Listen == "" {
		G_logger.Fatalln("Gstunnel config is error.")
	}
	if gsconfig.RSAServerPrivate == "" && gsconfig.RSAServerPublic == "" {
		G_logger.Fatalln("Gstunnel config is error. RSAServerPrivate and RSAServerPublic is not exists.")
	}

	//server rsa key

	key1, err := base64.StdEncoding.DecodeString(gsconfig.Key)
	checkError_exit(err)
	if len(key1) != gsbase.G_AesKeyLen {
		checkError_exit(fmt.Errorf("error: the key is not %d bytes", gsbase.G_AesKeyLen))
	}
	gsconfig.Key = string(key1)

	return &gsconfig
}
