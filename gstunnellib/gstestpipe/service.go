package gstestpipe

import (
	"errors"
	"io"
	"net"
)

type Service_server interface {
	GetClientConn() net.Conn
	GetServerConn() net.Conn
}

type service_serverPiPo struct {
	connList *pipe_conn
}

type service_serverNone struct {
	connList *pipe_conn
}

func NewServiceServerNone() Service_server {
	return &service_serverNone{newPipeConn()}
}

// ping pong
func NewServiceServerPiPo() Service_server {
	return new(service_serverPiPo)
}

func (ss *service_serverNone) GetClientConn() net.Conn {
	return ss.connList.client
}

func (ss *service_serverNone) GetServerConn() net.Conn {
	return ss.connList.server
}

func (ss *service_serverPiPo) GetClientConn() net.Conn {
	if ss.connList != nil {
		return ss.connList.client
	}
	ss.connList = newPipeConn()

	go func(pp *pipe_conn) {
		for {
			_, err := io.Copy(pp.server, pp.server)
			checkError_exit(err)
		}
	}(ss.connList)

	return ss.connList.client
}

func (ss *service_serverPiPo) GetServerConn() net.Conn {
	checkError_panic(errors.New("Service_serverPiPo can not use GetServerConn()."))
	return nil
}
