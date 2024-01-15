package gstestpipe

import (
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsnetconn"
)

type RawClientSocket struct {
	close           atomic.Bool
	serverAddr      string
	connListIn      []*gsnetconn.GSTNetConn
	lock_connListIn sync.Mutex
}

func NewRawClientSocket(serverAddr string) *RawClientSocket {
	return &RawClientSocket{serverAddr: serverAddr}
}

func (ss *RawClientSocket) createConn() net.Conn {
	client, err := net.Dial("tcp4", ss.serverAddr)
	checkError_panic(err)

	//gstunnellib.NetConnWriteAll(client, gspackoper.)
	gconnclient := gsnetconn.NewGSTNetConn(client)

	ss.lock_connListIn.Lock()
	ss.connListIn = append(ss.connListIn, gconnclient)
	ss.lock_connListIn.Unlock()
	return gconnclient
}

func (ss *RawClientSocket) Run() net.Conn {
	return ss.createConn()
}

func (ss *RawClientSocket) Close() {
	ss.close.Swap(true)

	ss.lock_connListIn.Lock()
	defer ss.lock_connListIn.Unlock()
	for _, v := range ss.connListIn {
		v.Close()
	}
}

type RawClientSocketEcho struct {
	RawClientSocket
}

func NewRawClientSocketEcho(serverAddr string) *RawClientSocketEcho {
	return &RawClientSocketEcho{RawClientSocket: RawClientSocket{serverAddr: serverAddr}}
}

func (ss *RawClientSocketEcho) createEchoHandler(client net.Conn) {
	go func() {

		_, err := io.Copy(client, client)
		if gstunnellib.IsErrorNetUsually(err) {
			checkError_info(err)
			return
		}
		if err != nil {
			checkError(err)
			return
		}

	}()
}

func (ss *RawClientSocketEcho) Run() {
	ss.createEchoHandler(ss.createConn())
}
