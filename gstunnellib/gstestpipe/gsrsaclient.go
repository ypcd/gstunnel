package gstestpipe

import (
	"net"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type gstRSAClientImp struct {
	gstClientImp

	//KeyRSAServer *gsrsa.RSA
	//KeyRSAClient *gsrsa.RSA
	apacknetRead,
	apacknetWrite gstunnellib.IGSRSAPackNet
}

func newGSTRSAClientImp(key string, KeyRSAServer, KeyRSAClient *gsrsa.RSA) *gstRSAClientImp {
	rs := &gstRSAClientImp{
		gstClientImp: *newGSTClientImp(key),
		//KeyRSAServer:  KeyRSAServer,
		//KeyRSAClient:  KeyRSAClient,
		apacknetRead:  gstunnellib.NewGSRSAPackNetImpWithGSTClient(key, KeyRSAServer, KeyRSAClient),
		apacknetWrite: gstunnellib.NewGSRSAPackNetImpWithGSTClient(key, KeyRSAServer, KeyRSAClient),
	}
	rs._vt = igstPipoHandle(rs)
	return rs
}

func NewGSTRSAClientImp(key string, KeyRSAServer, KeyRSAClient *gsrsa.RSA) IGsPiPe {
	return newGSTRSAClientImp(key, KeyRSAServer, KeyRSAClient)
}

func NewGSTRSAClientImp_DefaultKey() IGsPiPe {
	return newGSTRSAClientImp(
		gsbase.G_AesKeyDefault,
		gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
		gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen),
	)
}

func NewGstRSAPiPoErrorKeyNoKey() IGsPiPeErrorKey {
	return newGSTRSAClientImp(
		gsbase.G_AesKeyDefault,
		gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate)),
		gsrsa.NewGenRSAObj(gsbase.G_RSAKeyBitLen),
	)
}

/*
func (ss *gstRSAClientImp) ExchangeRSAkey_old(dst net.Conn, pack, unpack gstunnellib.IGSRSAPackNet) {

	//unpack := gstunnellib.NewGSRSAPackNetImp(ss.key, ss.serverRSAKey, ss.clientRSAKey)
	//pack := gstunnellib.NewGSRSAPackNetImp(ss.key, ss.serverRSAKey, ss.clientRSAKey)

	pubpack := pack.ClientPublicKeyPack()
	_, err := gstunnellib.NetConnWriteAll(dst, pubpack)
	checkError_panic(err)

	//go gstPipoHandleEx(obj, ss.key, src, unpack, pack)

}*/

func (gc *gstRSAClientImp) gstPipoHandle(key string, conn net.Conn) {
	defer conn.Close()

	gstclient_exchangeRSAkey_send(conn, gc.apacknetWrite)

	var rbuf, wbuf []byte
	var readSize, writeSize int64 = 0, 0

	rbuf = make([]byte, 4*1024)

	//	apacknetRead := gstunnellib.NewGSRSAPackNetImp(key, gc.KeyRSAServer, gc.KeyRSAClient)
	//	apacknetWrite := gstunnellib.NewGSRSAPackNetImp(key, gc.KeyRSAServer, gc.KeyRSAClient)
	for {
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
