package gstunnellib

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

/*
func Test_VersionPack_sendEx(t *testing.T) {
	con1, con2 := GetRDNetConn()
	ap1 := NewGsPack(gsbase.G_AesKeyDefault)
	wlent := int64(0)
	recvbuf := bytes.Buffer{}
	var GetSendbuf bytes.Buffer
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(&recvbuf, con2)
		wg.Done()
	}()
	versionPack_sendEx(con1, ap1, &wlent, &GetSendbuf)
	con1.Close()
	wg.Wait()
	if !(bytes.Equal(GetSendbuf.Bytes(), recvbuf.Bytes()) && GetSendbuf.Len() > 0) {
		t.Fatal("Error.")
	}
	t.Log(recvbuf.String(), recvbuf.Bytes(), recvbuf.Len())
}
*/
/*
	func Test_ChangeCryKey_sendEX(t *testing.T) {
		con1, con2 := GetRDNetConn()
		ap1 := NewGsPack(gsbase.G_AesKeyDefault)
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
		changeCryKey_sendEX(con1, ap1, &tmpv1, &wlent, &GetSendbuf)
		con1.Close()
		wg.Wait()
		if !(bytes.Equal(GetSendbuf.Bytes(), recvbuf.Bytes())) && GetSendbuf.Len() > 0 {
			t.Fatal("Error.")
		}
		t.Log(recvbuf.Bytes(), recvbuf.Len())
	}
*/
func Test_GetNetLocalRDPort(t *testing.T) {
	for i := 0; i < 1000*1; i++ {
		fmt.Println(GetNetLocalRDPort())
	}
}

func Test_GetRDNetConn(t *testing.T) {
	for i := 0; i < 100; i++ {
		conn1, conn2 := GetRDNetConn()
		fmt.Println("conn:", conn1.LocalAddr(), conn2.LocalAddr())
	}
}

func Test_GetRDNetConn2(t *testing.T) {
	wg := sync.WaitGroup{}
	dataLen := 1024 * 1024
	rbuf := make([]byte, dataLen)
	var sendData []byte

	//for i := 0; i < 100; i++ {
	conn1, conn2 := GetRDNetConn()
	wg.Add(1)
	go func() {
		defer wg.Done()
		sendData = gsrand.GetRDBytes(dataLen)
		_, err := NetConnWriteAll(conn1, sendData)
		checkError_panic(err)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.ReadAtLeast(conn2, rbuf, dataLen)
		checkError_panic(err)
	}()

	//}
	wg.Wait()
	if !bytes.Equal(sendData, rbuf) {
		panic("error")
	}
}

func Test_GetRDNetConn3(t *testing.T) {
	wg := sync.WaitGroup{}
	dataLen := 1024 * 1024
	rbuf := make([]byte, dataLen)
	var sendData []byte

	//for i := 0; i < 100; i++ {
	conn1, conn2 := GetRDNetConn()

	time.Sleep(time.Second * 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendData = gsrand.GetRDBytes(dataLen)
		_, err := NetConnWriteAll(conn1, sendData)
		checkError_panic(err)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.ReadAtLeast(conn2, rbuf, dataLen)
		checkError_panic(err)
	}()

	//}
	wg.Wait()
	if !bytes.Equal(sendData, rbuf) {
		panic("error")
	}
}

func Test_GetRDNetConn4(t *testing.T) {
	//wg := sync.WaitGroup{}
	//dataLen := 1024 * 1024
	rbuf := make([]byte, 1024)
	rbuf2 := make([]byte, 1024)

	//var sendData []byte

	//for i := 0; i < 100; i++ {
	conn1, conn2 := GetRDNetConn()

	time.Sleep(time.Second * 3)

	var err1, err2 error

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		conn1.SetReadDeadline(time.Now().Add(time.Second))
		_, err1 = conn1.Read(rbuf)
		checkError_info(err1)
	}()
	go func() {
		defer wg.Done()
		conn2.SetReadDeadline(time.Now().Add(time.Second))
		_, err2 = conn2.Read(rbuf2)
		checkError_info(err2)
	}()
	wg.Wait()
	//time.Sleep(time.Second * 2)
	fmt.Println(err1, err2)

}

func Test_IsNetConnClosed(t *testing.T) {
	c1, c2 := GetRDNetConn()
	c1.Close()
	bl := IsNetConnClosed(c1)
	bl2 := IsNetConnClosed(c2)
	_, _ = bl, bl2

	time.Sleep(time.Second * 1)
	bl = IsNetConnClosed_Read(c1)
	bl2 = IsNetConnClosed_Read(c2)
}
