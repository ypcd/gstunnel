package gstestpipe

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func gstconnEcho(gstconn net.Conn, key string) {
	rawdata := gsrand.GetRDBytes(1024)

	apacknetRead := gstunnellib.NewGsPackNet(key)
	apacknetWrite := gstunnellib.NewGsPackNet(key)
	data := apacknetWrite.Packing(rawdata)
	size, err := io.Copy(gstconn, bytes.NewBuffer(data))
	checkError_panic(err)
	fmt.Println("write size:", size)
	//fmt.Println(len(data), data)

	sizedata := make([]byte, 2)
	_, err = io.ReadAtLeast(gstconn, sizedata, 2)
	checkError_panic(err)
	sz := apacknetRead.GetGSPackSize(sizedata)
	fmt.Println("read body size:", sz)

	readbody := make([]byte, sz)
	_, err = io.ReadAtLeast(gstconn, readbody, len(readbody))
	checkError_panic(err)

	apacknetRead.WriteEncryData(append(sizedata, readbody...))
	dedata, _ := apacknetRead.GetDecryData()
	//fmt.Println(string(dedata))

	if !bytes.Equal(rawdata, dedata) {
		panic("error")
	}
}

func Test_GstServerSocketEcho(t *testing.T) {
	serverecho := NewGstServerSocketEcho_RandAddr()
	defer serverecho.Close()
	serverecho.Run()
	fmt.Println("server addr:", serverecho.GetServerAddr())

	client, err := net.Dial("tcp4", serverecho.GetServerAddr())
	checkError_panic(err)
	fmt.Println(client.LocalAddr(), client.RemoteAddr())

	gstconnEcho(client, serverecho.GetKey())
}

func Test_GstServerSocketEcho2(t *testing.T) {
	loopNum := 100

	serverecho := NewGstServerSocketEcho_RandAddr()
	defer serverecho.Close()
	serverecho.Run()
	fmt.Println("server addr:", serverecho.GetServerAddr())

	for i := 0; i < loopNum; i++ {
		client, err := net.Dial("tcp4", serverecho.GetServerAddr())
		checkError_panic(err)
		fmt.Println(client.LocalAddr(), client.RemoteAddr())

		gstconnEcho(client, serverecho.GetKey())

	}
}

func Test_GstServerSocketEcho3(t *testing.T) {
	loopNum := 100

	serverecho := NewGstServerSocketEcho_RandAddr()
	defer serverecho.Close()
	serverecho.Run()
	fmt.Println("server addr:", serverecho.GetServerAddr())

	wg := sync.WaitGroup{}

	for i := 0; i < loopNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client, err := net.Dial("tcp4", serverecho.GetServerAddr())
			checkError_panic(err)
			fmt.Println(client.LocalAddr(), client.RemoteAddr())

			gstconnEcho(client, serverecho.GetKey())

		}()
	}
	wg.Wait()
}

func Test_GstClientSocketEcho1(t *testing.T) {

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	clientEcho := NewGstClientSocketEcho_defaultKey(rawserver.GetServerAddr())
	defer clientEcho.Close()

	clientEcho.Run()

	client := <-rawserver.GetConnList()

	fmt.Println(client.LocalAddr(), client.RemoteAddr())

	gstconnEcho(client, clientEcho.GetKey())

}

func Test_GstClientSocketEcho2(t *testing.T) {
	total := 100

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	clientEcho := NewGstClientSocketEcho_defaultKey(rawserver.GetServerAddr())
	defer clientEcho.Close()
	for i := 0; i < total; i++ {
		clientEcho.Run()
	}

	for i := 0; i < total; i++ {
		client := <-rawserver.GetConnList()

		fmt.Println(client.LocalAddr(), client.RemoteAddr())

		gstconnEcho(client, clientEcho.GetKey())
	}
}

func Test_GstClientSocketEcho3(t *testing.T) {
	total := 100

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	clientEcho := NewGstClientSocketEcho_defaultKey(rawserver.GetServerAddr())
	defer clientEcho.Close()
	for i := 0; i < total; i++ {
		clientEcho.Run()
	}

	wg := sync.WaitGroup{}

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := <-rawserver.GetConnList()

			fmt.Println(client.LocalAddr(), client.RemoteAddr())

			gstconnEcho(client, clientEcho.GetKey())
		}()
	}
	wg.Wait()
}
