package gstestpipe

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
)

var key_defult string = "1234567890123456"

type GsClient interface {
	GetServerConn() net.Conn
}

type gsClientImp struct {
	connList *pipe_conn
	key      string
}

func NewGsClient(key string) GsClient {
	return &gsClientImp{key: key}
}

func NewGsClientDefultKey() GsClient {
	return &gsClientImp{key: key_defult}
}

func (gc *gsClientImp) GetServerConn() net.Conn {
	if gc.connList != nil {
		return gc.connList.server
	}
	gc.connList = newPipeConn()

	go func(gc *gsClientImp, pp *pipe_conn) {
		buf := make([]byte, 2*1024)
		var rbuf, wbuf []byte
		//var rpackbuf []byte

		apacknetRead := gstunnellib.NewGsPackNet(gc.key)
		apacknetWrite := gstunnellib.NewGsPackNet(gc.key)
		for {
			rbuf = buf
			re, err := pp.client.Read(rbuf)
			checkError(err)
			if err != nil {
				return
			}

			apacknetRead.WriteEncryData(rbuf[:re])
			wbuf, err = apacknetRead.GetDecryData()
			checkError_exit(err)
			if len(wbuf) > 0 {
				wbuf = apacknetWrite.Packing(wbuf)
				if gsbase.Deep_debug {
					logger.Println("gsclient packing data hash:", gshash.GetSha256Hex(wbuf))
				}
				_, err := io.Copy(pp.client, bytes.NewBuffer(wbuf))
				if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, io.EOF) {
					checkError(err)
					return
				} else {
					checkError_exit(err)
				}
			}
		}
	}(gc, gc.connList)

	return gc.connList.server
}
