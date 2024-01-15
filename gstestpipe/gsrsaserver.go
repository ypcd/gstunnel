package gstestpipe

import (
	"net"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type gstRSAServerImp struct {
	gstClientImp

	//KeyRSAServer *gsrsa.RSA
	//KeyRSAClient *gsrsa.RSA
	apacknetRead,
	apacknetWrite gstunnellib.IGSRSAPackNet
}

func NewGSTRSAServerImp(key string, KeyRSAServer *gsrsa.RSA) IGsPiPe {
	rs := &gstRSAServerImp{
		gstClientImp: gstClientImp{key: key},
		//KeyRSAServer:  KeyRSAServer,
		//KeyRSAClient:  KeyRSAClient,
		apacknetRead:  gstunnellib.NewGSRSAPackNetImpWithGSTServer(key, KeyRSAServer),
		apacknetWrite: gstunnellib.NewGSRSAPackNetImpWithGSTServer(key, KeyRSAServer),
	}
	rs._vt = igstPipoHandle(rs)
	//init_exchangeRSAkey()
	return rs
}

func (gc *gstRSAServerImp) ExchangeRSAkey_old(src net.Conn, pack, unpack gstunnellib.IGSRSAPackNet) {
	//pubkpack := clientPack.ClientPublicKeyPack()

	//obj := unpack
	var clientkey *gsrsa.RSA

	rbuf := make([]byte, 4*1024)
	//var wbuf []byte
	var rlent int64
	for {
		//Src.SetReadDeadline(time.Now().Add(NetworkTimeout))
		rlen, err := src.Read(rbuf)
		rlent += int64(rlen)
		//	recot_un_r.Run()
		if gstunnellib.IsErrorNetUsually(err) {
			checkError_info(err)
			return
		} else {
			checkError_panic(err)
		}
		/*
			if tmr_out.Run() {
				g_logger.Printf("Error: [%d] Time out, func exit.\n", Gctx.GetGsId())
				return
			}
		*/
		if rlen == 0 {
			g_logger.Println("Error: Src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			g_logger.Println("Error:", err)
			continue
		}
		//tmr_out.Boot()

		unpack.WriteEncryData(rbuf[:rlen])
		_, err = unpack.GetDecryData()
		checkError_panic(err)
		if unpack.IsExistsClientKey() {
			clientkey = unpack.GetClientRSAKey()
			break
		}
	}
	pack.SetClientRSAKey(clientkey)
	unpack.SetClientRSAKey(clientkey)

}

func (gc *gstRSAServerImp) gstPipoHandle(key string, conn net.Conn) {
	defer conn.Close()

	//gc.ExchangeRSAkey(conn, gc.apacknetWrite, gc.apacknetRead)

	buf := make([]byte, 4*1024)
	var rbuf, wbuf []byte
	var readSize, writeSize int64 = 0, 0

	//	apacknetRead := gstunnellib.NewGSRSAPackNetImp(key, gc.KeyRSAServer, gc.KeyRSAClient)
	//	apacknetWrite := gstunnellib.NewGSRSAPackNetImp(key, gc.KeyRSAServer, gc.KeyRSAClient)
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

		gc.apacknetRead.WriteEncryData(rbuf[:re])
		wbuf, err = gc.apacknetRead.GetDecryData()
		checkError_exit(err)
		if len(wbuf) > 0 {
			wbuf = gc.apacknetWrite.Packing(wbuf)
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
