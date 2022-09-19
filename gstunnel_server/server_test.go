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

// var logger_test log.Logger = *gstunnellib.CreateFileLogger("server_test.log")
var logger_test *log.Logger

func init() {
	logger_test = Logger
}

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
	gsc := gstestpipe.NewGsClientDefultKey()

	networkTimeout = time.Minute * 20
	Mt_model = mt_mode
	debug_server = true

	testReadTimeOut := time.Second * 1
	testCacheSize := 300 * 1024 * 1024

	run_pipe_test(ss, gsc)

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
			t.Logf("server read len: %d", re)
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
}

func inTest_server_NetPipe_m(t *testing.T) {
	logger_test.Println("[inTest_server_NetPipe] start.")
	ss := gstestpipe.NewServiceServerNone()
	gsc := gstestpipe.NewGsClientDefultKey()

	networkTimeout = time.Minute * 20
	Mt_model = true
	debug_server = true

	testReadTimeOut := time.Second * 1
	testCacheSize := 500 * 1024 * 1024

	run_pipe_test(ss, gsc)

	wg := sync.WaitGroup{}

	server := ss.GetServerConn()

	//SendData := []byte("123456")
	//rbuff := bytes.Buffer{}
	logger_test.Println("testCacheSize[MiB]:", testCacheSize/1024/1024)

	logger_test.Println("inTest_server_NetPipe data transfer start.")
	t1 := time.Now()

	for i := 0; i < 10; i++ {
		SendData := GetRDBytes_local(testCacheSize)
		rbuf := make([]byte, 0, testCacheSize)

		wg.Add(1)
		go func(server net.Conn) {
			defer wg.Done()
			buf := make([]byte, net_read_size)

			for {
				server.SetReadDeadline(time.Now().Add(testReadTimeOut))
				re, err := server.Read(buf)
				t.Logf("server read len: %d", re)
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
		//server.Close()
		wg.Wait()
		server.Close()
		forceGC()

		if !bytes.Equal(SendData, rbuf) {
			t.Fatal("Error: SendData != rbuf.")
		}
	}

	logger_test.Println("[inTest_server_NetPipe] end.")
	t2 := time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	//forceGC()

	//time.Sleep(time.Second * 60)
}

func inTest_server_NetPipe_go_init() {
	networkTimeout = time.Minute * 20
	Mt_model = true
	debug_server = true
}

func inTest_server_NetPipe_go(t *testing.T, gwg *sync.WaitGroup) {
	defer gwg.Done()

	logger_test.Println("[inTest_server_NetPipe] start.")
	ss := gstestpipe.NewServiceServerNone()
	gsc := gstestpipe.NewGsClientDefultKey()

	testReadTimeOut := time.Second * 1
	testCacheSize := 100 * 1024 * 1024

	run_pipe_test(ss, gsc)

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
			t.Logf("server read len: %d", re)
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
}

func Test_server_NetPipe(t *testing.T) {
	defer gstunnellib.RunTimeDebugInfoV1.WriteFile("debugInfo.json")
	logger_test.Println("[Test_server_NetPipe] start.")
	inTest_server_NetPipe(t, false)
	logger_test.Println("[Test_server_NetPipe] end.")
}
func Test_server_NetPipe_st(t *testing.T) {
	defer gstunnellib.RunTimeDebugInfoV1.WriteFile("debugInfo.json")
	logger_test.Println("[Test_server_NetPipe_st] start.")
	inTest_server_NetPipe(t, false)
	logger_test.Println("[Test_server_NetPipe_st] end.")
}
func Test_server_NetPipe_mt(t *testing.T) {
	defer gstunnellib.RunTimeDebugInfoV1.WriteFile("debugInfo.json")
	logger_test.Println("[Test_server_NetPipe_mt] start.")
	inTest_server_NetPipe(t, true)
	logger_test.Println("[Test_server_NetPipe_mt] end.")
}

func Test_server_NetPipe_m(t *testing.T) {
	defer gstunnellib.RunTimeDebugInfoV1.WriteFile("debugInfo.json")
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
		inTest_server_NetPipe(t, false)
		forceGC()
	}
	logger_test.Println("[Test_server_NetPipe_loop] end.")
}
