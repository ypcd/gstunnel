/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"errors"
	"net"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func main() {
	init_client_run()
	for {
		run()
	}
}

func createToGstServerConnHandler(chanconnlist <-chan net.Conn, gstServerAddr string) {

	const maxTryConnService int = 5000
	//var err error
	for {

		rawClient, ok := <-chanconnlist
		if !ok {
			checkError_NoExit(errors.New("'gstServer, ok := <-chanconnlist' is error"))
			return
		}
		g_log_List.GSIpLogger.Printf("ip: %s\n", rawClient.RemoteAddr().String())

		server_conn_error_total := 0
		tmr := timerm.CreateTimer(time.Second * 10)
		for {
			if tmr.Run() {
				g_Logger.Println("Error: server_conn_error_total: ", server_conn_error_total)
				tmr.Boot()
			}

			gstServer, err := net.Dial("tcp", gstServerAddr)
			//checkError(err)
			if err != nil {
				if server_conn_error_total > 1000 {
					g_Logger.Fatalln("Error: server_conn_error_total > 1000")
				}
				server_conn_error_total++
				g_Logger.Println("Error: [net.Dial('tcp', service)]:", err.Error())

				continue
			}
			//g_Logger.Println("conn.", service)

			//acc: 		src---client
			//dst: 		client---serever
			//pack: 	acc read recv, dst wirte send.
			//unpack:	dst read recv, acc wirte send.

			gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gsst)
			g_gsst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

			go srcTOdstP_count(rawClient, gstServer, gctx)
			go srcTOdstUn_count(gstServer, rawClient, gctx)
			g_Logger.Printf("go [%d].\n", gctx.GetGsId())
			break
		}
	}
}

func run() {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	var lstnaddr string
	var connaddr []string

	lstnaddr = g_gsconfig.Listen
	connaddr = g_gsconfig.GetServers()

	g_Logger.Println("Listen_Addr:", lstnaddr)
	g_Logger.Println("Conn_Addr:", connaddr)
	g_Logger.Println("Begin......")

	service := lstnaddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	chanConnList := make(chan net.Conn, 1024)
	defer close(chanConnList)

	for i := 0; i < g_connHandleGoNum; i++ {
		go createToGstServerConnHandler(chanConnList, connaddr[0])
	}
	for {
		rawClient, err := listener.Accept()
		if err != nil {
			g_Logger.Println("Error:", err)
			continue
		}
		chanConnList <- rawClient
	}
}

func run_pipe_test_listen(lstnAddr, connAddr string) {

	g_Logger.Println("Listen_Addr:", lstnAddr)
	g_Logger.Println("Conn_Addr:", connAddr)
	g_Logger.Println("Begin......")

	service := lstnAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	chanConnList := make(chan net.Conn, 1024)
	defer close(chanConnList)

	for i := 0; i < g_connHandleGoNum; i++ {
		go createToGstServerConnHandler(chanConnList, connAddr)
	}
	for {
		rawClient, err := listener.Accept()
		if err != nil {
			g_Logger.Println("Error:", err)
			continue
		}
		chanConnList <- rawClient
	}
}

func run_pipe_test_wg(sc gstestpipe.RawdataPiPe, gss gstestpipe.GsPiPe, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	acc := sc.GetServerConn()
	dst := gss.GetConn()

	g_Logger.Println("Test_Mt_model:", g_Mt_model)
	g_log_List.GSIpLogger.Printf("ip: %s\n", acc.RemoteAddr().String())

	gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gsst)
	g_gsst.GetStatusConnList().Add(gctx.GetGsId(), acc, dst)

	wg.Add(2)
	go srcTOdstP_wg(acc, dst, wg, gctx)
	go srcTOdstUn_wg(dst, acc, wg, gctx)
	g_Logger.Println("Gstunnel go.")

}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
}

func srcTOdstP_count(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	atomic.AddInt32(&g_goPackTotal, 1)
	srcTOdstP(src, dst, gctx)
	atomic.AddInt32(&g_goPackTotal, -1)

}

func srcTOdstUn_count(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	atomic.AddInt32(&g_goUnpackTotal, 1)
	srcTOdstUn(src, dst, gctx)
	atomic.AddInt32(&g_goUnpackTotal, -1)
}

func IsTheVersionConsistent_send(dst net.Conn, apack gstunnellib.GsPack, wlent *int64) error {
	return gstunnellib.IsTheVersionConsistent_send(dst, apack, wlent)
}

func ChangeCryKey_send(dst net.Conn, apack gstunnellib.GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return gstunnellib.ChangeCryKey_send(dst, apack, ChangeCryKey_Total, wlent)
}

func srcTOdstP(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	if g_Mt_model {
		srcTOdstP_mt(src, dst, gctx)
	} else {
		srcTOdstP_st(src, dst, gctx)
	}
}

func srcTOdstUn(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	if g_Mt_model {
		srcTOdstUn_mt(src, dst, gctx)
	} else {
		srcTOdstUn_st(src, dst, gctx)
	}
}

func srcTOdstP_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(src, dst, gctx)
	} else {
		srcTOdstP_st(src, dst, gctx)
	}
}

func srcTOdstUn_wg(src net.Conn, dst net.Conn, wg *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(src, dst, gctx)
	} else {
		srcTOdstUn_st(src, dst, gctx)
	}
}
