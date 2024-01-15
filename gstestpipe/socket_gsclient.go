package gstestpipe

import (
	"fmt"
	"net"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsnetconn"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type GstClientSocketEchoImp struct {
	RawClientSocket
	key                        string
	list_gstEchoHandler        []*gstEchoHandlerObj
	serverRSAKey, clientRSAKey *gsrsa.RSA
}

/*
	func NewGstClientSocketEcho_defaultKey(serverAddr string) *GstClientSocketEchoImp {
		return &GstClientSocketEchoImp{RawClientSocket: RawClientSocket{serverAddr: serverAddr},
			key: g_key_Default}
	}
*/
func NewGstRSAClientSocketEcho_defaultKey(serverAddr string, clientRSAKey *gsrsa.RSA) *GstClientSocketEchoImp {
	cs := &GstClientSocketEchoImp{RawClientSocket: RawClientSocket{serverAddr: serverAddr},
		key:          g_key_Default,
		serverRSAKey: gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
		clientRSAKey: clientRSAKey,
	}
	cs.init()
	return cs
}

func (ss *GstClientSocketEchoImp) createConn() net.Conn {
	client, err := net.Dial("tcp4", ss.serverAddr)
	checkError_panic(err)

	rsapack := gstunnellib.NewGSRSAPackNetImpWithGSTClient(ss.key, ss.serverRSAKey, ss.clientRSAKey)
	_, err = gstunnellib.NetConnWriteAll(client, rsapack.PackPOHello())
	checkError_panic(err)

	gconn := gsnetconn.NewGSTNetConn(client)

	ss.lock_connListIn.Lock()
	ss.connListIn = append(ss.connListIn, gconn)
	ss.lock_connListIn.Unlock()
	return gconn
}

func exchangeRSAkey_obj(pack *gsobj.GSTObj) {
	err := pack.VersionPack_send()
	checkError_panic(err)

	pubpack := pack.ClientPublicKeyPack()
	_, err = pack.NetConnWriteAll(pubpack)
	checkError_panic(err)
}

func gstclient_ChangeCryKeyFromGSTClient_send(dst net.Conn, pack gstunnellib.IGSRSAPackNet) (outkey string) {
	pubpack, key := pack.ChangeCryKeyFromGSTClient()

	_, err := gstunnellib.NetConnWriteAll(dst, pubpack)
	checkError_panic(err)
	return key
}

func gstclient_exchangeRSAkey_send(dst net.Conn, pack gstunnellib.IGSRSAPackNet) {
	pubpack := pack.PackVersion()
	_, err := gstunnellib.NetConnWriteAll(dst, pubpack)
	checkError_panic(err)

	pubpack = pack.ClientPublicKeyPack()
	_, err = gstunnellib.NetConnWriteAll(dst, pubpack)
	checkError_panic(err)
}

func gstclient_init_exchangeRSAkey_conn_close(conn net.Conn) {
	loopNum := 3 // g_time_gst_init_wait / time.Second
	for i := 0; i < int(loopNum); i++ {
		if gstunnellib.IsNetConnClosed_Read(conn) {
			return
		}
		time.Sleep(1 * time.Second)
		g_logger.Println("Info: Waiting for server to close GSTunnel initialization connection.")
	}

}

func gstclient_init_exchangeRSAkey(gstServerAddr string, key string, serverRSAKey, clientRSAKey *gsrsa.RSA) {
	maxTryConnService := 200
	var err error
	var connServiceError_count int = 0

	//connServiceError_count = 0
	var netconn1 net.Conn
	for {
		netconn1, err = net.Dial("tcp", gstServerAddr)
		checkError_info(err)
		connServiceError_count += 1
		if err == nil {
			break
		}
		if connServiceError_count > maxTryConnService {
			checkError(
				fmt.Errorf("connService_count > maxTryConnService(%d)", maxTryConnService))
		}
	}
	gstServer := gsnetconn.NewGSTNetConn(netconn1)

	//gctx := gstunnellib.NewGSContextImp(1, nil)

	/*
		pack := gsobj.NewGstObjWithClient(gstServer, gstServer, gctx,
			20,
			60,
			key,
			1024*6,
			serverRSAKey,
			clientRSAKey,
		)
		defer pack.Close()
	*/
	pack := gstunnellib.NewGSRSAPackNetImpWithGSTClient(key, serverRSAKey, clientRSAKey)
	gstclient_exchangeRSAkey_send(gstServer, pack)

	gstclient_init_exchangeRSAkey_conn_close(gstServer)
}

func (c *GstClientSocketEchoImp) init() {
	gstclient_init_exchangeRSAkey(c.serverAddr, c.key, c.serverRSAKey, c.clientRSAKey)
}

/*
func (c *GstClientSocketEchoImp) exchangeRSAkey(pack *gsobj.GSTObj) {
	err := pack.VersionPack_send()
	checkError_panic(err)

	pubpack := pack.ClientPublicKeyPack()
	_, err = pack.NetConnWriteAll(pubpack)
	checkError_panic(err)
}

func (c *GstClientSocketEchoImp) init_exchangeRSAkey_conn_close(conn net.Conn) {
	loopNum := 3 // g_time_gst_init_wait / time.Second
	for i := 0; i < int(loopNum); i++ {
		if gstunnellib.IsNetConnClosed_Read(conn) {
			return
		}
		time.Sleep(1 * time.Second)
		g_logger.Println("Info: Waiting for server to close GSTunnel initialization connection.")
	}

}

func (c *GstClientSocketEchoImp) init() {
	maxTryConnService := 200
	var err error
	var connServiceError_count int = 0

	//connServiceError_count = 0
	var netconn1 net.Conn
	for {
		netconn1, err = net.Dial("tcp", c.serverAddr)
		checkError_info(err)
		connServiceError_count += 1
		if err == nil {
			break
		}
		if connServiceError_count > maxTryConnService {
			checkError(
				fmt.Errorf("connService_count > maxTryConnService(%d)", maxTryConnService))
		}
	}
	gstServer := gsnetconn.NewGSTNetConn(netconn1)

	gctx := gstunnellib.NewGSContextImp(1, nil)

	pack := gsobj.NewGstObjWithClient(gstServer, gstServer, gctx,
		20,
		60,
		c.key,
		1024*6,
		c.serverRSAKey,
		c.clientRSAKey,
	)
	defer pack.Close()

	c.exchangeRSAkey(pack)

	c.init_exchangeRSAkey_conn_close(gstServer)
}
*/
/*
	func (ss *GstClientSocketEchoImp) ExchangeRSAkey_old(obj *gstEchoHandlerObj, client net.Conn) {
		unpack := gstunnellib.NewGSRSAPackNetImpWithGSTClient(ss.key, ss.serverRSAKey, ss.clientRSAKey)
		pack := gstunnellib.NewGSRSAPackNetImpWithGSTClient(ss.key, ss.serverRSAKey, ss.clientRSAKey)

		pubpack := pack.ClientPublicKeyPack()
		_, err := gstunnellib.NetConnWriteAll(client, pubpack)
		checkError_panic(err)

		go gstPipoHandleEx(obj, ss.key, client, unpack, pack)

}
*/
func (ss *GstClientSocketEchoImp) createGstEchoHandler(client net.Conn) {
	obj := gstEchoHandlerObj{}

	//go ss.ExchangeRSAkey(&obj, client)

	go gstPipoHandleEx(
		&obj,
		ss.key,
		client,
		gstunnellib.NewGSRSAPackNetImpWithGSTClient(
			ss.key, ss.serverRSAKey, ss.clientRSAKey),
		gstunnellib.NewGSRSAPackNetImpWithGSTClient(
			ss.key, ss.serverRSAKey, ss.clientRSAKey),
	)
	ss.list_gstEchoHandler = append(ss.list_gstEchoHandler, &obj)
}

func (ss *GstClientSocketEchoImp) Run() {
	ss.createGstEchoHandler(ss.createConn())
}

func (ss *GstClientSocketEchoImp) GetKey() string { return ss.key }
