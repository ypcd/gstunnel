package main

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe"
)

/*
const gf_connHandleGoNum = 16

func createConnHandler(chanconnlist <-chan net.Conn, rawServiceAddr string) {

		maxTryConnService := 5000
		var err error
		for {

			rawClient, ok := <-chanconnlist
			if !ok {
			checkError_info(errors.New("chanconnlist is closed"))
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
				//g_logger.Println("conn.")
			}

			gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)
			g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

			gstobjun := gsobj.NewGstObjWithClient(gstServer, rawClient, gctx,
				g_tmr_display_time,
				g_networkTimeout,
				g_key,
				g_net_read_size,
			)

			gstobjp := gsobj.NewGstObjWithClient(rawClient, gstServer, gctx,
				g_tmr_display_time,
				g_networkTimeout,
				g_key,
				g_net_read_size,
			)

			go srcTOdstUn(gstobjun)
			go srcTOdstP(gstobjp)
			g_logger.Printf("go [%d].\n", gctx.GetGsId())
		}
	}

func run_pipe_test_listen(lstnAddr, connAddr string) {

		g_logger.Println("Listen_Addr:", lstnAddr)
		g_logger.Println("Conn_Addr:", connAddr)
		g_logger.Println("Begin......")

		service := lstnAddr
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		checkError(err)
		listener, err := net.ListenTCP("tcp", tcpAddr)
		checkError(err)

		chanConnList := make(chan net.Conn, 10240)
		defer close(chanConnList)

		for i := 0; i < gf_connHandleGoNum; i++ {
			go createConnHandler(chanConnList, connAddr)
		}
		for {
			rawClient, err := listener.Accept()
			if err != nil {
				g_logger.Println("Error:", err)
				continue
			}
			chanConnList <- rawClient
		}
	}
*/
func run_pipe_test_wg_old(sc gstestpipe.IRawdataPiPe, gss gstestpipe.IGsPiPe, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_logger, gctx)

	rawClient := sc.GetServerConn()
	gstServer := gss.GetConn()

	g_logger.Println("Test_Mt_model:", g_Mt_model)
	g_log_List.GSIpLogger.Printf("ip: %s\n", rawClient.RemoteAddr().String())

	gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)
	g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

	gstobjun := gsobj.NewGstObjWithClient(gstServer, rawClient, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
		gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
	)

	gstobjp := gsobj.NewGstObjWithClient(rawClient, gstServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
		gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
	)

	ExchangeRSAkey(gstobjp, gstobjun, wg)

	g_logger.Println("Gstunnel go.")

}

func run_pipe_test_wg(rawClient, gstServer net.Conn, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_logger, gctx)

	wg.Add(1)
	defer wg.Done()
	//rawClient := sc.GetServerConn()
	//gstServer := gss.GetConn()

	g_logger.Println("Test_Mt_model:", g_Mt_model)
	g_log_List.GSIpLogger.Printf("ip: %s\n", rawClient.RemoteAddr().String())

	gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)
	g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

	gstobjun := gsobj.NewGstObjWithClient(gstServer, rawClient, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
		gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
	)

	gstobjp := gsobj.NewGstObjWithClient(rawClient, gstServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
		gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
	)

	ExchangeRSAkey(gstobjp, gstobjun, wg)

	g_logger.Println("Gstunnel go.")

}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
}

func srcTOdstP_count(obj *gsobj.GSTObj) {
	atomic.AddInt32(&g_goPackTotal, 1)
	srcTOdstP(obj)
	atomic.AddInt32(&g_goPackTotal, -1)

}

func srcTOdstUn_count(obj *gsobj.GSTObj) {
	atomic.AddInt32(&g_goUnpackTotal, 1)
	srcTOdstUn(obj)
	atomic.AddInt32(&g_goUnpackTotal, -1)
}

func srcTOdstP(obj *gsobj.GSTObj) {
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

func srcTOdstUn(obj *gsobj.GSTObj) {
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}

func srcTOdstP_wg(obj *gsobj.GSTObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

func srcTOdstUn_wg(obj *gsobj.GSTObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}

func ExchangeRSAkey(pack, unpack *gsobj.GSTObj, wg *sync.WaitGroup) {
	pubpack := pack.ClientPublicKeyPack()
	_, err := pack.NetConnWriteAll(pubpack)
	checkError_panic(err)

	wg.Add(2)
	go srcTOdstUn_wg(unpack, wg)
	go srcTOdstP_wg(pack, wg)
}
