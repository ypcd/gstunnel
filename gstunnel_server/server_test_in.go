package main

import (
	"bytes"
	randc "crypto/rand"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

var logger_test *log.Logger

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
	data := make([]byte, byteLen)
	_, err := randc.Reader.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

func forceGC() {
	logger_test.Println("[runtime.GC] start.")
	runtime.GC()
	logger_test.Println("[runtime.GC] end.")
}

func inTest_server_NetPipe(t *testing.T, mt_mode bool) {
	logger_test.Println("[inTest_server_NetPipe] start.")
	ss := gstestpipe.NewServiceServerNone()
	gsc := gstestpipe.NewGsPiPeDefultKey()

	Mt_model = mt_mode
	GValues.SetDebug(true)

	testReadTimeOut := time.Second * 1
	testCacheSize := 200 * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(ss.GetClientConn(), gsc.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := ss.GetServerConn()

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
		buf := make([]byte, net_read_size)

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

			rbuf = append(rbuf, buf[:re]...)

			if len(rbuf) == len(SendData) {
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

func inTest_server_NetPipe_errorData(t *testing.T, mt_mode bool) {
	logger_test.Println("[inTest_server_NetPipe_errorData] start.")
	ss := gstestpipe.NewServiceServerNone()
	gsc := gstestpipe.NewGsPiPeErrorKeyNoKey()

	Mt_model = mt_mode
	GValues.SetDebug(true)

	//testReadTimeOut := time.Second * 1
	testCacheSize := 200 * 1024 * 1024

	wg := new(sync.WaitGroup)
	run_pipe_test_wg(ss.GetClientConn(), gsc.GetServerConn(), wg)

	//wg := sync.WaitGroup{}

	//server := ss.GetServerConn()

	SendData := GetRDBytes_local(1024 * 1024 * 10)
	//SendData := []byte("123456")
	//rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_server_NetPipe data transfer start.")

	///////////////////////////////////////////////////////////////////////////
	_, err := io.Copy(gsc.GetClientConn(), bytes.NewBuffer(SendData))
	checkError_NoExit(err)

	//////////////////////////////////////////////////////////////////////////
	wg.Wait()
}

func inTest_server_NetPipe_errorKey(t *testing.T, mt_mode bool) {
	logger_test.Println("[inTest_server_NetPipe_errorKey] start.")
	ss := gstestpipe.NewServiceServerNone()
	key_error := "123456789012345678901234567890aa"
	gsc := gstestpipe.NewGsPiPeErrorKeyNoKey()

	Mt_model = mt_mode
	GValues.SetDebug(true)

	//testReadTimeOut := time.Second * 1
	testCacheSize := 200 * 1024 * 1024

	wg := new(sync.WaitGroup)
	run_pipe_test_wg(ss.GetClientConn(), gsc.GetServerConn(), wg)

	//wg := sync.WaitGroup{}

	//server := ss.GetServerConn()

	//SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	//rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_server_NetPipe_errorKey data transfer start.")

	///////////////////////////////////////////////////////////////////////////
	pack1 := gstunnellib.NewGsPackNet(key_error)
	pdata := pack1.Packing(gsrand.GetRDBytes(50 * 1024))
	_, err := io.Copy(gsc.GetClientConn(), bytes.NewBuffer(pdata))
	checkError_panic(err)

	wg.Wait()
	//////////////////////////////////////////////////////////////////////////

}

func inTest_server_NetPipe_go_init() {

	Mt_model = true
	GValues.SetDebug(true)
}

func inTest_server_NetPipe_go(t *testing.T, gwg *sync.WaitGroup) {
	defer gwg.Done()

	logger_test.Println("[inTest_server_NetPipe] start.")
	ss := gstestpipe.NewServiceServerNone()
	gsc := gstestpipe.NewGsPiPeDefultKey()

	testReadTimeOut := time.Second * 1
	testCacheSize := 100 * 1024 * 1024

	wg_run := new(sync.WaitGroup)
	run_pipe_test_wg(ss.GetClientConn(), gsc.GetConn(), wg_run)

	wg := sync.WaitGroup{}

	server := ss.GetServerConn()

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
		buf := make([]byte, net_read_size)

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

			rbuf = append(rbuf, buf[:re]...)
		}

	}(server)

	_, err := io.Copy(server, bytes.NewBuffer(SendData))
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
	ss := gstestpipe.NewServiceServerNone()
	gsc := gstestpipe.NewGsPiPeDefultKey()

	networkTimeout = time.Second * 1
	Mt_model = mt_mode
	GValues.SetDebug(true)

	//	testReadTimeOut := time.Second * 1
	//	testCacheSize := 200 * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg(ss.GetClientConn(), gsc.GetConn(), wg_run)

	time.Sleep(time.Second * 2)
	logger_test.Println("Func done.")
}
