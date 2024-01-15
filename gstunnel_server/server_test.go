package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

func init() {
	//runtime.SetCPUProfileRate(200)
	debug.SetMemoryLimit(1024 * 1024 * 3000)
	init_server_test()
	logger_test = g_logger
	g_networkTimeout = time.Second * 10
}

func sleep_race() {
	if gsbase.GetRaceState() {
		time.Sleep(time.Second * 2)
	}
}

func Test_server_NetPipe_st(t *testing.T) {
	defer sleep_race()
	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Print("\n\n")
	logger_test.Println("[Test_server_NetPipe_st] start.")
	inTest_server_NetPipe(t, false)
	logger_test.Print("[Test_server_NetPipe_st] end.\n\n")
}
func Test_server_NetPipe_mt(t *testing.T) {
	defer sleep_race()
	g_count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_mt] start.")
	inTest_server_NetPipe(t, true)
	logger_test.Print("[Test_server_NetPipe_mt] end.\n\n")
}

func Test_server_NetPipe_m(t *testing.T) {
	defer sleep_race()
	g_count_inTest_server_NetPipe++

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
	defer sleep_race()
	logger_test.Println("[Test_server_NetPipe_loop] start.")
	for i := 0; i < 6; i++ {
		logger_test.Printf("loop count: %d", i)
		inTest_server_NetPipe(t, false)
		forceGC()
	}
	logger_test.Print("[Test_server_NetPipe_loop] end.\n\n")
}

func Test_server_NetPipe_errorData(t *testing.T) {
	defer sleep_race()
	g_count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[notest_Test_server_NetPipe_errorData] start.")
	inTest_server_NetPipe_errorData(t, false)
	logger_test.Println("[notest_Test_server_NetPipe_errorData] end.")
}

func Test_server_NetPipe_errorKey(t *testing.T) {
	defer sleep_race()
	g_count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_NetPipe_errorKey] start.")
	inTest_server_NetPipe_errorKey(t, false)
	logger_test.Print("[Test_server_NetPipe_errorKey] end.\n\n")
}

func Test_server_timeout(t *testing.T) {
	defer sleep_race()
	g_count_inTest_server_NetPipe++

	if gstunnellib.G_RunTime_Debug {
		defer gstunnellib.G_RunTimeDebugInfo1.WriteFile("debugInfo.out.json")
	}
	logger_test.Println("[Test_server_timeout] start.")
	inTest_server_timeout(t, false)
	logger_test.Print("[Test_server_timeout] end.\n\n")
}

func Test_server_timeout2(t *testing.T) {
	defer sleep_race()
	g_count_inTest_server_NetPipe++

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
	//defer sleep_race()
	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}
	if g_count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because g_count_inTest_server_NetPipe > 0.")
		return
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1

	var MiBNumber int
	if !gsbase.GetRaceState() {
		MiBNumber = 512 * 2
	} else {
		MiBNumber = 512
	}
	testCacheSize := MiBNumber * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := rawServer.GetServerConn()

	SendData := GetRDBytes_local(testCacheSize)
	rbuf := make([]byte, 0, len(SendData))

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
			if gstunnellib.IsErrorNetUsually(err) {
				checkError_info(err)
				return
			} else {
				gstunnellib.CheckError_test(err, t)
			}

			rbuf = append(rbuf, buf[:re]...)
			rbufsz += re

			if rbufsz == len(SendData) {
				return
			}
		}

	}(server)

	_, err := gstunnellib.NetConnWriteAll(server, SendData)
	checkError(err)
	forceGC()
	//time.Sleep(time.Second * 6)
	wg.Wait()
	server.Close()
	forceGC()
	//time.Sleep(time.Second * 60)

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/s]:", float64(MiBNumber)/t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()
}

// 适合测试gstunnel是否正确传输数据。
// 但是因为计算数据哈希值，增加了cpu占用，导致吞吐量降低，所有不适合测试最大吞吐量。
func Test_server_NetPipe_nGiB_hash(t *testing.T) {
	//defer sleep_race()
	//samplesNum := 100
	//f, err := os.Create(fmt.Sprintf("cpuprof-%d.prof", samplesNum))
	//checkError_panic(err)
	//defer f.Close()
	//runtime.SetCPUProfileRate(samplesNum)
	//pprof.StartCPUProfile(f)

	//debug.SetGCPercent(-1)

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
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	MiBNumber := 100 * 100 / 2
	testCacheSize := MiBNumber * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := rawServer.GetServerConn()

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
			if gstunnellib.IsErrorNetUsually(err) {
				checkError_info(err)
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
		n, err := gstunnellib.NetConnWriteAll(server, SendData)
		checkError(err)
		_, err = hashw.Write(SendData[:n])
		checkError(err)
	}
	wg.Wait()
	server.Close()
	logger_test.Printf("Write done time: %s", time.Since(t1))
	//	forceGC()
	//time.Sleep(time.Second * 6)
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
	//defer sleep_race()
	//samplesNum := 100
	//f, err := os.Create(fmt.Sprintf("cpuprof-%d.prof", samplesNum))
	//checkError_panic(err)
	//defer f.Close()
	//runtime.SetCPUProfileRate(samplesNum)
	//pprof.StartCPUProfile(f)

	//debug.SetGCPercent(-1)

	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}
	if g_count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because g_count_inTest_server_NetPipe > 0.")
		return
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	//MiBNumber := 1000
	testCacheSize := 0
	const maxBytes = 1024 * 1024
	const writeNumber int = 10000 * 1
	var randomSZ int64 = 0

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	wg := sync.WaitGroup{}
	writeDone := atomic.Bool{}

	server := rawServer.GetServerConn()

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
			if gstunnellib.IsErrorNetUsually(err) {
				checkError_info(err)
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
		randomSZ = gsrand.GetRDInt_max(int64(maxBytes))
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
	//defer sleep_race()
	//samplesNum := 100
	//f, err := os.Create(fmt.Sprintf("cpuprof-%d.prof", samplesNum))
	//checkError_panic(err)
	//defer f.Close()
	//runtime.SetCPUProfileRate(samplesNum)
	//pprof.StartCPUProfile(f)

	//debug.SetGCPercent(-1)

	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}
	if g_count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because g_count_inTest_server_NetPipe > 0.")
		return
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	//MiBNumber := 1000
	testCacheSize := 0
	//const maxBytes = 1024
	const writeNumber int = 10000 * 10
	var randomSZ int64 = 0

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	wg := sync.WaitGroup{}
	writeDone := atomic.Bool{}

	server := rawServer.GetServerConn()

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
			if gstunnellib.IsErrorNetUsually(err) {
				checkError_info(err)
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

func in_Test_server_NetPipe_nGiB_wg(t *testing.T, inwg *sync.WaitGroup) {
	//defer sleep_race()
	defer inwg.Done()
	if gsbase.GetRaceState() {
		logger_test.Println("[Test_server_NetPipe_nGiB] Error: This func does not run because race is true.")
		return
	}
	if g_count_inTest_server_NetPipe > 0 {
		logger_test.Println("[Test_server_NetPipe_nGiB] This func does not run because g_count_inTest_server_NetPipe > 0.")
		return
	}
	logger_test.Println("[Test_server_NetPipe_nGiB] start.")
	//rawsever
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	g_Mt_model = true
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 1
	MiBNumber := 512
	testCacheSize := MiBNumber * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	//rawsever-server
	server := rawServer.GetServerConn()

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
			if gstunnellib.IsErrorNetUsually(err) {
				checkError_info(err)
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

	_, err := gstunnellib.NetConnWriteAll(server, SendData)
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
	//defer sleep_race()
	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go in_Test_server_NetPipe_nGiB_wg(t, &wg)
	}
	wg.Wait()
}
