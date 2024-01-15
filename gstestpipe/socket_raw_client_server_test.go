package gstestpipe

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func rawconnEcho_bug(rawconn net.Conn) {
	defer rawconn.Close()
	rawdata := gsrand.GetRDBytes(1024 * 1024)

	//wg := sync.WaitGroup{}

	_, err := io.Copy(rawconn, bytes.NewBuffer(rawdata))
	checkError_panic(err)
	//fmt.Println(rawdata)

	readdata := make([]byte, len(rawdata))
	_, err = io.ReadAtLeast(rawconn, readdata, len(rawdata))
	checkError_panic(err)
	//fmt.Println(readdata)
	if !bytes.Equal(rawdata, readdata) {
		panic("error")
	}
}

func rawconnEcho(rawconn net.Conn) {
	defer rawconn.Close()
	rawdata := gsrand.GetRDBytes(1024 * 1024)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := io.Copy(rawconn, bytes.NewBuffer(rawdata))
		checkError_panic(err)
		//fmt.Println(rawdata)
	}()

	readdata := make([]byte, len(rawdata))
	_, err := io.ReadAtLeast(rawconn, readdata, len(rawdata))
	checkError_panic(err)
	//fmt.Println(readdata)
	wg.Wait()
	if !bytes.Equal(rawdata, readdata) {
		panic("error")
	}
}

func Test_RawServerSocketEcho(t *testing.T) {
	serverecho := NewRawServerSocketEcho_RandAddr()
	defer serverecho.Close()
	serverecho.Run()

	client, err := net.Dial("tcp4", serverecho.GetServerAddr())
	checkError_panic(err)

	rawconnEcho(client)
}

func Test_RawServerSocketEcho2(t *testing.T) {
	total := 100

	serverecho := NewRawServerSocketEcho_RandAddr()
	defer serverecho.Close()
	serverecho.Run()

	for i := 0; i < total; i++ {

		client, err := net.Dial("tcp4", serverecho.GetServerAddr())
		checkError_panic(err)

		rawconnEcho(client)
	}
}

type connListMt struct {
	list []net.Conn
	lock sync.Mutex
}

func (cls *connListMt) add(conn net.Conn) {
	cls.lock.Lock()
	defer cls.lock.Unlock()

	cls.list = append(cls.list, conn)
}

func noUseTest_RawServerSocketEcho3_bug(t *testing.T) {
	total := 100

	serverecho := NewRawServerSocketEcho_RandAddr()
	defer serverecho.Close()
	serverecho.Run()
	timeout := time.Second * 1

	var clientList connListMt

	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := net.DialTimeout("tcp4", serverecho.GetServerAddr(), timeout)
			checkError_panic(err)
			clientList.add(client)
			rawconnEcho_bug(client)
		}()
	}
	wg.Wait()
	time.Sleep(time.Millisecond * 100)
}

func inTest_RawServerSocketEcho3_mt() {
	total := 100

	serverecho := NewRawServerSocketEcho_RandAddr()
	defer serverecho.Close()
	serverecho.Run()
	timeout := time.Second * 1

	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := net.DialTimeout("tcp4", serverecho.GetServerAddr(), timeout)
			checkError_panic(err)
			rawconnEcho(client)
		}()
	}
	wg.Wait()
}

func Test_RawServerSocketEcho3_mt(t *testing.T) {
	inTest_RawServerSocketEcho3_mt()
}

func noUseTest_RawServerSocketEcho3_mt_n(t *testing.T) {
	total := 1000
	for i := 0; i < total; i++ {
		inTest_RawServerSocketEcho3_mt()
	}
}

func Test_RawClientSocketEcho(t *testing.T) {
	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	clientEcho := NewRawClientSocketEcho(rawserver.GetServerAddr())
	defer clientEcho.Close()
	clientEcho.Run()

	client := <-rawserver.GetConnList()

	rawconnEcho(client)
}

func Test_RawClientSocketEcho2(t *testing.T) {
	total := 100

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	clientEcho := NewRawClientSocketEcho(rawserver.GetServerAddr())
	defer clientEcho.Close()

	for i := 0; i < total; i++ {
		clientEcho.Run()
		client := <-rawserver.GetConnList()
		rawconnEcho(client)
	}
}

func Test_RawClientSocketEcho3(t *testing.T) {
	total := 100

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	clientEcho := NewRawClientSocketEcho(rawserver.GetServerAddr())
	defer clientEcho.Close()
	for i := 0; i < total; i++ {
		clientEcho.Run()
	}

	var wg sync.WaitGroup

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			client := <-rawserver.GetConnList()
			rawconnEcho(client)
		}(i)
	}
	wg.Wait()
}
