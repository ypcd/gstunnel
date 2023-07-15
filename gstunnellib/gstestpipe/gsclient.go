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

var g_key_Default string = gsbase.G_AesKeyDefault

type GsPiPe interface {
	GetConn() net.Conn
}

type GsPiPeErrorKey interface {
	GetClientConn() net.Conn
	GetServerConn() net.Conn
}

type gstClientImp struct {
	connList   *pipe_conn
	key        string
	pipohandle bool
}

func NewGsPiPe(key string) GsPiPe {
	return &gstClientImp{key: key}
}

func NewGstPiPoDefaultKey() GsPiPe {
	return &gstClientImp{key: g_key_Default}
}

func NewGstPiPoErrorKeyNoKey() GsPiPeErrorKey {
	return &gstClientImp{key: g_key_Default}
}

func gstPipoHandle(key string, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4*1024)
	var rbuf, wbuf []byte
	var readSize, writeSize int64 = 0, 0

	apacknetRead := gstunnellib.NewGsPackNet(key)
	apacknetWrite := gstunnellib.NewGsPackNet(key)
	for {
		rbuf = buf
		re, err := conn.Read(rbuf)
		readSize += int64(re)
		if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) ||
			errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			checkError_info(err)
			g_logger.Printf("gsclient readSize:%d  writeSize:%d\n", readSize, writeSize)
			return
		}
		checkError_exit(err)

		apacknetRead.WriteEncryData(rbuf[:re])
		wbuf, err = apacknetRead.GetDecryData()
		checkError_exit(err)
		if len(wbuf) > 0 {
			wbuf = apacknetWrite.Packing(wbuf)
			if gsbase.G_Deep_debug {
				g_logger.Println("gsclient packing data hash:", gshash.GetSha256Hex(wbuf))
			}
			rn, err := io.Copy(conn, bytes.NewBuffer(wbuf))
			writeSize += rn
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) ||
				errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				checkError_info(err)
				return
			} else {
				checkError_exit(err)
			}
		}
	}
}

func gstPipoHandleEx(obj *gstEchoHandlerObj, key string, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4*1024)
	var rbuf, wbuf []byte
	//var readSize, writeSize int64 = 0, 0

	apacknetRead := gstunnellib.NewGsPackNet(key)
	apacknetWrite := gstunnellib.NewGsPackNet(key)
	for {
		rbuf = buf
		re, err := conn.Read(rbuf)
		obj.readSize += int64(re)
		if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) ||
			errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			checkError_info(err)
			g_logger.Printf("gsclient readSize:%d  writeSize:%d\n", obj.readSize, obj.writeSize)
			return
		}
		checkError_panic(err)

		apacknetRead.WriteEncryData(rbuf[:re])
		wbuf, err = apacknetRead.GetDecryData()
		checkError_exit(err)
		if len(wbuf) > 0 {
			wbuf = apacknetWrite.Packing(wbuf)
			if gsbase.G_Deep_debug {
				g_logger.Println("gsclient packing data hash:", gshash.GetSha256Hex(wbuf))
			}
			rn, err := io.Copy(conn, bytes.NewBuffer(wbuf))
			obj.writeSize += rn
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) ||
				errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				checkError_info(err)
				return
			} else {
				checkError_exit(err)
			}
		}
	}
}

func (gc *gstClientImp) createPiPoHandler() {
	if gc.pipohandle {
		return
	}
	go gstPipoHandle(gc.key, gc.connList.client)
	gc.pipohandle = true
}

func (gc *gstClientImp) GetConn() net.Conn {
	return gc.GetServerConn()
}

func (gc *gstClientImp) GetClientConn() net.Conn {
	if gc.connList != nil {
		return gc.connList.client
	}
	gc.connList = newPipeConn()
	gc.createPiPoHandler()

	return gc.connList.client
}

func (gc *gstClientImp) GetServerConn() net.Conn {
	if gc.connList != nil {
		return gc.connList.server
	}
	gc.connList = newPipeConn()
	gc.createPiPoHandler()

	return gc.connList.server
}
