package gstunnellib

import (
	"encoding/json"
	. "gstunnellib/gsrand"
	"io"
	"os"
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
	Listen             string
	Servers            []string
	Key                string
	Debug              bool
	Tmr_display_time   int
	Tmr_changekey_time int
	Mt_model           bool
}

func (gs *GsConfig) GetServer_rand() string {
	return gs.Servers[GetRDCInt_max(int64(len(gs.Servers)))]
}
func (gs *GsConfig) GetServers() []string {
	return gs.Servers
}

func CreateGsconfig(confn string) *GsConfig {
	f, err := os.Open(confn)
	checkError(err)
	defer f.Close()

	buf, err := io.ReadAll(f)
	checkError(err)

	//fmt.Println(string(buf))
	var gsconfig GsConfig

	gsconfig.Debug = false
	gsconfig.Tmr_display_time = 5
	gsconfig.Tmr_changekey_time = 60
	gsconfig.Mt_model = true

	err = json.Unmarshal(buf, &gsconfig)
	checkError(err)
	/*
		if gsconfig.Tmr_display_time == 0 {
			gsconfig.Tmr_display_time = 5
		}
		if gsconfig.Tmr_changekey_time == 0 {
			gsconfig.Tmr_changekey_time = 60
		}
	*/
	if gsconfig.Servers == nil {
		logger.Fatalln("gsconfig.Servers==nil")
	}
	return &gsconfig
}
