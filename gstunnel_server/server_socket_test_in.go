package main

import (
	"bytes"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

type testSocketEchoHandlerObj struct {
	testReadTimeOut  time.Duration
	testCacheSize    int
	testCacheSizeMiB int
	wg               *sync.WaitGroup
}

func test_server_socket_echo_handler(obj *testSocketEchoHandlerObj, server net.Conn) {
	defer func() {
		if obj.wg != nil {
			obj.wg.Done()
		}
	}()

	testReadTimeOut := obj.testReadTimeOut
	testCacheSize := obj.testCacheSize
	//testCacheSizeMiB := obj.testCacheSizeMiB

	SendData := GetRDBytes_local(testCacheSize)
	//SendData := []byte("123456")
	rbuf := make([]byte, 0, len(SendData))
	//rbuff := bytes.Buffer{}
	//logger_test.Println("testCacheSize[MiB]:", testCacheSizeMiB)

	//logger_test.Println("inTest_server_NetPipe data transfer start.")
	//t1 := time.Now()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, g_net_read_size)

		for {
			server.SetReadDeadline(time.Now().Add(testReadTimeOut))
			re, err := server.Read(buf)
			//t.Logf("server read len: %d", re)
			if gstunnellib.IsErrorNetUsually(err) {
				server.Close()
				//time.Sleep(time.Second * 2)
				checkError_panic(err)
				return
			} else {
				checkError_panic(err)
			}

			rbuf = append(rbuf, buf[:re]...)

			if len(rbuf) == len(SendData) {
				return
			}
		}

	}(server)

	//_, err := gstunnellib.NetConnWriteAll(server, SendData)
	_, err := gstunnellib.NetConnWriteAll(server, SendData)

	checkError(err)
	//time.Sleep(time.Second * 6)
	wg.Wait()
	server.Close()
	//time.Sleep(time.Second * 60)

	if !bytes.Equal(SendData, rbuf) {
		panic("Error: SendData != rbuf.")
	}
	//logger_test.Println("[inTest_server_NetPipe] end.")
	//t2 := time.Now()
	//logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	//logger_test.Println("[pipe run MiB/sec]:", float64(testCacheSizeMiB)/(t2.Sub(t1).Seconds()))

	//time.Sleep(time.Hour * 60)
}

func inTest_server_socket(t *testing.T, mt_mode bool, mbNum int) {
	g_count_inTest_server_NetPipe++

	var t1, t2 time.Time

	logger_test.Println("[inTest_server_NetPipe] start.")
	t1 = time.Now()
	rawServer := gstestpipe.NewRawServerSocket_RandAddr()
	defer rawServer.Close()
	rawServer.Run()

	listenAddr := gstunnellib.GetNetLocalRDPort()

	g_Mt_model = mt_mode
	g_Values.SetDebug(true)

	handleobj := testSocketEchoHandlerObj{
		testReadTimeOut:  time.Second * 10,
		testCacheSize:    mbNum * 1024 * 1024,
		testCacheSizeMiB: mbNum}
	//g_networkTimeout = handleobj.testReadTimeOut

	gstServerT := newGstServer(listenAddr, rawServer.GetServerAddr())
	defer gstServerT.close()
	go gstServerT.run()

	RSAclientKey := gsrsa.NewGenRSAObj(4096)
	tinit1 := time.Now()
	gstClient := gstestpipe.NewGstRSAClientSocketEcho_defaultKey(listenAddr, RSAclientKey)
	defer gstClient.Close()
	logger_test.Println("gstclient init time:", time.Since(tinit1))

	time.Sleep(time.Millisecond * 10)
	gstClient.Run()

	for i := 0; i < len(rawServer.GetConnList()); i++ {
		server := <-rawServer.GetConnList()
		test_server_socket_echo_handler(&handleobj, server)
	}
	t2 = time.Now()
	logger_test.Println("[pipe run time(sec)]:", t2.Sub(t1).Seconds())
	logger_test.Println("[pipe run MiB/sec]:", float64(handleobj.testCacheSizeMiB)/(t2.Sub(t1).Seconds()))
}

func inTest_server_socket_mt(t *testing.T, mt_mode bool, connes int) {
	g_count_inTest_server_NetPipe++

	mtNum := connes

	time_run_begin := time.Now()
	logger_test.Println("[inTest_server_socket_mt] start.")

	rawServer := gstestpipe.NewRawServerSocket_RandAddr()
	defer rawServer.Close()
	rawServer.Run()

	listenAddr := gstunnellib.GetNetLocalRDPort()

	gstServerT := newGstServer(listenAddr, rawServer.GetServerAddr())
	defer gstServerT.close()
	go gstServerT.run()

	RSAclientKey := gsrsa.NewGenRSAObj(4096)
	gstClient := gstestpipe.NewGstRSAClientSocketEcho_defaultKey(listenAddr, RSAclientKey)
	defer gstClient.Close()

	g_Mt_model = mt_mode
	g_Values.SetDebug(true)
	wg := sync.WaitGroup{}

	handleobj := testSocketEchoHandlerObj{
		testReadTimeOut:  time.Second * 60,
		testCacheSize:    2 * 1024 * 1024,
		testCacheSizeMiB: 2,
		wg:               &wg}
	g_networkTimeout = handleobj.testReadTimeOut

	time.Sleep(time.Millisecond * 10)
	time_connSocket_begin := time.Now()
	for i := 0; i < mtNum; i++ {
		gstClient.Run()
	}
	time_connSocket := time.Since(time_connSocket_begin)

	for i := 0; i < mtNum; i++ {
		server := <-rawServer.GetConnList()
		wg.Add(1)
		go test_server_socket_echo_handler(&handleobj, server)
	}
	wg.Wait()
	time_run_end := time.Now()
	os.Stdout.Sync()
	time.Sleep(time.Millisecond * 100)
	logger_test.Println("time_connSocket:", time_connSocket)
	logger_test.Println("time_run:", time_run_end.Sub(time_run_begin))
}

func inTest_server_socket_mt_old_listen(t *testing.T, mt_mode bool, connes int) {
	g_count_inTest_server_NetPipe++

	mtNum := connes

	time_run_begin := time.Now()
	logger_test.Println("[inTest_server_socket_mt_old] start.")

	rawServer := gstestpipe.NewRawServerSocket_RandAddr()
	defer rawServer.Close()
	rawServer.Run()

	listenAddr := gstunnellib.GetNetLocalRDPort()
	RSAclientKey := gsrsa.NewGenRSAObj(4096)
	gstClient := gstestpipe.NewGstRSAClientSocketEcho_defaultKey(listenAddr, RSAclientKey)
	defer gstClient.Close()

	g_Mt_model = mt_mode
	g_Values.SetDebug(true)
	wg := sync.WaitGroup{}

	handleobj := testSocketEchoHandlerObj{
		testReadTimeOut:  time.Second * 60,
		testCacheSize:    2 * 1024 * 1024,
		testCacheSizeMiB: 2,
		wg:               &wg}
	g_networkTimeout = handleobj.testReadTimeOut

	go run_pipe_test_listen_old(listenAddr, rawServer.GetServerAddr())

	time.Sleep(time.Millisecond * 10)
	time_connSocket_begin := time.Now()
	for i := 0; i < mtNum; i++ {
		gstClient.Run()
	}
	time_connSocket := time.Since(time_connSocket_begin)

	for i := 0; i < mtNum; i++ {
		server := <-rawServer.GetConnList()
		wg.Add(1)
		go test_server_socket_echo_handler(&handleobj, server)
	}
	wg.Wait()
	time_run_end := time.Now()
	os.Stdout.Sync()
	time.Sleep(time.Millisecond * 100)
	logger_test.Println("time_connSocket:", time_connSocket)
	logger_test.Println("time_run:", time_run_end.Sub(time_run_begin))
}
