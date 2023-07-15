package gstestpipe

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gshash"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_GsPiPeImp(t *testing.T) {
	gc := NewGstPiPoDefaultKey()
	wg := sync.WaitGroup{}

	SendData := gsrand.GetRDBytes(50 * 1024)
	//SendData := []byte("123456")

	var rbuf []byte
	//rbuff := bytes.Buffer{}
	apack := gstunnellib.NewGsPackNet(g_key_Default)

	server := gc.GetConn()
	wg.Add(1)
	go func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, 2*1024)

		for {
			re, err := server.Read(buf)
			t.Logf("server read len: %d", re)
			CheckError_test_noExit(err, t)
			if err != nil && re == 0 {
				return
			}

			fmt.Println("hash:", gshash.GetSha256Hex(buf[:re]))
			apack.WriteEncryData(buf[:re])
			//_, err := pp.server.Read(buf)
			//CheckError_test(err, t)
			//pp.client
			data, err := apack.GetDecryData()
			CheckError_test(err, t)
			if data != nil {
				rbuf = append(rbuf, data...)
			}
		}

	}(server)

	wbuf := apack.Packing(SendData)
	server.Write(wbuf)
	time.Sleep(time.Second * 3)
	server.Close()
	wg.Wait()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}

}
