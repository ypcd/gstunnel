package gstestpipe

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsbase"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func in_ExchangeRSAkey_client(dst net.Conn, pack, unpack gstunnellib.IGSRSAPackNet) {
	pubpack := pack.ClientPublicKeyPack()
	_, err := gstunnellib.NetConnWriteAll(dst, pubpack)
	checkError_panic(err)
}

func Test_gstRSAServerImp(t *testing.T) {
	key1 := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	serverrsakey := gsrsa.NewGenRSAObj(4096)
	clientrsakey := gsrsa.NewGenRSAObj(4096)

	gc := NewGSTRSAServerImp(string(key1), serverrsakey)
	wg := sync.WaitGroup{}

	SendData := gsrand.GetRDBytes(50 * 1024)
	//SendData := []byte("123456")

	var rbuf []byte
	//rbuff := bytes.Buffer{}
	pack := gstunnellib.NewGSRSAPackNetImpWithGSTClient(string(key1), serverrsakey, clientrsakey)
	unpack := gstunnellib.NewGSRSAPackNetImpWithGSTClient(string(key1), serverrsakey, clientrsakey)

	server := gc.GetConn()
	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, 2*1024)

		for {
			re, err := server.Read(buf)
			t.Logf("server read len: %d", re)
			checkError_info(err)
			if err != nil && re == 0 {
				return
			}

			fmt.Println("hash:", gshash.GetSha256Hex(buf[:re]))
			unpack.WriteEncryData(buf[:re])
			//_, err := pp.server.Read(buf)
			//CheckError_test(err, t)
			//pp.client
			data, err := unpack.GetDecryData()
			CheckError_test(err, t)
			if data != nil {
				rbuf = append(rbuf, data...)
			}
			if len(rbuf) == len(SendData) {
				return
			}
		}

	}(server)

	in_ExchangeRSAkey_client(server, pack, unpack)
	changekeypack, _ := pack.ChangeCryKeyFromGSTClient()
	gstunnellib.NetConnWriteAll(server, changekeypack)
	wbuf := pack.Packing(SendData)
	gstunnellib.NetConnWriteAll(server, wbuf)
	//	time.Sleep(time.Second * 3)
	wg.Wait()
	server.Close()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}

}
