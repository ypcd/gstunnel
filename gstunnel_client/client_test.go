package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

var g_count_inTest_server_NetPipe int = 0

func init() {
	init_client_test()
	logger_test = g_Logger
	g_networkTimeout = time.Second * 10
}

func Test_client_NetPipe_st(t *testing.T) {
	g_count_inTest_server_NetPipe++
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Print("\n\n")
	logger_test.Println("[Test_client_NetPipe_st] start.")
	inTest_client_NetPipe(t, false)
	logger_test.Print("[Test_client_NetPipe_st] end.\n\n")
}
func Test_client_NetPipe_mt(t *testing.T) {
	g_count_inTest_server_NetPipe++
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

// 适合测试gstunnel是否正确传输数据。
// 但是因为计算数据哈希值，增加了cpu占用，导致吞吐量降低，所有不适合测试最大吞吐量。
func Test_server_NetPipe_nGiB_hash(t *testing.T) {

	//samplesNum := 100
	//f, err := os.Create(fmt.Sprintf("cpuprof-%d.prof", samplesNum))
	//checkError_panic(err)
	//defer f.Close()
	//runtime.SetCPUProfileRate(samplesNum)
	//pprof.StartCPUProfile(f)

	//debug.SetGCPercent(-1)
	//debug.SetMemoryLimit(1024 * 1024 * 1024 * 10)

	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}

	/*
		if g_count_inTest_server_NetPipe > 0 {
			logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because g_count_inTest_server_NetPipe > 0.")
			return
		}
	*/
	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	rawClient := gstestpipe.NewServiceServerNone()
	gstServer := gstestpipe.NewGstPiPoDefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	MiBNumber := 100 * 100 / 2
	testCacheSize := MiBNumber * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawClient.GetClientConn(), gstServer.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := rawClient.GetServerConn()

	var SendData []byte
	//rbuf := make([]byte, 0, len(SendData))

	hashw := sha256.New()
	hashr := sha256.New()

	logger_test.Println("testCacheSize[MiB]:", MiBNumber)
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

			hashr.Write(buf[:re])
			rbufsz += re

			if rbufsz == testCacheSize {
				return
			}
		}

	}(server)

	for i := 0; i < MiBNumber/100; i++ {
		SendData = GetRDBytes_local(100 * 1024 * 1024)
		_, err := gstunnellib.NetConnWriteAll(server, SendData)
		checkError(err)
		_, err = hashw.Write(SendData)
		checkError(err)
	}
	server.Close()
	logger_test.Printf("Write done time: %s", time.Since(t1))
	//	forceGC()
	//time.Sleep(time.Second * 6)
	wg.Wait()
	//	forceGC()
	//time.Sleep(time.Second * 60)

	if !bytes.Equal(hashr.Sum(nil), hashw.Sum(nil)) {
		t.Fatal("Error: hashr != hashw")
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] end.")
	t2 := time.Now()
	logger_test.Println("testCacheSize[MiB]:", MiBNumber)
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/s]:", float64(MiBNumber)/t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()

	//pprof.StopCPUProfile()
}

func Test_server_NetPipe_nGiB_hash_random(t *testing.T) {

	//samplesNum := 100
	//f, err := os.Create(fmt.Sprintf("cpuprof-%d.prof", samplesNum))
	//checkError_panic(err)
	//defer f.Close()
	//runtime.SetCPUProfileRate(samplesNum)
	//pprof.StartCPUProfile(f)

	//debug.SetGCPercent(-1)
	//debug.SetMemoryLimit(1024 * 1024 * 1024 * 10)

	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}

	if g_count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because g_count_inTest_server_NetPipe > 0.")
		return
	}

	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	rawClient := gstestpipe.NewServiceServerNone()
	gstServer := gstestpipe.NewGstPiPoDefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	//MiBNumber := 1000
	testCacheSize := 0
	const maxBytes = 1024 * 1024
	const writeNumber int = 10000 * 1
	var randomSZ int64 = 0

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawClient.GetClientConn(), gstServer.GetConn(), wg_run)

	wg := sync.WaitGroup{}
	writeDone := atomic.Bool{}

	server := rawClient.GetServerConn()

	var SendData []byte
	//rbuf := make([]byte, 0, len(SendData))

	hashw := sha256.New()
	hashr := sha256.New()

	//rfd, err := os.Create(os.DevNull)
	//checkError_panic(err)
	//wfd, err := os.Create(os.DevNull)
	//checkError_panic(err)

	logger_test.Println("Test_server_NetPipe_nGiB data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		//defer rfd.Close()

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
			rbuf := buf[:re]

			//fdn, err := rfd.Write(rbuf)
			//checkError_panic(err)
			//if fdn != len(rbuf) {
			//	panic("fdn!=len(rbuf)")
			//}

			hashr.Write(rbuf)
			rbufsz += re

			if writeDone.Load() {
				if rbufsz == testCacheSize {
					return
				}
			}
		}

	}(server)

	for i := 0; i < writeNumber; i++ {
		randomSZ = gsrand.GetRDCInt_max(int64(maxBytes))
		//	fmt.Println("random sz:", randomSZ, float64(randomSZ)/1024/1024)
		SendData = GetRDBytes_local(int(randomSZ))
		wlen, err := gstunnellib.NetConnWriteAll(server, SendData)
		checkError(err)

		//fdn, err := wfd.Write(SendData)
		//checkError_panic(err)
		//if fdn != len(SendData) {
		//	panic("fdn != len(SendData)")
		//}

		testCacheSize += int(wlen)
		_, err = hashw.Write(SendData)
		checkError(err)
	}
	//wfd.Close()
	writeDone.Store(true)

	go func() {
		time.Sleep(time.Second)
		server.Close()
	}()
	logger_test.Printf("Write done time: %s", time.Since(t1))
	//	forceGC()
	//time.Sleep(time.Second * 6)
	wg.Wait()
	server.Close()
	//	forceGC()
	//time.Sleep(time.Second * 60)

	if !bytes.Equal(hashr.Sum(nil), hashw.Sum(nil)) {
		t.Fatal("Error: hashr != hashw")
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] end.")
	t2 := time.Now()
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/s]:", float64(testCacheSize)/1024/1024/t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()

	//pprof.StopCPUProfile()
}

func Test_server_NetPipe_nGiB_hash_diZeng(t *testing.T) {

	//samplesNum := 100
	//f, err := os.Create(fmt.Sprintf("cpuprof-%d.prof", samplesNum))
	//checkError_panic(err)
	//defer f.Close()
	//runtime.SetCPUProfileRate(samplesNum)
	//pprof.StartCPUProfile(f)

	//debug.SetGCPercent(-1)
	//debug.SetMemoryLimit(1024 * 1024 * 1024 * 10)

	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}

	if g_count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because g_count_inTest_server_NetPipe > 0.")
		return
	}

	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	rawClient := gstestpipe.NewServiceServerNone()
	gstServer := gstestpipe.NewGstPiPoDefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	//MiBNumber := 1000
	testCacheSize := 0
	//const maxBytes = 1024
	const writeNumber int = 10000 * 10
	var randomSZ int64 = 0

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawClient.GetClientConn(), gstServer.GetConn(), wg_run)

	wg := sync.WaitGroup{}
	writeDone := atomic.Bool{}

	server := rawClient.GetServerConn()

	var SendData []byte
	//rbuf := make([]byte, 0, len(SendData))

	hashw := sha256.New()
	hashr := sha256.New()

	//rfd, err := os.Create(os.DevNull)
	//checkError_panic(err)
	//wfd, err := os.Create(os.DevNull)
	//checkError_panic(err)

	logger_test.Println("Test_server_NetPipe_nGiB data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		//defer rfd.Close()

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
			rbuf := buf[:re]

			//fdn, err := rfd.Write(rbuf)
			//checkError_panic(err)
			//if fdn != len(rbuf) {
			//	panic("fdn!=len(rbuf)")
			//}

			hashr.Write(rbuf)
			rbufsz += re

			if writeDone.Load() {
				if rbufsz == testCacheSize {
					return
				}
			}
		}

	}(server)

	for i := 0; i < writeNumber; i++ {
		randomSZ = int64(i)
		if i%10000 == 0 {
			fmt.Println("random sz:", randomSZ, float64(randomSZ)/1024/1024)
		}
		SendData = GetRDBytes_local(int(randomSZ))
		wlen, err := gstunnellib.NetConnWriteAll(server, SendData)
		checkError(err)

		//fdn, err := wfd.Write(SendData)
		//checkError_panic(err)
		//if fdn != len(SendData) {
		//	panic("fdn != len(SendData)")
		//}

		testCacheSize += int(wlen)
		_, err = hashw.Write(SendData)
		checkError(err)
	}
	//wfd.Close()
	writeDone.Store(true)

	go func() {
		time.Sleep(time.Second)
		server.Close()
	}()
	logger_test.Printf("Write done time: %s", time.Since(t1))
	//	forceGC()
	//time.Sleep(time.Second * 6)
	wg.Wait()
	server.Close()
	//	forceGC()
	//time.Sleep(time.Second * 60)

	if !bytes.Equal(hashr.Sum(nil), hashw.Sum(nil)) {
		t.Fatal("Error: hashr != hashw")
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] end.")
	t2 := time.Now()
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/s]:", float64(testCacheSize)/1024/1024/t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()

	//pprof.StopCPUProfile()
}
