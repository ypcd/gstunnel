package gstestpipe

import (
	"errors"
	"io"
	"net"
)

type IRawdataPiPe interface {
	GetClientConn() net.Conn
	GetServerConn() net.Conn
}

type rawdataPiPePiPo struct {
	connList *pipe_conn
}

// ping pong
func NewServiceServerPiPo() IRawdataPiPe {
	return new(rawdataPiPePiPo)
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

type rawdataPiPeNone struct {
	connList *pipe_conn
}

// IRawdataPiPe
func NewServiceServerNone() *rawdataPiPeNone {
	return &rawdataPiPeNone{newPipeConn()}
}

func NewServiceServerNoneNoServer() *rawdataPiPeNone {
	p1 := &rawdataPiPeNone{newPipeConn()}
	p1.connList.server = nil
	return p1
}

func NewRawClientNone() *rawdataPiPeNone {
	return &rawdataPiPeNone{newPipeConn()}
}

func (ss *rawdataPiPeNone) GetClientConn() net.Conn {
	return ss.connList.client
}

func (ss *rawdataPiPeNone) GetServerConn() net.Conn {
	return ss.connList.server
}
