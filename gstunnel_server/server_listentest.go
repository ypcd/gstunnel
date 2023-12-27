package main

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
)

/*
type gstServer_test struct {
	gstServer
}

func newGstServer_test(gstClient, rawServer string) *gstServer_test {
	return &gstServer_test{
		gstServer{
			listenAddr:        gstClient,
			rawServiceAddr:    rawServer,
			chanConnList:      make(chan net.Conn, 10240),
			maxTryConnService: 5000,
		},
	}
}

func (s *gstServer_test) run_pipe_test_listen(listenAddr, rawServiceAddr string) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	//var listenAddr, rawServiceAddr string

	//listenAddr = g_gsconfig.Listen
	//rawServiceAddr = g_gsconfig.GetServer_rand()

	g_Logger.Println("Listen_Addr:", listenAddr)
	g_Logger.Println("Conn_Addr:", rawServiceAddr)
	g_Logger.Println("Begin......")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", listenAddr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//var rawService net.Conn

	for i := 0; i < g_connHandleGoNum; i++ {
		go s.createToRawServiceConnHandler(s.chanConnList, rawServiceAddr)
	}

	for {
		gstClient, err := listener.Accept()
		if err != nil {
			checkError_NoExit(err)
			continue
		}
		s.chanConnList <- gstClient
	}
}

func (s *gstServer_test) run() {
	s.run_pipe_test_listen(s.listenAddr, s.rawServiceAddr)
}

*/

/*
func createToRawServiceConnHandler(chanconnlist <-chan net.Conn, rawServiceAddr string) {

		const maxTryConnService int = 5000
		var err error
		for {

			gstClient, ok := <-chanconnlist
			if !ok {
				checkError_NoExit(errors.New("'gstServer, ok := <-chanconnlist' is error"))
				return
			}
			g_log_List.GSIpLogger.Printf("Gstunnel client ip: %s\n", gstClient.RemoteAddr().String())

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

			go srcTOdstUn(gsobj.GstObj)
			go srcTOdstP(gsobj.GstObj)
			g_Logger.Printf("go [%d].\n", gctx.GetGsId())
		}
	}

// 性能是old版的4倍

	func run_pipe_test_listen(listenAddr, rawServiceAddr string) {
		//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

		//var listenAddr, rawServiceAddr string

		//listenAddr = g_gsconfig.Listen
		//rawServiceAddr = g_gsconfig.GetServer_rand()

		g_Logger.Println("Listen_Addr:", listenAddr)
		g_Logger.Println("Conn_Addr:", rawServiceAddr)
		g_Logger.Println("Begin......")

		tcpAddr, err := net.ResolveTCPAddr("tcp4", listenAddr)
		checkError(err)
		listener, err := net.ListenTCP("tcp", tcpAddr)
		checkError(err)

		//var rawService net.Conn

		chanConnList := make(chan net.Conn, 10240)

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
*/
func run_pipe_test_listen_old(listenAddr, rawServiceAddr string) {
	//defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	g_Logger.Println("Listen_Addr:", listenAddr)
	g_Logger.Println("Conn_Addr:", rawServiceAddr)
	g_Logger.Println("Begin......")

	service := listenAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		gstClient, err := listener.Accept()
		if err != nil {
			checkError_NoExit(err)
			continue
		}
		g_log_List.GSIpLogger.Printf("ip: %s\n", gstClient.RemoteAddr().String())

		service := rawServiceAddr

		connServiceError_count := 0
		const maxTryConnService int = 5000
		var rawServer net.Conn
		for {
			rawServer, err = net.Dial("tcp", service)
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
		g_gstst.GetStatusConnList().Add(gctx.GetGsId(), gstClient, rawServer)

		gstobjun := gsobj.NewGstObj(gstClient, rawServer, gctx,
			g_tmr_display_time,
			g_networkTimeout,
			g_key,
			g_net_read_size,
		)

		gstobjp := gsobj.NewGstObj(rawServer, gstClient, gctx,
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

func run_pipe_test_wg(rawServer net.Conn, gstClient net.Conn, wg *sync.WaitGroup) {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	g_Logger.Println("Test_Mt_model:", g_Mt_model)
	g_log_List.GSIpLogger.Printf("ip: %s\n", gstClient.RemoteAddr().String())

	gctx := gstunnellib.NewGsContextImp(g_gid.GenerateId(), g_gstst)
	g_gstst.GetStatusConnList().Add(gctx.GetGsId(), gstClient, rawServer)

	gstobjun := gsobj.NewGstObj(gstClient, rawServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
	)

	gstobjp := gsobj.NewGstObj(rawServer, gstClient, gctx,
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

// service to gstunnel client
func srcTOdstP(obj *gsobj.GstObj) {
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func srcTOdstUn(obj *gsobj.GstObj) {
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}

// service to gstunnel client
func srcTOdstP_wg(obj *gsobj.GstObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func srcTOdstUn_wg(obj *gsobj.GstObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}
