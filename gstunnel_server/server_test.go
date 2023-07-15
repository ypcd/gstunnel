package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

func init() {
	init_server_test()
	logger_test = g_Logger
	g_networkTimeout = time.Second * 10
}

func Test_server_NetPipe_st(t *testing.T) {
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Print("\n\n")
	logger_test.Println("[Test_server_NetPipe_st] start.")
	inTest_server_NetPipe(t, false)
	logger_test.Print("[Test_server_NetPipe_st] end.\n\n")
}
func Test_server_NetPipe_mt(t *testing.T) {
	count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_mt] start.")
	inTest_server_NetPipe(t, true)
	logger_test.Print("[Test_server_NetPipe_mt] end.\n\n")
}

func Test_server_NetPipe_m(t *testing.T) {
	count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe] start.")
	//inTest_server_NetPipe_m(t)
	gwg := sync.WaitGroup{}

	inTest_server_NetPipe_go_init()
	for i := 0; i < 2; i++ {
		gwg.Add(1)
		go inTest_server_NetPipe_go(t, &gwg)
	}
	gwg.Wait()
	logger_test.Print("[Test_server_NetPipe] end.\n\n")
}

func Test_server_NetPipe_loop(t *testing.T) {
	logger_test.Println("[Test_server_NetPipe_loop] start.")
	for i := 0; i < 6; i++ {
		logger_test.Printf("loop count: %d", i)
		inTest_server_NetPipe(t, false)
		forceGC()
	}
	logger_test.Print("[Test_server_NetPipe_loop] end.\n\n")
}

func Test_server_NetPipe_errorData(t *testing.T) {
	count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[notest_Test_server_NetPipe_errorData] start.")
	inTest_server_NetPipe_errorData(t, false)
	logger_test.Println("[notest_Test_server_NetPipe_errorData] end.")
}

func Test_server_NetPipe_errorKey(t *testing.T) {
	count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_errorKey] start.")
	inTest_server_NetPipe_errorKey(t, false)
	logger_test.Print("[Test_server_NetPipe_errorKey] end.\n\n")
}

func Test_server_timeout(t *testing.T) {
	count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_timeout] start.")
	inTest_server_timeout(t, false)
	logger_test.Print("[Test_server_timeout] end.\n\n")
}

func Test_server_timeout2(t *testing.T) {
	count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_timeout] start.")
	inTest_server_timeout(t, true)
	logger_test.Print("[Test_server_timeout] end.\n\n")
}

/*
func Test_init1(t *testing.T) {
	init_server_run()
	init_server_run()
}

func Test_init2(t *testing.T) {
	init_server_test()
	init_server_test()
}

func Test_init3(t *testing.T) {
	init_server_test()
	init_server_run()
}
*/

func Test_server_NetPipe_nGiB(t *testing.T) {
	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}
	if count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because count_inTest_server_NetPipe > 0.")
		return
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	ss := gstestpipe.NewServiceServerNone()
	gst := gstestpipe.NewGstPiPoDefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	MiBNumber := 2 * 1024
	testCacheSize := MiBNumber * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(ss.GetClientConn(), gst.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := ss.GetServerConn()

	SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	//rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("Test_server_NetPipe_nGiB data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, g_net_read_size)
		rbufsz := 0

		for {
			server.SetReadDeadline(time.Now().Add(testReadTimeOut))
			re, err := server.Read(buf)
			//t.Logf("server read len: %d", re)
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				gstunnellib.CheckError_test_noExit(err, t)
				return
			} else {
				gstunnellib.CheckError_test(err, t)
			}

			//rbuf = append(rbuf, buf[:re]...)
			rbufsz += re

			if rbufsz == len(SendData) {
				return
			}
		}

	}(server)

	_, err := io.Copy(server, bytes.NewBuffer(SendData))
	checkError(err)
	forceGC()
	//time.Sleep(time.Second * 6)
	wg.Wait()
	server.Close()
	forceGC()
	//time.Sleep(time.Second * 60)

	//	if !bytes.Equal(SendData, rbuf) {
	//		t.Fatal("Error: SendData != rbuf.")
	//	}
	logger_test.Println("[Test_server_NetPipe_nGiB] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/s]:", float64(MiBNumber)/t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()
}

func in_Test_server_NetPipe_nGiB_wg(t *testing.T, inwg *sync.WaitGroup) {
	defer inwg.Done()
	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}
	if count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because count_inTest_server_NetPipe > 0.")
		return
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	//rawsever
	ss := gstestpipe.NewServiceServerNone()
	gst := gstestpipe.NewGstPiPoDefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	MiBNumber := 2 * 1024
	testCacheSize := MiBNumber * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(ss.GetClientConn(), gst.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	//rawsever-server
	server := ss.GetServerConn()

	SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	//rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("Test_server_NetPipe_nGiB data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, g_net_read_size)
		rbufsz := 0

		for {
			server.SetReadDeadline(time.Now().Add(testReadTimeOut))
			re, err := server.Read(buf)
			//t.Logf("server read len: %d", re)
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				gstunnellib.CheckError_test_noExit(err, t)
				return
			} else {
				gstunnellib.CheckError_test(err, t)
			}

			//rbuf = append(rbuf, buf[:re]...)
			rbufsz += re

			if rbufsz == len(SendData) {
				return
			}
		}

	}(server)

	_, err := io.Copy(server, bytes.NewBuffer(SendData))
	checkError(err)
	forceGC()
	//time.Sleep(time.Second * 6)
	wg.Wait()
	server.Close()
	forceGC()
	//time.Sleep(time.Second * 60)

	//	if !bytes.Equal(SendData, rbuf) {
	//		t.Fatal("Error: SendData != rbuf.")
	//	}
	logger_test.Println("[Test_server_NetPipe_nGiB] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/s]:", float64(MiBNumber)/t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()
}

func Test_server_NetPipe_nGiB_mt(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go in_Test_server_NetPipe_nGiB_wg(t, &wg)
	}
	wg.Wait()
}
