package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsmap"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

const g_time_gst_init_wait time.Duration = time.Second * 60

type gstClient struct {
	listenAddr        string
	gstServerAddr     string
	chanConnList      chan net.Conn
	listener          net.Listener
	maxTryConnService int
	connHandleGoNum   int
	clientRSAKey      *gsrsa.RSA
	table_GSTConn     gsmap.IGSMapLock //debug
}

func newGstClient(listenAddr, gstServerAddr string) *gstClient {
	if debug_gstClient {
		return &gstClient{
			listenAddr:        listenAddr,
			gstServerAddr:     gstServerAddr,
			chanConnList:      make(chan net.Conn, 10240),
			maxTryConnService: 5000,
			connHandleGoNum:   16,
			clientRSAKey:      gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen),
			table_GSTConn:     gsmap.NewGSMapLock(), //debug
		}
	} else {
		return &gstClient{
			listenAddr:        listenAddr,
			gstServerAddr:     gstServerAddr,
			chanConnList:      make(chan net.Conn, 10240),
			maxTryConnService: 5000,
			connHandleGoNum:   16,
			clientRSAKey:      gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen),
		}
	}
}

/*
func (c *gstClient) exchangeRSAkey_old(pack, unpack *gsobj.GSTObj) {

	pubpack := pack.ClientPublicKeyPack()
	_, err := pack.NetConnWriteAll(pubpack)
	checkError_panic(err)
	go c.srcTOdstUn(unpack)
	go c.srcTOdstP(pack)
}*/

func (c *gstClient) init_exchangeRSAkey_conn_close(conn net.Conn) {
	loopNum := g_time_gst_init_wait / time.Second
	for i := 0; i < int(loopNum); i++ {
		if gstunnellib.IsNetConnClosed_Read(conn) {
			return
		}
		time.Sleep(1 * time.Second)
		g_logger.Println("Info: Waiting for server to close GSTunnel initialization connection.")
	}
}

func (c *gstClient) exchangeRSAkey(pack *gsobj.GSTObj) {
	err := pack.VersionPack_send()
	checkError_panic(err)

	pubpack := pack.ClientPublicKeyPack()
	_, err = pack.NetConnWriteAll(pubpack)
	checkError_panic(err)
}

func (c *gstClient) init() {
	maxTryConnService := c.maxTryConnService
	var err error
	var connServiceError_count int = 0

	//connServiceError_count = 0
	var gstServer net.Conn
	for {
		gstServer, err = net.Dial("tcp", c.gstServerAddr)
		checkError_NoExit(err)
		connServiceError_count += 1
		if err == nil {
			break
		}
		if connServiceError_count > maxTryConnService {
			checkError(
				fmt.Errorf("connService_count > maxTryConnService(%d)", maxTryConnService))
		}
	}
	gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)

	pack := gsobj.NewGstObjWithClient(gstServer, gstServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
		c.clientRSAKey,
	)
	defer pack.Close()

	c.exchangeRSAkey(pack)

	c.init_exchangeRSAkey_conn_close(gstServer)
}

func (c *gstClient) newGSTConn(rawClient, gstServer net.Conn) {
	gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)
	//g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

	//clientRSAKey := gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen)

	unpack := gsobj.NewGstObjWithClient(gstServer, rawClient, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
		c.clientRSAKey,
	)

	pack := gsobj.NewGstObjWithClient(rawClient, gstServer, gctx,
		g_tmr_display_time,
		g_networkTimeout,
		g_key,
		g_net_read_size,
		g_key_rsa_server,
		c.clientRSAKey,
	)

	g_logger.Printf("go [%d].\n", gctx.GetGsId())

	if debug_gstClient {
		c.table_GSTConn.Add(gctx.GetGsId(), gsobj.NewGSTConn(pack, unpack))
	}
	//c.exchangeRSAkey(gstobjp, gstobjun)
	go c.srcTOdstUn(unpack)
	go c.srcTOdstP(pack)
}

func (c *gstClient) createConnHandler(chanconnlist <-chan net.Conn, gstServerAddr string) {

	maxTryConnService := c.maxTryConnService
	var err error
	for {

		rawClient, ok := <-chanconnlist
		if !ok {
			checkError_info(errors.New("chanconnlist is not ok"))
			return
		}
		g_log_List.GSIpLogger.Printf("Raw client ip: %s\n", rawClient.RemoteAddr().String())

		connServiceError_count := 0
		var gstServer net.Conn
		for {
			gstServer, err = net.Dial("tcp", gstServerAddr)
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

		go c.newGSTConn(rawClient, gstServer)
	}
}

func (c *gstClient) run() {
	//defer gstunnellib.Panic_Recover_GSCtx(g_logger, gctx)

	g_logger.Println("GSTClient Listen_Addr:", c.listenAddr)
	g_logger.Println("GSTServer_Addr:", c.gstServerAddr)
	g_logger.Println("Begin......")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", c.listenAddr)
	checkError(err)
	c.listener, err = net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	c.init()

	for i := 0; i < c.connHandleGoNum; i++ {
		go c.createConnHandler(c.chanConnList, c.gstServerAddr)
	}

	for {
		gstClient, err := c.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			checkError_info(err)
			return
		} else if err != nil {
			checkError_NoExit(err)
			continue
		}
		c.chanConnList <- gstClient
	}
}

func (c *gstClient) close() {
	close(c.chanConnList)
	c.listener.Close()
}

// service to gstunnel client
func (c *gstClient) srcTOdstP(obj *gsobj.GSTObj) {
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func (c *gstClient) srcTOdstUn(obj *gsobj.GSTObj) {
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}

/*
// service to gstunnel client
func (c *gstClient) srcTOdstP_wg(obj *gsobj.GSTObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstP_mt(obj)
	} else {
		srcTOdstP_st(obj)
	}
}

// gstunnel client to service
func (c *gstClient) srcTOdstUn_wg(obj *gsobj.GSTObj, wg *sync.WaitGroup) {
	defer wg.Done()
	if g_Mt_model {
		srcTOdstUn_mt(obj)
	} else {
		srcTOdstUn_st(obj)
	}
}
*/
