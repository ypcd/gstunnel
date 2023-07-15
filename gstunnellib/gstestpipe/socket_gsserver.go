package gstestpipe

import (
	"net"
)

type gstEchoHandlerObj struct {
	readSize, writeSize int64
}

type GstServerSocketEcho struct {
	RawServerSocket
	echo_running        bool
	key                 string
	list_gstEchoHandler []*gstEchoHandlerObj
}

func NewGstServerSocketEcho_RandAddr() *GstServerSocketEcho {
	return &GstServerSocketEcho{RawServerSocket: RawServerSocket{connList: make(chan net.Conn, 1024),
		serverAddr: GetRandAddr()},
		key: g_key_Default}
}

func (ss *GstServerSocketEcho) gstEchoHandler() {
	if ss.echo_running {
		return
	}

	go func() {
		for !ss.close.Load() {
			for conn1 := range ss.connList {
				obj := gstEchoHandlerObj{}
				go gstPipoHandleEx(&obj, ss.key, conn1)
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
