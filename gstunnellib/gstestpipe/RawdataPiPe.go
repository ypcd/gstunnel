package gstestpipe

import (
	"errors"
	"io"
	"net"
)

type RawdataPiPe interface {
	GetClientConn() net.Conn
	GetServerConn() net.Conn
}

type rawdataPiPePiPo struct {
	connList *pipe_conn
}

type rawdataPiPeNone struct {
	connList *pipe_conn
}

func NewServiceServerNone() RawdataPiPe {
	return &rawdataPiPeNone{newPipeConn()}
}

// ping pong
func NewServiceServerPiPo() RawdataPiPe {
	return new(rawdataPiPePiPo)
}

func NewSrcClientNone() RawdataPiPe {
	return &rawdataPiPeNone{newPipeConn()}
}

func (ss *rawdataPiPeNone) GetClientConn() net.Conn {
	return ss.connList.client
}

func (ss *rawdataPiPeNone) GetServerConn() net.Conn {
	return ss.connList.server
}

func (ss *rawdataPiPePiPo) GetClientConn() net.Conn {
	if ss.connList != nil {
		return ss.connList.client
	}
	ss.connList = newPipeConn()

	go func(pp *pipe_conn) {

		_, err := io.Copy(pp.server, pp.server)
		checkError_exit(err)

	}(ss.connList)

	return ss.connList.client
}

func (ss *rawdataPiPePiPo) GetServerConn() net.Conn {
	checkError_panic(
		errors.New("RawdataPiPePiPo can not use GetServerConn()"))
	return nil
}
