package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

var logger_test *log.Logger
var g_count_inTest_server_NetPipe int = 0

func Test_func1(t *testing.T) {
	logger_test.Println("ok.")
	v1 := []byte{}
	var v2 []byte = nil

	_, _ = v1, v2
	logger_test.Println(bytes.Equal(v1, v2))
	logger_test.Println(v1 == nil)
	logger_test.Println(v2 == nil)

	v3 := []byte("123456")
	v4 := v3[6:]
	_ = v4
}

func GetRDBytes_local(byteLen int) []byte {
	return gsrand.GetRDBytes(byteLen)
}

func forceGC() {
	//logger_test.Println("[runtime.GC] start.")
	//runtime.GC()
	//logger_test.Println("[runtime.GC] end.")
}

func inTest_server_NetPipe(t *testing.T, mt_mode bool) {
	g_count_inTest_server_NetPipe++

	logger_test.Println("[inTest_server_NetPipe] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	g_Mt_model = mt_mode
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 2
	testCacheSizeMiB := 200
	testCacheSize := testCacheSizeMiB * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := rawServer.GetServerConn()

	SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSizeMiB)

	logger_test.Println("inTest_server_NetPipe data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, g_net_read_size)

		for {
			server.SetReadDeadline(time.Now().Add(testReadTimeOut))
			re, err := server.Read(buf)
			//t.Logf("server read len: %d", re)
			if gstunnellib.IsErrorNetUsually(err) {
				if len(rbuf) == len(SendData) {
					checkError_info(err)
					return
				} else {
					panic(err)
				}
			} else {
				gstunnellib.CheckError_test(err, t)
			}

			rbuf = append(rbuf, buf[:re]...)

			if len(rbuf) == len(SendData) {
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
		panic("Error: SendData != rbuf.")
	}
	logger_test.Println("[inTest_server_NetPipe] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/sec]:", float64(testCacheSizeMiB)/(t2.Sub(t1).Seconds()))
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()
}

func inTest_server_NetPipe_errorData(t *testing.T, mt_mode bool) {
	defer gserror.Panic_Recover(g_logger)
	logger_test.Println("[inTest_server_NetPipe_errorData] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGstRSAPiPoErrorKeyNoKey()

	g_Mt_model = mt_mode
	g_Values.SetDebug(true)

	//testReadTimeOut := time.Second * 1
	testCacheSize := 200 * 1024 * 1024

	wg := new(sync.WaitGroup)
	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetServerConn(), wg)

	//wg := sync.WaitGroup{}

	//server := rawServer.GetServerConn()

	SendData := GetRDBytes_local(1024 * 1024 * 10)
	//SendData := []byte("123456")
	//rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_server_NetPipe data transfer start.")

	///////////////////////////////////////////////////////////////////////////
	n, err := gstunnellib.NetConnWriteAll(gstClient.GetClientConn(), SendData)
	_ = n
	checkError_NoExit(err)

	//////////////////////////////////////////////////////////////////////////
	wg.Wait()
}

func in2Test_serviceNone(rawServer gstestpipe.IRawdataPiPe) {
	//rawServer := NewServiceServerNone()
	wg := sync.WaitGroup{}

	SendData := []byte("a")
	rbuf := make([]byte, len(SendData))

	server := rawServer.GetServerConn()

	wg.Add(1)
	go func(client net.Conn) {
		defer wg.Done()
		//buf := make([]byte, 1024*1024)
		_, err := io.ReadAtLeast(client, rbuf, len(rbuf))
		//t.Logf("server read len: %d", re)
		checkError_info(err)

	}(rawServer.GetClientConn())

	server.Write(SendData)
	wg.Wait()

	if !bytes.Equal(SendData, rbuf) {
		panic("Error: SendData != rbuf.")
	}

}

func inTest_serviceNone(rawServer gstestpipe.IRawdataPiPe) {
	go in2Test_serviceNone(rawServer)
}

func inTest_server_NetPipe_errorKey(t *testing.T, mt_mode bool) {
	logger_test.Println("[inTest_server_NetPipe_errorKey] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGstRSAPiPoErrorKeyNoKey()

	t1 := time.Now()
	rawServerConn := rawServer.GetClientConn()

	rawServerDebug := gstunnellib.NewGSTNetConnDebug(rawServerConn)
	rawServerConn.SetDeadline(time.Now().Add(time.Second * 60))

	g_networkTimeout_old := g_networkTimeout
	g_networkTimeout = time.Second * 3
	defer func() { g_networkTimeout = g_networkTimeout_old }()

	g_Mt_model = mt_mode
	g_Values.SetDebug(true)

	//testReadTimeOut := time.Second * 1
	//testCacheSize := 200 * 1024 * 1024
	//inTest_serviceNone(rawServer)

	wg := new(sync.WaitGroup)
	gstunnellib.IsClosedPanic(rawServerConn)
	run_pipe_test_wg(rawServerDebug, gstClient.GetServerConn(), wg)
	gstunnellib.IsClosedPanic(rawServerConn)

	fmt.Println("time:", time.Since(t1))
	//wg := sync.WaitGroup{}
	//inTest_serviceNone(rawServer)

	//server := rawServer.GetServerConn()

	//SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	//rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	//logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_server_NetPipe_errorKey data transfer start.")

	///////////////////////////////////////////////////////////////////////////
	//等待交换密钥完成，再写入错误数据。

	time.Sleep(time.Millisecond * 100)
	key_error := gsrand.GetrandStringPlus(gsbase.G_AesKeyLen)
	pack1 := gstunnellib.NewGsPackNet(key_error)
	pdata := pack1.Packing(gsrand.GetRDBytes(50 * 1024))
	_, err := gstunnellib.NetConnWriteAll(gstClient.GetClientConn(), pdata)
	checkError_panic(err)
	gstunnellib.IsClosedPanic(rawServerConn)

	wg.Wait()
	//time.Sleep(13 * time.Second)
	//////////////////////////////////////////////////////////////////////////
	//time.Sleep(time.Millisecond * 2000)
	gstunnellib.IsClosedPanic(rawServerConn)
	logger_test.Println("inTest_server_NetPipe_errorKey data transfer end.")

}

func inTest_server_NetPipe_go_init() {

	g_Mt_model = true
	g_Values.SetDebug(true)
}

func inTest_server_NetPipe_go(t *testing.T, gwg *sync.WaitGroup) {
	defer gwg.Done()

	logger_test.Println("[inTest_server_NetPipe] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	testReadTimeOut := time.Second * 1
	testCacheSize := 100 * 1024 * 1024

	wg_run := new(sync.WaitGroup)
	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := rawServer.GetServerConn()

	SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_server_NetPipe data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, g_net_read_size)

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
		}

	}(server)

	_, err := gstunnellib.NetConnWriteAll(server, SendData)
	checkError(err)
	forceGC()
	//time.Sleep(time.Second * 6)
	wg.Wait()
	server.Close()
	forceGC()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}
	logger_test.Println("[inTest_server_NetPipe] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()
}

func inTest_server_timeout(t *testing.T, mt_mode bool) {
	logger_test.Println("[inTest_server_NetPipe] start.")
	rawServer := gstestpipe.NewServiceServerNone()
	gstClient := gstestpipe.NewGSTRSAClientImp_DefaultKey()

	old_networkTimeout := g_networkTimeout
	g_networkTimeout = time.Second * 1
	defer func() {
		g_networkTimeout = old_networkTimeout
	}()
	g_Mt_model = mt_mode
	g_Values.SetDebug(true)

	//	testReadTimeOut := time.Second * 1
	//	testCacheSize := 200 * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(rawServer.GetClientConn(), gstClient.GetConn(), wg_run)

	time.Sleep(time.Second * 2)
	logger_test.Println("Func done.")
}
