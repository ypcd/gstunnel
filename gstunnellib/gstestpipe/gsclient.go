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

var key_defult string = "12345678901234567890123456789012"

type GsPiPe interface {
	GetConn() net.Conn
}

type GsPiPeErrorKey interface {
	GetClientConn() net.Conn
	GetServerConn() net.Conn
}

type gsClientImp struct {
	connList *pipe_conn
	key      string
}

func NewGsPiPe(key string) GsPiPe {
	return &gsClientImp{key: key}
}

func NewGsPiPeDefultKey() GsPiPe {
	return &gsClientImp{key: key_defult}
}

func NewGsPiPeErrorKeyNoKey() GsPiPeErrorKey {
	return &gsClientImp{key: key_defult}
}

func (gc *gsClientImp) GetConn() net.Conn {
	if gc.connList != nil {
		return gc.connList.server
	}
	gc.connList = newPipeConn()

	go func(gc *gsClientImp, pp *pipe_conn) {
		buf := make([]byte, 4*1024)
		var rbuf, wbuf []byte
		var readSize, writeSize int64 = 0, 0

		apacknetRead := gstunnellib.NewGsPackNet(gc.key)
		apacknetWrite := gstunnellib.NewGsPackNet(gc.key)
		for {
			rbuf = buf
			re, err := pp.client.Read(rbuf)
			readSize += int64(re)
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) ||
				errors.Is(err, io.EOF) {
				checkError_info(err)
				logger.Printf("gsclient readSize:%d  writeSize:%d\n", readSize, writeSize)
				return
			}
			checkError_exit(err)

			apacknetRead.WriteEncryData(rbuf[:re])
			wbuf, err = apacknetRead.GetDecryData()
			checkError_exit(err)
			if len(wbuf) > 0 {
				wbuf = apacknetWrite.Packing(wbuf)
				if gsbase.Deep_debug {
					logger.Println("gsclient packing data hash:", gshash.GetSha256Hex(wbuf))
				}
				rn, err := io.Copy(pp.client, bytes.NewBuffer(wbuf))
				writeSize += rn
				if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) ||
					errors.Is(err, io.EOF) {
					checkError_info(err)
					return
				} else {
					checkError_exit(err)
				}
			}
		}
	}(gc, gc.connList)

	return gc.connList.server
}

func (gc *gsClientImp) GetClientConn() net.Conn {
	if gc.connList != nil {
		return gc.connList.client
	}
	gc.connList = newPipeConn()
	return gc.connList.client
}

func (gc *gsClientImp) GetServerConn() net.Conn {
	if gc.connList != nil {
		return gc.connList.server
	}
	gc.connList = newPipeConn()
	return gc.connList.server
}
