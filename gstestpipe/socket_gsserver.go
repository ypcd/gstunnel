package gstestpipe

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

var g_timeout time.Duration = time.Second * 6

type gstEchoHandlerObj struct {
	readSize, writeSize int64
}

type GstServerSocketEcho struct {
	RawServerSocket
	echo_running               bool
	key                        string
	list_gstEchoHandler        []*gstEchoHandlerObj
	serverRSAKey, clientRSAKey *gsrsa.RSA
	lock_clientRSAKey          sync.Mutex
}

/*
func NewGstServerSocketEcho_RandAddr() *GstServerSocketEcho {
	return &GstServerSocketEcho{RawServerSocket: RawServerSocket{connList: make(chan net.Conn, 10240),
		serverAddr: gstunnellib.GetNetLocalRDPort()},
		key: g_key_Default}
}
*/

func NewGstRSAServerSocketEcho_RandAddr(serverRSAKey *gsrsa.RSA) *GstServerSocketEcho {
	return &GstServerSocketEcho{RawServerSocket: RawServerSocket{connList: make(chan net.Conn, 10240),
		serverAddr: gstunnellib.GetNetLocalRDPort()},
		key:          g_key_Default,
		serverRSAKey: serverRSAKey,
	}
}

func NewGstRSAServerSocketEcho_RandAddr_DefaultKey() *GstServerSocketEcho {
	return &GstServerSocketEcho{RawServerSocket: RawServerSocket{connList: make(chan net.Conn, 10240),
		serverAddr: gstunnellib.GetNetLocalRDPort()},
		key:          g_key_Default,
		serverRSAKey: gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
	}
}

// newGSTPipoHandleEx_Init return
// 0--error
// 1--ok
// 2--ClientPubKey
func (ss *GstServerSocketEcho) newGSTPipoHandleEx_Init(obj *gstEchoHandlerObj, key string, gstClient net.Conn, apacknetRead, apacknetWrite gstunnellib.IGSRSAPackNet) int {

	rp := apacknetRead
	//gstClient := conn
	//gstClient.SetDeadline(time.Now().Add(g_timeout))
	re, err := rp.UnpackOneGSTPackFromNetConn(gstClient, g_timeout)
	if gserror.IsErrorNetUsually(err) {
		checkError_info(err)
		return 0
	} else {
		checkError_panic(err)
	}

	if re.IsPOVersion() {
		if re.Version == gsbase.G_Version {
			g_logger.Println("gstunnel POVersion is ok.")
		} else {
			err = errors.New("GSTCLient version is error. " + "error version: " + re.Version)
			panic(err)
		}
		re, err := rp.UnpackOneGSTPackFromNetConn(gstClient, g_timeout)
		if gserror.IsErrorNetUsually(err) {
			checkError_info(err)
			return 0
		} else {
			checkError_panic(err)
		}
		clientip := gstunnellib.GetIP(gstClient.RemoteAddr())

		if re.IsClientPubKey() {
			//s.addClientPubKeyFromIP(clientip, re.ClientKey)
			ss.lock_clientRSAKey.Lock()
			ss.clientRSAKey = re.ClientKey
			ss.lock_clientRSAKey.Unlock()
			gstClient.Close()
			g_logger.Printf("gstunnel [%s] ClientPubKey exchange completed.", clientip)
			return 2
		}
	}
	//if !
	clientip := gstunnellib.GetIP(gstClient.RemoteAddr())
	var newkey string

	if re.IsClientPubKey() {
		//s.addClientPubKeyFromIP(clientip, re.ClientKey)
		ss.lock_clientRSAKey.Lock()
		ss.clientRSAKey = re.ClientKey
		ss.lock_clientRSAKey.Unlock()
		gstClient.Close()
		g_logger.Printf("gstunnel [%s] ClientPubKey exchange completed.", clientip)
		return 2
	}
	if re.IsChangeCryKeyFromGSTClient() {
		newkey = string(re.Key)
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
	ss.lock_clientRSAKey.Lock()
	if ss.clientRSAKey == nil {
		panic("ss.clientRSAKey==nil")
	}
	ss.lock_clientRSAKey.Unlock()

	readpack := gstunnellib.NewGSRSAPackNetImpWithGSTServer(newkey, ss.serverRSAKey)
	writepack := gstunnellib.NewGSRSAPackNetImpWithGSTServer(key, ss.serverRSAKey)
	ss.lock_clientRSAKey.Lock()
	readpack.SetClientRSAKey(ss.clientRSAKey)
	writepack.SetClientRSAKey(ss.clientRSAKey)
	ss.lock_clientRSAKey.Unlock()

	gstPipoHandleEx(obj, key, gstClient, readpack, writepack)
	return 1
}

func (ss *GstServerSocketEcho) newGSTPipoHandleEx(obj *gstEchoHandlerObj, key string, conn net.Conn, apacknetRead, apacknetWrite gstunnellib.IGSRSAPackNet) {
	defer gserror.Panic_Recover(g_logger)
	ss.newGSTPipoHandleEx_Init(obj, key, conn, apacknetRead, apacknetWrite)
	/*
		switch re {
		case 0: //eroor
			return
		case 2: //ClientPubKey
			return
		case 1: //ok
			gstPipoHandleEx(obj, ss.key, conn, gstunnellib.NewGSRSAPackNetImpWithGSTServer(ss.key, ss.serverRSAKey), gstunnellib.NewGSRSAPackNetImpWithGSTServer(ss.key, ss.serverRSAKey))
		}
	*/
}

func (ss *GstServerSocketEcho) gstEchoHandler() {
	if ss.echo_running {
		return
	}

	go func() {
		for !ss.close.Load() {
			for conn1 := range ss.connList {
				obj := gstEchoHandlerObj{}

				go ss.newGSTPipoHandleEx(&obj, ss.key, conn1, gstunnellib.NewGSRSAPackNetImpWithGSTServer(ss.key, ss.serverRSAKey), gstunnellib.NewGSRSAPackNetImpWithGSTServer(ss.key, ss.serverRSAKey))
				ss.list_gstEchoHandler = append(ss.list_gstEchoHandler, &obj)
			}
		}
	}()
	ss.echo_running = true
}

func (ss *GstServerSocketEcho) Run() {
	ss.listen()
	ss.gstEchoHandler()
}

func (ss *GstServerSocketEcho) GetKey() string { return ss.key }
