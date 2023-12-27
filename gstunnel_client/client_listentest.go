package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

func createConnHandler(chanconnlist <-chan net.Conn, rawServiceAddr string) {

	maxTryConnService := 5000
	var err error
	for {

		rawClient, ok := <-chanconnlist
		if !ok {
			checkError_NoExit(errors.New("'gstClient, ok := <-chanconnlist' is error"))
			return
		}
		g_log_List.GSIpLogger.Printf("Raw client ip: %s\n", rawClient.RemoteAddr().String())

		connServiceError_count := 0
		var gstServer net.Conn
		for {
			gstServer, err = net.Dial("tcp", rawServiceAddr)
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
		g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

		gstobjun := gsobj.NewGstObj(gstServer, rawClient, gctx,
			g_tmr_display_time,
			g_networkTimeout,
			g_key,
			g_net_read_size,
		)

		gstobjp := gsobj.NewGstObj(rawClient, gstServer, gctx,
			g_tmr_display_time,
			g_networkTimeout,
			g_key,
			g_net_read_size,
		)

		go srcTOdstUn(gstobjun)
		go srcTOdstP(gstobjp)
		g_Logger.Printf("go [%d].\n", gctx.GetGsId())
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

	chanConnList := make(chan net.Conn, 10240)
	defer close(chanConnList)

	for i := 0; i < g_connHandleGoNum; i++ {
		go createConnHandler(chanConnList, connAddr)
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

func run_pipe_test_wg_old(sc gstestpipe.RawdataPiPe, gss gstestpipe.GsPiPe, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	rawClient := sc.GetServerConn()
	gstServer := gss.GetConn()

	g_Logger.Println("Test_Mt_model:", g_Mt_model)
	g_log_List.GSIpLogger.Printf("ip: %s\n", rawClient.RemoteAddr().String())

	gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gstst)
	g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

	gstobjun := gsobj.NewGstObj(gstServer, rawClient, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
	)

	gstobjp := gsobj.NewGstObj(rawClient, gstServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
	)

	wg.Add(2)
	go srcTOdstUn_wg(gstobjun, wg)
	go srcTOdstP_wg(gstobjp, wg)

	g_Logger.Println("Gstunnel go.")

}

func run_pipe_test_wg(rawClient, gstServer net.Conn, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	//rawClient := sc.GetServerConn()
	//gstServer := gss.GetConn()

	g_Logger.Println("Test_Mt_model:", g_Mt_model)
	g_log_List.GSIpLogger.Printf("ip: %s\n", rawClient.RemoteAddr().String())

	gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gstst)
	g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

	gstobjun := gsobj.NewGstObj(gstServer, rawClient, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
	)

	gstobjp := gsobj.NewGstObj(rawClient, gstServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
	)

	wg.Add(2)
	go srcTOdstUn_wg(gstobjun, wg)
	go srcTOdstP_wg(gstobjp, wg)

	g_Logger.Println("Gstunnel go.")

}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
}

func srcTOdstP_count(obj *gsobj.GstObj) {
	atomic.AddInt32(&g_goPackTotal, 1)
	srcTOdstP(obj)
	atomic.AddInt32(&g_goPackTotal, -1)

}

func srcTOdstUn_count(obj *gsobj.GstObj) {
	atomic.AddInt32(&g_goUnpackTotal, 1)
	srcTOdstUn(obj)
	atomic.AddInt32(&g_goUnpackTotal, -1)
}

func srcTOdstP(obj *gsobj.GstObj) {
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

func srcTOdstUn(obj *gsobj.GstObj) {
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}

func srcTOdstP_wg(obj *gsobj.GstObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

func srcTOdstUn_wg(obj *gsobj.GstObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}
