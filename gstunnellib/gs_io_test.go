package gstunnellib

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_NetConnWriteAll(t *testing.T) {
	testCacheSize := 160 * 1024 * 1024
	SendData := gsrand.GetRDBytes(testCacheSize)
	client, server := net.Pipe()
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
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				CheckError_test_noExit(err, t)
				return
			} else {
				CheckError_test(err, t)
			}
			rbufsz += re

			rbuf = append(rbuf, buf...)

			if rbufsz == testCacheSize {
				return
			}
		}
	}()

	_, err := NetConnWriteAll(client, SendData)
	CheckError_panic(err)

	wg.Wait()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}
}
