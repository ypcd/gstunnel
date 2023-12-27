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

func inTest_client_NetPipe(t *testing.T, mt_mode bool) {
	logger_test.Println("[inTest_client_NetPipe] start.")
	rawClient := gstestpipe.NewSrcClientNone()
	gstServer := gstestpipe.NewGstPiPoDefaultKey()

	g_Mt_model = mt_mode
	g_Values.SetDebug(true)

	testReadTimeOut := time.Second * 100
	testCacheSize := 200 * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg_old(rawClient, gstServer, wg_run)

	wg := sync.WaitGroup{}

	server := rawClient.GetClientConn()

	SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_client_NetPipe data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, g_net_read_size)

		for {
			server.SetReadDeadline(time.Now().Add(testReadTimeOut))
			re, err := server.Read(buf)
			//t.Logf("server read len: %d", re)
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				gstunnellib.CheckErrorEx_info(err, g_Logger)
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
	logger_test.Println("[inTest_client_NetPipe] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()
}

func inTest_client_NetPipe_go_init() {

	g_Mt_model = true
	g_Values.SetDebug(true)
}

func inTest_client_NetPipe_go(t *testing.T, gwg *sync.WaitGroup) {
	defer gwg.Done()

	logger_test.Println("[inTest_client_NetPipe] start.")
	rawClient := gstestpipe.NewServiceServerNone()
	gstServer := gstestpipe.NewGstPiPoDefaultKey()

	testReadTimeOut := time.Second * 3
	testCacheSize := 100 * 1024 * 1024

	wg_run := new(sync.WaitGroup)

	run_pipe_test_wg_old(rawClient, gstServer, wg_run)

	wg := sync.WaitGroup{}

	server := rawClient.GetClientConn()

	SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_client_NetPipe data transfer start.")
	t1 := time.Now()

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, g_net_read_size)

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
	logger_test.Println("[inTest_client_NetPipe] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
	wg_run.Wait()
}

func inTest_client_timeout(t *testing.T, mt_mode bool) {
	logger_test.Println("[inTest_client_NetPipe] start.")
	rawClient := gstestpipe.NewSrcClientNone()
	gstServer := gstestpipe.NewGstPiPoDefaultKey()

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

	run_pipe_test_wg_old(rawClient, gstServer, wg_run)

	time.Sleep(time.Second * 2)
	logger_test.Println("Func done.")
}
