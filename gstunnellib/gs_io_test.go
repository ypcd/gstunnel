package gstunnellib

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_NetConnWriteAll(t *testing.T) {
	testCacheSize := 160 * 1024 * 1024
	SendData := gsrand.GetRDBytes(testCacheSize)
	client, server := GetRDNetConn()
	testReadTimeOut := time.Second * 1

	rbuf := make([]byte, 0, len(SendData))

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 4096)
		rbufsz := 0

		for {
			server.SetReadDeadline(time.Now().Add(testReadTimeOut))
			re, err := server.Read(buf)
			//t.Logf("server read len: %d", re)
			if IsErrorNetUsually(err) {
				checkError_info(err)
				return
			} else {
				checkError_panic(err)
			}
			rbufsz += re

			rbuf = append(rbuf, buf[:re]...)

			if rbufsz == testCacheSize {
				return
			}
		}
	}()

	n, err := NetConnWriteAll(client, SendData)
	_ = n
	checkError_panic(err)

	wg.Wait()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}
}
