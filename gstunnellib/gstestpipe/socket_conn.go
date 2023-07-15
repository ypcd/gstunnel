package gstestpipe

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

func getnetconn() (net.Conn, net.Conn) {
	var conna, connc net.Conn

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		lst, err := net.Listen("tcp4", "127.0.0.1:12345")
		if err != nil {
			fmt.Println("conna, connc:", conna, connc)
		}
		checkError_exit(err)
		defer lst.Close()
		conna, err = lst.Accept()
		checkError_exit(err)
	}()
	//	time.Sleep(1 * time.Second)
	connc, err := net.Dial("tcp4", "127.0.0.1:12345")
	checkError_exit(err)
	wg.Wait()
	return conna, connc
}

func newSocketConn(server string) *pipe_conn {
	pp1 := new(pipe_conn)
	var err error
	pp1.client, err = net.Dial("tcp4", server)
	checkError_panic(err)
	return pp1
}

func GetRandAddr() string {
	port := 1024 + gsrand.GetRDCInt_max(49151-1024)
	return "127.0.0.1:" + strconv.Itoa(int(port))
}
