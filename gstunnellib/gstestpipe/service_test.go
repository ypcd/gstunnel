package gstestpipe

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func Test_pipeConn(t *testing.T) {
	const sendSize int = 1024 * 1024
	SendData := gsrand.GetRDBytes(sendSize)

	pp1 := newPipeConn()

	rbuf := make([]byte, 1024*1024)
	//rbuf := bytes.Buffer{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(pp *pipe_conn) {
		defer wg.Done()

		re, err := io.ReadAtLeast(pp.server, rbuf, sendSize)
		t.Logf("server read len: %d", re)
		CheckError_test(err, t)
	}(pp1)

	re, err := pp1.client.Write(SendData)
	_ = re
	CheckError_test(err, t)
	//pp1.client.Close()
	wg.Wait()

	rdata := rbuf
	_ = rdata

	sendlen := len(SendData)
	rbuflen := len(rdata)
	t.Logf("sendLen: %d, recvLen: %d", sendlen, rbuflen)

	bl := string(SendData) == string(rdata)
	_ = bl

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}
}

func Test_netpipe(t *testing.T) {
	pp1 := newPipeConn()
	wg := sync.WaitGroup{}

	SendData := gsrand.GetRDBytes(1024 * 1024)
	rbuf := make([]byte, len(SendData))
	//rbuff := bytes.Buffer{}

	go func(pp *pipe_conn) {
		//buf := make([]byte, 1024*1024)
		for {
			_, err := io.Copy(pp.server, pp.server)
			CheckError_test(err, t)
			//_, err := pp.server.Read(buf)
			//CheckError_test(err, t)
			//pp.client
		}
	}(pp1)

	wg.Add(1)
	go func(pp *pipe_conn) {
		defer wg.Done()
		//buf := make([]byte, 1024*1024)
		re, err := io.ReadAtLeast(pp.client, rbuf, len(rbuf))
		t.Logf("server read len: %d", re)
		CheckError_test(err, t)
		//_, err := pp.server.Read(buf)
		//CheckError_test(err, t)
		//pp.client

	}(pp1)

	pp1.client.Write(SendData)
	wg.Wait()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}
}

func Test_servicePiPo(t *testing.T) {
	ss := NewServiceServerPiPo()
	wg := sync.WaitGroup{}

	SendData := gsrand.GetRDBytes(1024 * 1024)
	rbuf := make([]byte, len(SendData))
	//rbuff := bytes.Buffer{}

	client := ss.GetClientConn()

	wg.Add(1)
	go func(client net.Conn) {
		defer wg.Done()
		//buf := make([]byte, 1024*1024)
		re, err := io.ReadAtLeast(client, rbuf, len(rbuf))
		t.Logf("server read len: %d", re)
		CheckError_test(err, t)
		//_, err := pp.server.Read(buf)
		//CheckError_test(err, t)
		//pp.client

	}(client)

	client.Write(SendData)
	wg.Wait()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}

}

func Test_serviceNone(t *testing.T) {
	ss := NewServiceServerNone()
	wg := sync.WaitGroup{}

	SendData := gsrand.GetRDBytes(1024 * 1024)
	rbuf := make([]byte, len(SendData))

	server := ss.GetServerConn()

	wg.Add(1)
	go func(client net.Conn) {
		defer wg.Done()
		//buf := make([]byte, 1024*1024)
		re, err := io.ReadAtLeast(client, rbuf, len(rbuf))
		t.Logf("server read len: %d", re)
		CheckError_test(err, t)

	}(ss.GetClientConn())

	server.Write(SendData)
	wg.Wait()

	if !bytes.Equal(SendData, rbuf) {
		t.Fatal("Error: SendData != rbuf.")
	}

}
