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
	ap1 := NewGsPack("1234567890123456")
	wlent := int64(0)
	rbuf := bytes.Buffer{}
	var sendbuf bytes.Buffer
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(&rbuf, con2)
		wg.Done()
	}()
	//IsTheVersionConsistent_send(con1, ap1, &wlent)
	IsTheVersionConsistent_sendEx(con1, ap1, &wlent, &sendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(sendbuf.Bytes(), rbuf.Bytes())) && sendbuf.Len() > 0 {
		t.Fatal("Error.")
	}
	t.Log(rbuf.Bytes(), rbuf.Len())
}

func Test_ChangeCryKey_sendEX(t *testing.T) {
	con1, con2 := net.Pipe()
	ap1 := NewGsPack("1234567890123456")
	wlent := int64(0)
	rbuf := bytes.Buffer{}
	var sendbuf bytes.Buffer
	tmpv1 := 0
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(&rbuf, con2)
		wg.Done()
	}()
	//IsTheVersionConsistent_send(con1, ap1, &wlent)
	ChangeCryKey_sendEX(con1, ap1, &tmpv1, &wlent, &sendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(sendbuf.Bytes(), rbuf.Bytes())) && sendbuf.Len() > 0 {
		t.Fatal("Error.")
	}
	t.Log(rbuf.Bytes(), rbuf.Len())
}
