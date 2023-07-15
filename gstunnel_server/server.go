/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"errors"
	"fmt"
	"net"
	_ "net/http/pprof"
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

func main() {
	init_server_run()
	for {
		run()
	}
}

func createToRawServiceConnHandler(chanconnlist <-chan net.Conn, rawServiceAddr string) {

	const maxTryConnService int = 5000
	var err error
	for {

		gstClient, ok := <-chanconnlist
		if !ok {
			checkError_NoExit(errors.New("'gstServer, ok := <-chanconnlist' is error"))
			return
		}
		g_log_List.GSIpLogger.Printf("ip: %s\n", gstClient.RemoteAddr().String())

		connServiceError_count := 0
		var rawService net.Conn
		for {
			rawService, err = net.Dial("tcp", rawServiceAddr)
			checkError_NoExit(err)
			connServiceError_count += 1
			if err == nil {
				break
			}
			if connServiceError_count > maxTryConnService {
				checkError(
					fmt.Errorf("connService_count > maxTryConnService(%d)", maxTryConnService))
			}
			//g_Logger.Println("conn.")
		}

		gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gstst)
		g_gstst.GetStatusConnList().Add(gctx.GetGsId(), gstClient, rawService)

		go srcTOdstUn(gstClient, rawService, gctx)
		go srcTOdstP(rawService, gstClient, gctx)
		g_Logger.Printf("go [%d].\n", gctx.GetGsId())
	}
}

func run() {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	var lstnaddr, rawServiceAddr string

	lstnaddr = g_gsconfig.Listen
	rawServiceAddr = g_gsconfig.GetServer_rand()

	g_Logger.Println("Listen_Addr:", lstnaddr)
	g_Logger.Println("Conn_Addr:", rawServiceAddr)
	g_Logger.Println("Begin......")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", lstnaddr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//var rawService net.Conn

	chanConnList := make(chan net.Conn, 1024)
	defer close(chanConnList)

	for i := 0; i < g_connHandleGoNum; i++ {
		go createToRawServiceConnHandler(chanConnList, rawServiceAddr)
	}

	for {
		gstClient, err := listener.Accept()
		if err != nil {
			checkError_NoExit(err)
			continue
		}
		chanConnList <- gstClient
	}
}

// 性能是old版的4倍
func run_pipe_test_listen(lstnAddr, rawServiceAddr string) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	//var lstnaddr, rawServiceAddr string

	//lstnaddr = g_gsconfig.Listen
	//rawServiceAddr = g_gsconfig.GetServer_rand()

	g_Logger.Println("Listen_Addr:", lstnAddr)
	g_Logger.Println("Conn_Addr:", rawServiceAddr)
	g_Logger.Println("Begin......")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", lstnAddr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//var rawService net.Conn

	chanConnList := make(chan net.Conn, 1024)

	for i := 0; i < g_connHandleGoNum; i++ {
		go createToRawServiceConnHandler(chanConnList, rawServiceAddr)
	}

	for {
		gstClient, err := listener.Accept()
		if err != nil {
			checkError_NoExit(err)
			continue
		}
		chanConnList <- gstClient
	}
}

func run_pipe_test_listen_old(lstnAddr, rawServiceAddr string) {
	//defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	g_Logger.Println("Listen_Addr:", lstnAddr)
	g_Logger.Println("Conn_Addr:", rawServiceAddr)
	g_Logger.Println("Begin......")

	service := lstnAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		gstServer, err := listener.Accept()
		if err != nil {
			checkError_NoExit(err)
			continue
		}
		g_log_List.GSIpLogger.Printf("ip: %s\n", gstServer.RemoteAddr().String())

		service := rawServiceAddr

		connServiceError_count := 0
		const maxTryConnService int = 5000
		var dst net.Conn
		for {
			dst, err = net.Dial("tcp", service)
			checkError_NoExit(err)
			connServiceError_count += 1
			if err == nil {
				break
			}
			if connServiceError_count > maxTryConnService {
				checkError(
					errors.New(
						fmt.Sprintf("connService_count > maxTryConnService(%d)", maxTryConnService),
					),
				)
			}
			//Logger.Println("conn.")
		}

		gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gstst)
		g_gstst.GetStatusConnList().Add(gctx.GetGsId(), gstServer, dst)

		go srcTOdstUn(gstServer, dst, gctx)
		go srcTOdstP(dst, gstServer, gctx)
		g_Logger.Printf("go [%d].\n", gctx.GetGsId())
	}
}

func run_pipe_test_wg(dst net.Conn, gstServer net.Conn, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	g_Logger.Println("Test_Mt_model:", g_Mt_model)
	g_log_List.GSIpLogger.Printf("ip: %s\n", gstServer.RemoteAddr().String())

	gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gstst)
	g_gstst.GetStatusConnList().Add(gctx.GetGsId(), gstServer, dst)

	wg.Add(2)
	go srcTOdstUn_wg(gstServer, dst, wg, gctx)
	go srcTOdstP_wg(dst, gstServer, wg, gctx)
	g_Logger.Println("Gstunnel go.")

}

func IsTheVersionConsistent_send(dst net.Conn, apack gstunnellib.GsPack, wlent *int64) error {
	return gstunnellib.IsTheVersionConsistent_send(dst, apack, wlent)
}

func ChangeCryKey_send(dst net.Conn, apack gstunnellib.GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return gstunnellib.ChangeCryKey_send(dst, apack, ChangeCryKey_Total, wlent)
}

// service to gstunnel client
func srcTOdstP(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	if g_Mt_model {
		srcTOdstP_mt(src, dst, gctx)
	} else {
		srcTOdstP_st(src, dst, gctx)
	}
}

// gstunnel client to service
func srcTOdstUn(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	if g_Mt_model {
		srcTOdstUn_mt(src, dst, gctx)
	} else {
		srcTOdstUn_st(src, dst, gctx)
	}
}

// service to gstunnel client
func srcTOdstP_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(src, dst, gctx)
	} else {
		srcTOdstP_st(src, dst, gctx)
	}
}

// gstunnel client to service
func srcTOdstUn_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(src, dst, gctx)
	} else {
		srcTOdstUn_st(src, dst, gctx)
	}
}
