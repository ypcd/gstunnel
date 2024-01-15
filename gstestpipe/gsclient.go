package gstestpipe

import (
	"net"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
)

var g_key_Default string = gsbase.G_AesKeyDefault

type IGsPiPe interface {
	GetConn() net.Conn
}

type igstPipoHandle interface {
	gstPipoHandle(key string, conn net.Conn)
}

type IGsPiPeErrorKey interface {
	GetClientConn() net.Conn
	GetServerConn() net.Conn
}

type gstClientImp struct {
	connList   *pipe_conn
	key        string
	pipohandle bool
	_vt        igstPipoHandle
}

func newGSTClientImp(key string) *gstClientImp {
	return &gstClientImp{key: key}
}

func NewGsPiPe(key string) IGsPiPe {
	return newGSTClientImp(key)
}

func NewGstPiPoDefaultKey() IGsPiPe {
	return &gstClientImp{key: g_key_Default}
}

func NewGstPiPoErrorKeyNoKey() IGsPiPeErrorKey {
	return &gstClientImp{key: g_key_Default}
}

func (gc *gstClientImp) gstPipoHandle(key string, conn net.Conn) {
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
		if gstunnellib.IsErrorNetUsually(err) {
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
			rn, err := gstunnellib.NetConnWriteAll(conn, wbuf)
			writeSize += rn
			if gstunnellib.IsErrorNetUsually(err) {
				checkError_info(err)
				return
			} else {
				checkError_exit(err)
			}
		}
	}
}

/*
	func gstPipoHandleEx_old(obj *gstEchoHandlerObj, key string, conn net.Conn) {
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
			if gstunnellib.IsErrorNetUsually(err) {
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
				rn, err := gstunnellib.NetConnWriteAll(conn, wbuf)
				obj.writeSize += rn
				if gstunnellib.IsErrorNetUsually(err) {
					checkError_info(err)
					return
				} else {
					checkError_exit(err)
				}
			}
		}
	}
*/
func gstPipoHandleEx(obj *gstEchoHandlerObj, key string, conn net.Conn, apacknetRead, apacknetWrite gstunnellib.IGSRSAPackNet) {
	defer conn.Close()

	buf := make([]byte, 4*1024)
	var rbuf, wbuf []byte
	//var readSize, writeSize int64 = 0, 0

	//apacknetRead := gstunnellib.NewGSRSAPackNetImp(key, serverRSAKey, clientRSAKey)
	//apacknetWrite := gstunnellib.NewGSRSAPackNetImp(key, serverRSAKey, clientRSAKey)

	for {
		rbuf = buf
		re, err := conn.Read(rbuf)
		obj.readSize += int64(re)
		if gstunnellib.IsErrorNetUsually(err) {
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
			rn, err := gstunnellib.NetConnWriteAll(conn, wbuf)
			obj.writeSize += rn
			if gstunnellib.IsErrorNetUsually(err) {
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
	if gc._vt != nil {
		go gc._vt.gstPipoHandle(gc.key, gc.connList.client)
	} else {
		go gc.gstPipoHandle(gc.key, gc.connList.client)
	}
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
