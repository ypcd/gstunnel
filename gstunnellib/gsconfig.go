package gstunnellib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	. "github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

type GsConfig struct {
	Listen  string
	Servers []string
	Key     string

	Tmr_display_time   int
	Tmr_changekey_time int
	NetworkTimeout     int

	Debug    bool
	Mt_model bool
	WebUI    bool
}

func (gs *GsConfig) GetServer_rand() string {
	return gs.Servers[GetRDCInt_max(int64(len(gs.Servers)))]
}
func (gs *GsConfig) GetServers() []string {
	return gs.Servers
}

func CreateGsconfig(confg string) *GsConfig {
	f, err := os.Open(confg)
	CheckError_exit(err)
	defer f.Close()

	buf, err := io.ReadAll(f)
	CheckError_exit(err)

	//fmt.Println(string(buf))
	var gsconfig GsConfig

	gsconfig.Tmr_display_time = 6
	gsconfig.Tmr_changekey_time = 60
	gsconfig.NetworkTimeout = 60
	gsconfig.Debug = false
	gsconfig.Mt_model = true
	gsconfig.WebUI = false

	err = json.Unmarshal(buf, &gsconfig)
	CheckError(err)

	if gsconfig.Servers == nil || gsconfig.Key == "" || gsconfig.Listen == "" {
		g_logger.Fatalln("Gstunnel config is error.")
	}

	key1, err := base64.StdEncoding.DecodeString(gsconfig.Key)
	CheckError_exit(err)
	if len(key1) != gsbase.G_AesKeyLen {
		CheckError_exit(fmt.Errorf("error: the key is not %d bytes", gsbase.G_AesKeyLen))
	}
	gsconfig.Key = string(key1)

	return &gsconfig
}
