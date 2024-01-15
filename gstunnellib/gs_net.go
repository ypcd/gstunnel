package gstunnellib

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gsnetconn"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrand"
)

const g_debug_netconn_isclosed bool = false

/*
	func init() {
		fmt.Println("gs_net init().")
	}
*/
/*
func versionPack_sendEx(dst net.Conn, apack *gsnetconn.GSTObj, wlent *int64, GetSendbuf *bytes.Buffer) error {

	if true {
		buf := apack.PackVersion()

		if GetSendbuf != nil {
			io.Copy(GetSendbuf, bytes.NewBuffer(buf))
		}

		wlen, err := NetConnWriteAll(dst, buf)
		*wlent += int64(wlen)
		if IsErrorNetUsually(err) {
			checkError_panic(err)
			return err
		} else {
			checkError_panic(err)
		}
	}
	return nil
}

func VersionPack_send(obj *gsnetconn.GSTObj) error {
	return versionPack_sendEx(obj.Dst, obj, &obj.Wlent, nil)
}


func changeCryKey_sendEX_fromClient(dst net.Conn, apack *gsnetconn.GSTObj, ChangeCryKey_Total *int, wlent *int64, GetSendbuf *bytes.Buffer) error {

	buf := apack.ChangeCryKeyFromGSTClient()

	if GetSendbuf != nil {
		io.Copy(GetSendbuf, bytes.NewReader(buf))
	}
	//tmr.Boot()
	*ChangeCryKey_Total += 1
	//outf2.Write(buf)
	wlen, err := NetConnWriteAll(dst, buf)
	*wlent += int64(wlen)
	if IsErrorNetUsually(err) {
		checkError_panic(err)
		return err
	} else {
		checkError_panic(err)
	}
	return nil
}

func changeCryKey_sendEX_fromServer(dst net.Conn, apack *gsnetconn.GSTObj, ChangeCryKey_Total *int, wlent *int64, GetSendbuf *bytes.Buffer) error {

	buf := apack.ChangeCryKeyFromGSTServer()

	if GetSendbuf != nil {
		io.Copy(GetSendbuf, bytes.NewReader(buf))
	}
	//tmr.Boot()
	*ChangeCryKey_Total += 1
	//outf2.Write(buf)
	wlen, err := NetConnWriteAll(dst, buf)
	*wlent += int64(wlen)
	if IsErrorNetUsually(err) {
		checkError_panic(err)
		return err
	} else {
		checkError_panic(err)
	}
	return nil
}


func ChangeCryKey_send_fromClient(obj *gsnetconn.GSTObj) error {
	return changeCryKey_sendEX_fromClient(obj.Dst, obj, &obj.ChangeCryKey_Total, &obj.Wlent, nil)
}

func ChangeCryKey_send_fromServer(obj *gsnetconn.GSTObj) error {
	return changeCryKey_sendEX_fromServer(obj.Dst, obj, &obj.ChangeCryKey_Total, &obj.Wlent, nil)
}

func GetNetLocalRDPort() string {
	return fmt.Sprintf("127.0.0.1:%d", gsrand.GetRD_netPortNumber())
}

func GetRDNetConn() (net.Conn, net.Conn) {
	var conna, connc net.Conn

	serverAddr := ""

	wg := sync.WaitGroup{}
	listenDone := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()

		lst, err := net.Listen("tcp4", "127.0.0.1:")
		checkError_exit(err)
		serverAddr = lst.Addr().String()

		listenDone <- 1
		defer lst.Close()
		conna, err = lst.Accept()
		checkError_exit(err)
	}()
	<-listenDone
	time.Sleep(time.Millisecond * 10)
	connc, err := net.Dial("tcp4", serverAddr)
	checkError_exit(err)
	wg.Wait()
	return conna, connc
}
*/

/*
func changeCryKey_sendEX(dst net.Conn, apack IGSPack, ChangeCryKey_Total *int, wlent *int64, GetSendbuf *bytes.Buffer) error {

		buf := apack.ChangeCryKey()

		if GetSendbuf != nil {
			io.Copy(GetSendbuf, bytes.NewReader(buf))
		}
		//tmr.Boot()
		*ChangeCryKey_Total += 1
		//outf2.Write(buf)
		wlen, err := NetConnWriteAll(dst, buf)
		*wlent += int64(wlen)
		if IsErrorNetUsually(err) {
			checkError_panic(err)
			return err
		} else {
			checkError_panic(err)
		}
		return nil
	}
*/

/*
	func ChangeCryKey_send(dst net.Conn, apack IGSPack, ChangeCryKey_Total *int, wlent *int64) error {
		return changeCryKey_sendEX(dst, apack, ChangeCryKey_Total, wlent, nil)
	}
*/

func GetNetLocalRDPort() string {
	return fmt.Sprintf("127.0.0.1:%d", gsrand.GetRD_netPortNumber())
}

func setTcpConnKeepAlive(conn net.Conn) {
	tcpconn1, ok := conn.(*net.TCPConn)
	if !ok {
		panic("!ok")
	}
	tcpconn1.SetKeepAlive(true)
	tcpconn1.SetKeepAlivePeriod(time.Second)
}

func GetRDNetConn() (net.Conn, net.Conn) {
	var conna, connc net.Conn

	serverAddr := ""

	wg := sync.WaitGroup{}
	listenDone := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()

		lst, err := net.Listen("tcp4", "127.0.0.1:")
		checkError_panic(err)
		serverAddr = lst.Addr().String()

		listenDone <- 1
		//defer lst.Close()
		conna, err = lst.Accept()
		checkError_exit(err)
	}()
	<-listenDone
	time.Sleep(time.Millisecond * 10)
	connc, err := net.Dial("tcp4", serverAddr)
	checkError_exit(err)
	wg.Wait()

	setTcpConnKeepAlive(conna)
	setTcpConnKeepAlive(connc)

	return gsnetconn.NewGSTNetConn(conna), gsnetconn.NewGSTNetConn(connc)
}

func GetRDNetConnLiten() (net.Conn, net.Conn, net.Listener) {
	var conna, connc net.Conn

	serverAddr := ""

	wg := sync.WaitGroup{}
	listenDone := make(chan int)
	var lst net.Listener
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()

		lst, err = net.Listen("tcp4", "127.0.0.1:")
		checkError_panic(err)
		serverAddr = lst.Addr().String()

		listenDone <- 1
		//defer lst.Close()
		conna, err = lst.Accept()
		checkError_exit(err)
	}()
	<-listenDone
	time.Sleep(time.Millisecond * 10)
	connc, err = net.Dial("tcp4", serverAddr)
	checkError_exit(err)
	wg.Wait()

	setTcpConnKeepAlive(conna)
	setTcpConnKeepAlive(connc)

	return gsnetconn.NewGSTNetConn(conna), gsnetconn.NewGSTNetConn(connc), lst
}

func IsNetConnClosed(conn net.Conn) bool {
	err := conn.SetDeadline(time.Now().Add(time.Second * 60))
	if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
		return true
	} else {
		//checkError_panic(err)
		return false
	}
}

func IsNetConnClosed_Read(conn net.Conn) bool {
	rbuf := make([]byte, 1)
	err := conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	_, err = conn.Read(rbuf)
	if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
		return true
	} else {
		//checkError_panic(err)
		return false
	}
}

func isClosed_Read(conn net.Conn) bool {

	rbuf := make([]byte, 1)
	err := conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
		return true
	}
	_, err = conn.Read(rbuf)
	if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
		return true
	} else {
		return false
	}
}

func IsClosedPanic(conn net.Conn) {
	if g_debug_netconn_isclosed {
		if isClosed_Read(conn) {
			panic("closed")
		}
	}
	// TcpDoKeepAlive(conn)
}

func TcpDoKeepAlive(conn net.Conn) {
	go conn.Write([]byte("1"))
}

func GetIP(conn net.Addr) string {
	s1 := conn.String()
	ix := strings.Index(s1, ":")
	if ix != -1 {
		return s1[:ix]
	}
	return s1
}

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type GSTNetConnDebug struct {
	conn net.Conn
}

func (c *GSTNetConnDebug) Read(b []byte) (n int, err error)   { return c.conn.Read(b) }
func (c *GSTNetConnDebug) Write(b []byte) (n int, err error)  { return c.conn.Write(b) }
func (c *GSTNetConnDebug) Close() error                       { return c.conn.Close() }
func (c *GSTNetConnDebug) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *GSTNetConnDebug) RemoteAddr() net.Addr               { return c.conn.RemoteAddr() }
func (c *GSTNetConnDebug) SetDeadline(t time.Time) error      { return c.conn.SetDeadline(t) }
func (c *GSTNetConnDebug) SetReadDeadline(t time.Time) error  { return c.conn.SetReadDeadline(t) }
func (c *GSTNetConnDebug) SetWriteDeadline(t time.Time) error { return c.conn.SetWriteDeadline(t) }

func NewGSTNetConnDebug(conn net.Conn) *GSTNetConnDebug { return &GSTNetConnDebug{conn} }
