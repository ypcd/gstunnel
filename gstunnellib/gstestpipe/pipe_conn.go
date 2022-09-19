package gstestpipe

import "net"

type pipe_conn struct {
	server, client net.Conn
}

func newPipeConn() *pipe_conn {
	pp1 := new(pipe_conn)
	pp1.server, pp1.client = net.Pipe()
	return pp1
}
