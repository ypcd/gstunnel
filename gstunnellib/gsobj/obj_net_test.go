package gsobj

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

func Test_VersionPack_sendEx(t *testing.T) {
	con1, con2 := gstunnellib.GetRDNetConn()
	//ap1 := gstunnellib.NewGsPack(gsbase.G_AesKeyDefault)
	//wlent := int64(0)
	recvbuf := bytes.Buffer{}
	var GetSendbuf bytes.Buffer

	ap1 := newGSTObj_net_test(con1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(&recvbuf, con2)
		wg.Done()
	}()
	VersionPack_sendEx(con1, ap1, &ap1.Wlent, &GetSendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(GetSendbuf.Bytes(), recvbuf.Bytes()) && GetSendbuf.Len() > 0) {
		t.Fatal("Error.")
	}
	t.Log(recvbuf.String(), recvbuf.Bytes(), recvbuf.Len())
}

func Test_changeCryKey_sendEX_fromClient(t *testing.T) {
	con1, con2 := gstunnellib.GetRDNetConn()
	ap1 := newGSTObj_net_test(con1)
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
	//VersionPack_send(con1, ap1, &wlent)
	changeCryKey_sendEX_fromClient(con1, ap1, &tmpv1, &wlent, &GetSendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(GetSendbuf.Bytes(), recvbuf.Bytes())) && GetSendbuf.Len() > 0 {
		t.Fatal("Error.")
	}
	//t.Log(recvbuf.Bytes(), recvbuf.Len())
}

func Test_changeCryKey_sendEX_fromServer(t *testing.T) {
	con1, con2 := gstunnellib.GetRDNetConn()
	ap1 := newGSTObj_net_test(con1)
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
	//VersionPack_send(con1, ap1, &wlent)
	changeCryKey_sendEX_fromServer(con1, ap1, &tmpv1, &wlent, &GetSendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(GetSendbuf.Bytes(), recvbuf.Bytes())) && GetSendbuf.Len() > 0 {
		t.Fatal("Error.")
	}
	//t.Log(recvbuf.Bytes(), recvbuf.Len())
}
