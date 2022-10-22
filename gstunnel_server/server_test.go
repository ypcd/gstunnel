package main

import (
	"sync"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

func Test_server_NetPipe_st(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_st] start.")
	inTest_server_NetPipe(t, false)
	logger_test.Println("[Test_server_NetPipe_st] end.")
}
func Test_server_NetPipe_mt(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_mt] start.")
	inTest_server_NetPipe(t, true)
	logger_test.Println("[Test_server_NetPipe_mt] end.")
}

func Test_server_NetPipe_m(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe] start.")
	//inTest_server_NetPipe_m(t)
	gwg := sync.WaitGroup{}

	inTest_server_NetPipe_go_init()
	for i := 0; i < 2; i++ {
		go inTest_server_NetPipe_go(t, &gwg)
		gwg.Add(1)
	}
	gwg.Wait()
	logger_test.Println("[Test_server_NetPipe] end.")
}

func Test_server_NetPipe_loop(t *testing.T) {
	logger_test.Println("[Test_server_NetPipe_loop] start.")
	for i := 0; i < 6; i++ {
		logger_test.Printf("loop count: %d", i)
		inTest_server_NetPipe(t, false)
		forceGC()
	}
	logger_test.Println("[Test_server_NetPipe_loop] end.")
}

func Test_server_NetPipe_errorData(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_st] start.")
	inTest_server_NetPipe_errorData(t, false)
	logger_test.Println("[Test_server_NetPipe_st] end.")
}

func Test_server_NetPipe_errorKey(t *testing.T) {
	if gstunnellib.RunTime_Debug {
		defer gstunnellib.RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_st] start.")
	inTest_server_NetPipe_errorKey(t, false)
	logger_test.Println("[Test_server_NetPipe_st] end.")
}
