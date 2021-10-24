package gstunnellib

import (
	"testing"
)

func Test_GsConfig(t *testing.T) {
	gs := CreateGsconfig("config.client.json")
	t.Log(gs)
}

func Test_GsConfig_getserver(t *testing.T) {
	gs := CreateGsconfig("config3.client.json")
	t.Log(gs)
	re := gs.GetServers()
	t.Log(re)
	for i := 0; i < 10; i++ {
		t.Log(gs.GetServer_rand())
	}
}
