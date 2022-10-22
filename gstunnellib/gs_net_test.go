package gstunnellib

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
)

func Test_IsTheVersionConsistent_sendEx(t *testing.T) {
	con1, con2 := net.Pipe()
	ap1 := NewGsPack("12345678901234567890123456789012")
	wlent := int64(0)
	recvbuf := bytes.Buffer{}
	var GetSendbuf bytes.Buffer
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(&recvbuf, con2)
		wg.Done()
	}()
	isTheVersionConsistent_sendEx(con1, ap1, &wlent, &GetSendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(GetSendbuf.Bytes(), recvbuf.Bytes()) && GetSendbuf.Len() > 0) {
		t.Fatal("Error.")
	}
	t.Log(recvbuf.String(), recvbuf.Bytes(), recvbuf.Len())
}

func Test_ChangeCryKey_sendEX(t *testing.T) {
	con1, con2 := net.Pipe()
	ap1 := NewGsPack("12345678901234567890123456789012")
	wlent := int64(0)
	recvbuf := bytes.Buffer{}
	var GetSendbuf bytes.Buffer
	tmpv1 := 0
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(&recvbuf, con2)
		wg.Done()
	}()
	//IsTheVersionConsistent_send(con1, ap1, &wlent)
	changeCryKey_sendEX(con1, ap1, &tmpv1, &wlent, &GetSendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(GetSendbuf.Bytes(), recvbuf.Bytes())) && GetSendbuf.Len() > 0 {
		t.Fatal("Error.")
	}
	t.Log(recvbuf.Bytes(), recvbuf.Len())
}
