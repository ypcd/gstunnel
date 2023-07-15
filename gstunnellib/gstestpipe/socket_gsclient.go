package gstestpipe

import (
	"net"
)

type GstClientSocketEchoImp struct {
	RawClientSocket
	key                 string
	list_gstEchoHandler []*gstEchoHandlerObj
}

func NewGstClientSocketEcho_defaultKey(serverAddr string) *GstClientSocketEchoImp {
	return &GstClientSocketEchoImp{RawClientSocket: RawClientSocket{serverAddr: serverAddr},
		key: g_key_Default}
}

func (ss *GstClientSocketEchoImp) createGstEchoHandler(client net.Conn) {
	obj := gstEchoHandlerObj{}
	go gstPipoHandleEx(&obj, ss.key, client)
	ss.list_gstEchoHandler = append(ss.list_gstEchoHandler, &obj)
}

func (ss *GstClientSocketEchoImp) Run() {
	ss.createGstEchoHandler(ss.createConn())
}

func (ss *GstClientSocketEchoImp) GetKey() string { return ss.key }
