package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
)

type gstClient struct {
	listenAddr        string
	rawServiceAddr    string
	chanConnList      chan net.Conn
	maxTryConnService int
}

func newGstClient(listenAddr, rawServiceAddr string) *gstClient {
	return &gstClient{
		listenAddr:        listenAddr,
		rawServiceAddr:    rawServiceAddr,
		chanConnList:      make(chan net.Conn, 10240),
		maxTryConnService: 5000,
	}
}

func (s *gstClient) createConnHandler(chanconnlist <-chan net.Conn, rawServiceAddr string) {

	maxTryConnService := s.maxTryConnService
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

		go s.srcTOdstUn(gstobjun)
		go s.srcTOdstP(gstobjp)
		g_Logger.Printf("go [%d].\n", gctx.GetGsId())
	}
}

func (s *gstClient) run() {
	//defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	g_Logger.Println("Listen_Addr:", s.listenAddr)
	g_Logger.Println("Conn_Addr:", s.rawServiceAddr)
	g_Logger.Println("Begin......")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.listenAddr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//var rawServer net.Conn

	for i := 0; i < g_connHandleGoNum; i++ {
		go s.createConnHandler(s.chanConnList, s.rawServiceAddr)
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

func (s *gstClient) close() {
	close(s.chanConnList)
}

// service to gstunnel client
func (s *gstClient) srcTOdstP(obj *gsobj.GstObj) {
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func (s *gstClient) srcTOdstUn(obj *gsobj.GstObj) {
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}

/*
// service to gstunnel client
func (s *gstClient) srcTOdstP_wg(obj *gsobj.GstObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func (s *gstClient) srcTOdstUn_wg(obj *gsobj.GstObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}
*/
