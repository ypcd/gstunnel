package main

import (
	"sync"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

func Test_client_NetPipe_st(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_client_NetPipe_st] start.")
	inTest_client_NetPipe(t, false)
	logger_test.Println("[Test_client_NetPipe_st] end.")
}
func Test_client_NetPipe_mt(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_client_NetPipe_mt] start.")
	inTest_client_NetPipe(t, true)
	logger_test.Println("[Test_client_NetPipe_mt] end.")
}

func Test_client_NetPipe_m(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
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
	logger_test.Println("[Test_client_NetPipe_m] end.")
}

func Test_client_NetPipe_loop(t *testing.T) {
	logger_test.Println("[Test_client_NetPipe_loop] start.")
	for i := 0; i < 6; i++ {
		logger_test.Printf("loop count: %d", i)
		inTest_client_NetPipe(t, false)
		forceGC()
	}
	logger_test.Println("[Test_client_NetPipe_loop] end.")
}
