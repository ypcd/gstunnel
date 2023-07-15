package main

import (
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

func init() {
	init_client_test()
	logger_test = g_Logger
	g_networkTimeout = time.Second * 10
}

func Test_client_NetPipe_st(t *testing.T) {
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Print("\n\n")
	logger_test.Println("[Test_client_NetPipe_st] start.")
	inTest_client_NetPipe(t, false)
	logger_test.Print("[Test_client_NetPipe_st] end.\n\n")
}
func Test_client_NetPipe_mt(t *testing.T) {
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println()
	logger_test.Println("[Test_client_NetPipe_mt] start.")
	inTest_client_NetPipe(t, true)
	logger_test.Print("[Test_client_NetPipe_mt] end.\n\n")
}

func Test_client_NetPipe_m(t *testing.T) {
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_client_NetPipe_m] start.")
	//inTest_client_NetPipe_m(t)
	gwg := sync.WaitGroup{}

	inTest_client_NetPipe_go_init()
	for i := 0; i < 2; i++ {
		gwg.Add(1)
		go inTest_client_NetPipe_go(t, &gwg)
	}
	gwg.Wait()
	logger_test.Print("[Test_client_NetPipe_m] end.\n\n")
}

func Test_client_NetPipe_loop(t *testing.T) {
	logger_test.Println("[Test_client_NetPipe_loop] start.")
	for i := 0; i < 6; i++ {
		logger_test.Printf("loop count: %d", i)
		inTest_client_NetPipe(t, false)
		forceGC()
	}
	logger_test.Print("[Test_client_NetPipe_loop] end.\n\n")
}

func Test_client_timeout(t *testing.T) {
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_client_timeout] start.")
	inTest_client_timeout(t, false)
	logger_test.Print("[Test_client_timeout] end.\n\n")
}

func Test_client_timeout2(t *testing.T) {
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_client_timeout] start.")
	inTest_client_timeout(t, true)
	logger_test.Print("[Test_client_timeout] end.\n\n")
}
