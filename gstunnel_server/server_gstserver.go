package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsmap"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type gstServer struct {
	listenAddr            string
	rawServiceAddr        string
	chanConnList          chan net.Conn
	listener              net.Listener
	maxTryConnService     int
	connHandleGoNum       int
	table_ip_clientPubKey *gsmap.GsMapStr_RSA
	//chan_exchangeRSAkey_handle chan *exchangeRSAkey_handle
	table_GSTConn gsmap.IGSMapLock //debug
}

type exchangeRSAkey_handle struct {
	pack, unpack *gsobj.GSTObj
	clientip     string
}

func newExchangeRSAkey_handle(pack, unpack *gsobj.GSTObj, clientip string) *exchangeRSAkey_handle {
	return &exchangeRSAkey_handle{pack, unpack, clientip}
}

func newGstServer(listenAddr, rawServiceAddr string) *gstServer {
	if debug_gstServer {
		return &gstServer{
			listenAddr:            listenAddr,
			rawServiceAddr:        rawServiceAddr,
			chanConnList:          make(chan net.Conn, 10240),
			maxTryConnService:     100 * 6,
			connHandleGoNum:       16,
			table_ip_clientPubKey: gsmap.NewGSMapStr_RSA(),
			//chan_exchangeRSAkey_handle: make(chan *exchangeRSAkey_handle, 1000),
			table_GSTConn: gsmap.NewGSMapLock(), //debug
		}
	} else {
		return &gstServer{
			listenAddr:            listenAddr,
			rawServiceAddr:        rawServiceAddr,
			chanConnList:          make(chan net.Conn, 10240),
			maxTryConnService:     100 * 6,
			connHandleGoNum:       16,
			table_ip_clientPubKey: gsmap.NewGSMapStr_RSA(),
			//chan_exchangeRSAkey_handle: make(chan *exchangeRSAkey_handle, 1000),
		}
	}
}

func (s *gstServer) getIP(conn net.Addr) string {
	return gstunnellib.GetIP(conn)
}

func (s *gstServer) getClientPubKeyFromIP(clientip string) (*gsrsa.RSA, bool) {
	return s.table_ip_clientPubKey.Get(clientip)
}

func (s *gstServer) addClientPubKeyFromIP(clientip string, clientkey *gsrsa.RSA) {
	s.table_ip_clientPubKey.Add(clientip, clientkey)
}

/*
func (s *gstServer) exchangeRSAkey(pack, unpack *gsobj.GSTObj, clientip string) {
	//pubkpack := clientPack.ClientPublicKeyPack()

	obj := unpack
	var clientkey *gsrsa.RSA
	for {
		//obj.Src.SetReadDeadline(time.Now().Add(obj.NetworkTimeout))
		rlen, err := obj.ReadNetSrc(obj.Rbuf)
		obj.Rlent += int64(rlen)
		//	recot_un_r.Run()
		if gstunnellib.IsErrorNetUsually(err) {
			checkError_info_GSCtx(err, obj.Gctx)
			return
		} else {
			checkError_panic_GSCtx(err, obj.Gctx)
		}

			//if tmr_out.Run() {
			//	g_logger.Printf("Error: [%d] Time out, func exit.\n", obj.Gctx.GetGsId())
			//	return
			//}

		if rlen == 0 {
			g_logger.Println("Error: obj.Src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			g_logger.Println("Error:", err)
			continue
		}
		//tmr_out.Boot()

		obj.WriteEncryData(obj.Rbuf[:rlen])
		obj.Wbuf, err = obj.GetDecryData()
		checkError_panic_GSCtx(err, obj.Gctx)
		if obj.IsExistsClientKey() {
			clientkey = obj.GetClientRSAKey()
			break
		}
	}
	s.addClientPubKeyFromIP(clientip, clientkey)
	pack.SetClientRSAKey(clientkey)
	unpack.SetClientRSAKey(clientkey)

	unpack.Close()
	pack.Close()
}
*/

func (s *gstServer) newGSTConn_handle(gstClient, rawServer net.Conn, key string, clientrsa *gsrsa.RSA) {
	gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)
	//g_gstst.GetStatusConnList().Add(gctx.GetGsId(), gstClient, rawServer)

	unpack := gsobj.NewGstObjWithServer(gstClient, rawServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		key,
		g_net_read_size,
		g_key_rsa_server,
	)

	pack := gsobj.NewGstObjWithServer(rawServer, gstClient, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
	)

	g_logger.Printf("go [%d].\n", gctx.GetGsId())

	pack.SetClientRSAKey(clientrsa)
	unpack.SetClientRSAKey(clientrsa)

	if debug_gstServer {
		s.table_GSTConn.Add(gctx.GetGsId(), gsobj.NewGSTConn(pack, unpack))
	}

	go s.srcTOdstUn(unpack)
	go s.srcTOdstP(pack)

	//clientaddr := gstClient.RemoteAddr().String()
	//_ = clientaddr
	//clientip := s.getIP(gstClient.RemoteAddr())
	//clientrsa, exists := s.getClientPubKeyFromIP(clientip)

	/*
		if exists {
			gstobjpack.SetClientRSAKey(clientrsa)
			gstobjun.SetClientRSAKey(clientrsa)
			go s.srcTOdstUn(gstobjun)
			go s.srcTOdstP(gstobjpack)
		} else {
			//s.chan_exchangeRSAkey_handle <- newExchangeRSAkey_handle(gstobjpack, gstobjun, clientip)
			//panic("client rsa not exists")
		}
	*/

}

func (s *gstServer) newGSTConn_createRawServer(gstClient net.Conn, key string) {
	clientip := s.getIP(gstClient.RemoteAddr())

	clientkey, exists := s.getClientPubKeyFromIP(clientip)

	if !exists {
		gstClient.Close()
		checkError_info(errors.New("GSTClient key is not exists"))
	}

	connServiceError_count := 0
	var rawServer net.Conn
	var err error
	for {
		rawServer, err = net.Dial("tcp", s.rawServiceAddr)
		checkError_NoExit(err)
		connServiceError_count += 1
		if err == nil {
			break
		}
		if connServiceError_count > s.maxTryConnService {
			checkError(
				fmt.Errorf("connService_count > maxTryConnService(%d)", s.maxTryConnService))
		}
		time.Sleep(time.Millisecond * 100)
		//g_logger.Println("conn.")
	}
	s.newGSTConn_handle(gstClient, rawServer, key, clientkey)
}

func (s *gstServer) newGSTConn_Init(gstClient net.Conn) {

	//one pack
	//unpack
	//potype
	//
	//gstClient.SetDeadline(time.Now()+)
	rp := gstunnellib.NewGSRSAPackNetImpWithGSTServer(g_key, g_key_rsa_server)
	re, err := rp.UnpackOneGSTPackFromNetConn(gstClient, g_networkTimeout)
	if gserror.IsErrorNetUsually(err) {
		checkError_info(err)
		return
	} else {
		checkError_panic(err)
	}
	//clientpubkey exchangeRSAkey------
	if re.IsPOVersion() {
		if re.Version == gsbase.G_Version {
			g_logger.Println("gstunnel POVersion is ok.")
		} else {
			err = errors.New("GSTCLient version is error. " + "error version: " + re.Version)
			panic(err)
		}
		re, err = rp.UnpackOneGSTPackFromNetConn(gstClient, g_networkTimeout)
		if gserror.IsErrorNetUsually(err) {
			checkError_info(err)
			return
		} else {
			checkError_panic(err)
		}
		clientip := s.getIP(gstClient.RemoteAddr())

		if re.IsClientPubKey() {
			s.addClientPubKeyFromIP(clientip, re.ClientKey)
			gstClient.Close()
			g_logger.Printf("gstunnel [%s] ClientPubKey exchange completed.", clientip)
			return
		}
		panic("GSTClient init pack is error")
	}
	//------

	clientip := s.getIP(gstClient.RemoteAddr())
	//key := g_key

	if re.IsClientPubKey() {
		s.addClientPubKeyFromIP(clientip, re.ClientKey)
		gstClient.Close()
		g_logger.Printf("gstunnel [%s] ClientPubKey exchange completed.", clientip)
		return
	}
	if re.IsChangeCryKeyFromGSTClient() {
		key := string(re.Key)
		s.newGSTConn_createRawServer(gstClient, key)
		return
	}
	if re.IsPOGen() {
		panic("error: (*gstServer) newGSTConn_Init(): Init func should not show the POGenOper type package.")
	}

	if re.IsChangeCryKeyFromGSTServer() {
		panic("error: (*gstServer) newGSTConn_Init(): Init func should not show the POChangeCryKeyFromGSTServer type package.")
	}

	if re.IsChangeCryKey() {
		panic("error: (*gstServer) newGSTConn_Init(): Init func should not show the POChangeCryKey type package.")
	}
	if re.IsPOHello() {
		s.newGSTConn_createRawServer(gstClient, g_key)
		return
	}
	panic("gstunnel pack type is error")
}

func (s *gstServer) newGSTConn(gstClient net.Conn) {
	defer gserror.Panic_Recover(g_logger)
	s.newGSTConn_Init(gstClient)
}

/*
func (s *gstServer) go_exchangeRSAkey_handle() {

	//	1	doing	2	done
	//
	map1 := make(map[string]int)

	//chan_eckhandle_sleep := make(chan *exchangeRSAkey_handle, 1000)

	go func() {
		for ehpack := range s.chan_exchangeRSAkey_handle {
			_, ok := map1[ehpack.clientip]
			if !ok {
				go s.exchangeRSAkey(ehpack.pack, ehpack.unpack, ehpack.clientip)
				map1[ehpack.clientip] = 1
			} else {
				/*
					go func(ehpack *exchangeRSAkey_handle) {
						_, exists := s.getClientPubKeyFromIP(ehpack.clientip)
						for !exists {
							time.Sleep(time.Millisecond * 10)
							_, exists = s.getClientPubKeyFromIP(ehpack.clientip)
						}
						clientrsa, _ := s.getClientPubKeyFromIP(ehpack.clientip)
						ehpack.pack.SetClientRSAKey(clientrsa)
						ehpack.unpack.SetClientRSAKey(clientrsa)
						go s.srcTOdstUn(ehpack.unpack)
						go s.srcTOdstP(ehpack.pack)
					}(ehpack)
*/
/*
				go func(ehpack *exchangeRSAkey_handle) {
					time.Sleep(time.Second * 2)
					ehpack.pack.Close()
					ehpack.unpack.Close()
				}(ehpack)
			}
		}
	}()
}
*/
/*
func (s *gstServer) createConnHandler_old(chanconnlist <-chan net.Conn, rawServiceAddr string) {

		maxTryConnService := s.maxTryConnService
		var err error
		for {

			gstClient, ok := <-chanconnlist
			if !ok {
				checkError_info(errors.New("chanconnlist is closed"))
				return
			}
			g_log_List.GSIpLogger.Printf("Gstunnel client ip: %s\n", gstClient.RemoteAddr().String())
			////
			clientip := s.getIP(gstClient.RemoteAddr())
			_, exists := s.getClientPubKeyFromIP(clientip)

			if !exists {
				gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)

				gstobjun := gsobj.NewGstObjWithServer(gstClient, gstClient, gctx,
					g_tmr_display_time,
					g_networkTimeout,
					g_key,
					g_net_read_size,
					g_key_rsa_server,
				)
				s.chan_exchangeRSAkey_handle <- newExchangeRSAkey_handle(gstobjun, gstobjun, clientip)

			} else {

				connServiceError_count := 0
				var rawServer net.Conn
				for {
					rawServer, err = net.Dial("tcp", rawServiceAddr)
					checkError_NoExit(err)
					connServiceError_count += 1
					if err == nil {
						break
					}
					if connServiceError_count > maxTryConnService {
						checkError(
							fmt.Errorf("connService_count > maxTryConnService(%d)", maxTryConnService))
					}
					time.Sleep(time.Millisecond * 100)
					//g_logger.Println("conn.")
				}
				go s.newGSTConn(gstClient, rawServer)
			}
		}
	}
*/
func (s *gstServer) createConnHandler(chanconnlist <-chan net.Conn, rawServiceAddr string) {

	//	maxTryConnService := s.maxTryConnService
	//var err error
	for {

		gstClient, ok := <-chanconnlist
		if !ok {
			checkError_info(errors.New("chanconnlist is closed"))
			return
		}
		g_log_List.GSIpLogger.Printf("Gstunnel client ip: %s\n", gstClient.RemoteAddr().String())
		go s.newGSTConn(gstClient)
	}
}

func (s *gstServer) run() {
	//defer gstunnellib.Panic_Recover_GSCtx(g_logger, gctx)
	//s.go_exchangeRSAkey_handle()

	g_logger.Println("Listen_Addr:", s.listenAddr)
	g_logger.Println("Conn_Addr:", s.rawServiceAddr)
	g_logger.Println("Begin......")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.listenAddr)
	checkError(err)
	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//var rawServer net.Conn

	for i := 0; i < s.connHandleGoNum; i++ {
		go s.createConnHandler(s.chanConnList, s.rawServiceAddr)
	}

	for {
		gstClient, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			checkError_info(err)
			return
		} else if err != nil {
			checkError_NoExit(err)
			continue
		}
		s.chanConnList <- gstClient
	}
}

func (s *gstServer) close() {
	close(s.chanConnList)
	//close(s.chan_exchangeRSAkey_handle)
	s.listener.Close()
}

// service to gstunnel client
func (s *gstServer) srcTOdstP(obj *gsobj.GSTObj) {
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func (s *gstServer) srcTOdstUn(obj *gsobj.GSTObj) {
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}

/*
// service to gstunnel client
func (s *gstServer) srcTOdstP_wg(obj *gsobj.GSTObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func (s *gstServer) srcTOdstUn_wg(obj *gsobj.GSTObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}
*/
