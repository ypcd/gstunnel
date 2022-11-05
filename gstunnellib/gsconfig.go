package gstunnellib

import (
	"encoding/json"
	"io"
	"os"

	. "github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

type gsConfig_1 struct {
	Listen             string
	Server             string
	Key                string
	Debug              bool
	Tmr_display_time   int
	Tmr_changekey_time int
	Mt_model           bool
}

type GsConfig struct {
	Listen  string
	Servers []string
	Key     string

	Tmr_display_time   int
	Tmr_changekey_time int
	NetworkTimeout     int

	Debug    bool
	Mt_model bool
}

func (gs *GsConfig) GetServer_rand() string {
	return gs.Servers[GetRDCInt_max(int64(len(gs.Servers)))]
}
func (gs *GsConfig) GetServers() []string {
	return gs.Servers
}

func CreateGsconfig(confg string) *GsConfig {
	f, err := os.Open(confg)
	checkError_exit(err)
	defer f.Close()

	buf, err := io.ReadAll(f)
	checkError(err)

	//fmt.Println(string(buf))
	var gsconfig GsConfig

	gsconfig.Tmr_display_time = 6
	gsconfig.Tmr_changekey_time = 60
	gsconfig.NetworkTimeout = 60
	gsconfig.Debug = false
	gsconfig.Mt_model = true

	err = json.Unmarshal(buf, &gsconfig)
	checkError(err)

	if gsconfig.Servers == nil || gsconfig.Key == "" || gsconfig.Listen == "" {
		logger.Fatalln("Gstunnel config is error.")
	}
	return &gsconfig
}
