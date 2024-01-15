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
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

func gstconnEchoRSAClient(gstconn net.Conn, key string, serverRSAKey, clientRSAKey *gsrsa.RSA) {
	newkey := gstclient_ChangeCryKeyFromGSTClient_send(gstconn, gstunnellib.NewGSRSAPackNetImpWithGSTClient(key, serverRSAKey, clientRSAKey))

	rawdata := gsrand.GetRDBytes(1024)

	apacknetRead := gstunnellib.NewGSRSAPackNetImpWithGSTClient(key, serverRSAKey, clientRSAKey)
	apacknetWrite := gstunnellib.NewGSRSAPackNetImpWithGSTClient(newkey, serverRSAKey, clientRSAKey)
	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	func() {
		defer wg.Done()
		data := apacknetWrite.Packing(rawdata)
		size, err := io.Copy(gstconn, bytes.NewBuffer(data))
		checkError_panic(err)
		fmt.Println("write size:", size)
		//fmt.Println(len(data), data)
	}()

	var dedata []byte
	for {
		sizedata := make([]byte, 2)
		_, err := io.ReadAtLeast(gstconn, sizedata, 2)
		checkError_panic(err)
		sz := apacknetRead.GetGSPackSize(sizedata)
		fmt.Println("read body size:", sz)

		readbody := make([]byte, sz)
		_, err = io.ReadAtLeast(gstconn, readbody, len(readbody))
		checkError_panic(err)

		apacknetRead.WriteEncryData(append(sizedata, readbody...))
		dedata, _ = apacknetRead.GetDecryData()
		//fmt.Println(string(dedata))
		if len(dedata) == len(rawdata) {
			break
		}
	}
	if !bytes.Equal(rawdata, dedata) {
		panic("error")
	}
}

func Test_GstRSAServerSocketEcho(t *testing.T) {
	//key1 := gsrand.GetRDBytes(gsbase.G_AesKeyLen)
	serverRSAKey := gsrsa.NewGenRSAObj(4096)
	clientRSAKey := gsrsa.NewGenRSAObj(4096)

	//var clientRSAKey := nil

	GSTserverecho := NewGstRSAServerSocketEcho_RandAddr(serverRSAKey)
	defer GSTserverecho.Close()
	GSTserverecho.Run()
	g_logger.Println("server addr:", GSTserverecho.GetServerAddr())

	gstclient_init_exchangeRSAkey(GSTserverecho.GetServerAddr(), GSTserverecho.GetKey(), serverRSAKey, clientRSAKey)

	client, err := net.Dial("tcp4", GSTserverecho.GetServerAddr())
	checkError_panic(err)
	g_logger.Println(client.LocalAddr(), client.RemoteAddr())

	gstconnEchoRSAClient(client, GSTserverecho.GetKey(), serverRSAKey, clientRSAKey)
}

func Test_GstRSAServerSocketEcho2(t *testing.T) {
	const loopNum = 30

	serverRSAKey := gsrsa.NewGenRSAObj(4096)
	clientRSAKey := gsrsa.NewGenRSAObj(4096)

	serverecho := NewGstRSAServerSocketEcho_RandAddr(serverRSAKey)
	defer serverecho.Close()
	serverecho.Run()
	fmt.Println("server addr:", serverecho.GetServerAddr())

	for i := 0; i < loopNum; i++ {
		gstclient_init_exchangeRSAkey(serverecho.GetServerAddr(), serverecho.GetKey(), serverRSAKey, clientRSAKey)

		client, err := net.Dial("tcp4", serverecho.GetServerAddr())
		checkError_panic(err)
		fmt.Println(client.LocalAddr(), client.RemoteAddr())

		gstconnEchoRSAClient(client, serverecho.GetKey(), serverRSAKey, clientRSAKey)

	}
}

func Test_GstRSAServerSocketEcho3(t *testing.T) {
	loopNum := 30

	serverRSAKey := gsrsa.NewGenRSAObj(4096)
	clientRSAKey := gsrsa.NewGenRSAObj(4096)

	serverecho := NewGstRSAServerSocketEcho_RandAddr(serverRSAKey)
	defer serverecho.Close()
	serverecho.Run()
	fmt.Println("server addr:", serverecho.GetServerAddr())

	wg := sync.WaitGroup{}

	for i := 0; i < loopNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gstclient_init_exchangeRSAkey(serverecho.GetServerAddr(), serverecho.GetKey(), serverRSAKey, clientRSAKey)

			client, err := net.Dial("tcp4", serverecho.GetServerAddr())
			checkError_panic(err)
			fmt.Println(client.LocalAddr(), client.RemoteAddr())

			gstconnEchoRSAClient(client, serverecho.GetKey(), serverRSAKey, clientRSAKey)

		}()
	}
	wg.Wait()
}

/*
func gstconnEchoRSAServer(gstconn net.Conn, key string, serverRSAKey, clientRSAKey *gsrsa.RSA) {
	rawdata := gsrand.GetRDBytes(1024)

	apacknetRead := gstunnellib.NewGSRSAPackNetImpWithGSTServer(key, serverRSAKey)
	apacknetWrite := gstunnellib.NewGSRSAPackNetImpWithGSTServer(key, serverRSAKey)
	data := apacknetWrite.Packing(rawdata)
	size, err := io.Copy(gstconn, bytes.NewBuffer(data))
	checkError_panic(err)
	fmt.Println("write size:", size)
	//fmt.Println(len(data), data)

	var dedata []byte
	for {
		sizedata := make([]byte, 2)
		_, err = io.ReadAtLeast(gstconn, sizedata, 2)
		checkError_panic(err)
		sz := apacknetRead.GetGSPackSize(sizedata)
		fmt.Println("read body size:", sz)

		readbody := make([]byte, sz)
		_, err = io.ReadAtLeast(gstconn, readbody, len(readbody))
		checkError_panic(err)

		apacknetRead.WriteEncryData(append(sizedata, readbody...))
		dedata, _ = apacknetRead.GetDecryData()
		//fmt.Println(string(dedata))
		if len(dedata) == len(rawdata) {
			break
		}
	}
	if !bytes.Equal(rawdata, dedata) {
		panic("error")
	}
}

func Test_GstRSAClientSocketEcho1(t *testing.T) {

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	serverRSAKey := gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate))
	clientRSAKey := gsrsa.NewGenRSAObj(4096)

	clientEcho := NewGstRSAClientSocketEcho_defaultKey(rawserver.GetServerAddr(), clientRSAKey)
	defer clientEcho.Close()

	clientEcho.Run()

	client := <-rawserver.GetConnList()

	fmt.Println(client.LocalAddr(), client.RemoteAddr())

	gstconnEchoRSAServer(client, clientEcho.GetKey(), serverRSAKey, clientRSAKey)

}

func Test_GstRSAClientSocketEcho2(t *testing.T) {
	total := 100

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	serverRSAKey := gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate))
	clientRSAKey := gsrsa.NewGenRSAObj(4096)

	clientEcho := NewGstRSAClientSocketEcho_defaultKey(rawserver.GetServerAddr(), clientRSAKey)
	defer clientEcho.Close()
	for i := 0; i < total; i++ {
		clientEcho.Run()
	}

	for i := 0; i < total; i++ {
		client := <-rawserver.GetConnList()

		fmt.Println(client.LocalAddr(), client.RemoteAddr())

		gstconnEchoRSAServer(client, clientEcho.GetKey(), serverRSAKey, clientRSAKey)
	}
}

func Test_GstRSAClientSocketEcho3(t *testing.T) {
	total := 100

	rawserver := NewRawServerSocket_RandAddr()
	defer rawserver.Close()
	rawserver.Run()

	serverRSAKey := gsrsa.NewRSAObjFromBase64([]byte(gsbase.G_DefaultRSAKeyPrivate))
	clientRSAKey := gsrsa.NewGenRSAObj(4096)

	clientEcho := NewGstRSAClientSocketEcho_defaultKey(rawserver.GetServerAddr(), clientRSAKey)
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

			gstconnEchoRSAServer(client, clientEcho.GetKey(), serverRSAKey, clientRSAKey)
		}()
	}
	wg.Wait()
}
*/
