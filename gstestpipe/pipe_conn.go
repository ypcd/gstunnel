package gstestpipe

import (
	"net"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

type pipe_conn struct {
	server, client net.Conn
	lst            net.Listener
}

type pipe_connLiten struct {
	pipe_conn
	lst net.Listener
}

/*
func newPipeConn_old() *pipe_conn {
	pp1 := new(pipe_conn)
	pp1.server, pp1.client = net.Pipe()
	return pp1
}
*/

func newPipeConn() *pipe_conn {
	pp1 := new(pipe_conn)
	pp1.server, pp1.client = gstunnellib.GetRDNetConn()
	return pp1
}
