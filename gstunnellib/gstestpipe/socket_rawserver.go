package gstestpipe

import (
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// no echo
type RawServerSocket struct {
	connList        chan net.Conn
	close           atomic.Bool
	listening       bool
	serverAddr      string
	connListIn      []net.Conn
	lock_connListIn sync.Mutex
	listen_close    atomic.Bool
	listener        net.Listener
}

func NewRawServerSocket_RandAddr() *RawServerSocket {
	return &RawServerSocket{connList: make(chan net.Conn, 1024),
		serverAddr: GetRandAddr()}
}

func (ss *RawServerSocket) listen() {
	if ss.listening {
		return
	}
	ss.listening = true

	go func() {
		defer ss.listen_close.Swap(true)

		lst, err := net.Listen("tcp4", ss.serverAddr)
		checkError_exit(err)
		defer lst.Close()
		ss.listener = lst
		for !ss.close.Load() {
			conna, err := lst.Accept()
			if errors.Is(err, net.ErrClosed) {
				return
			}
			checkError_exit(err)
			ss.lock_connListIn.Lock()
			ss.connListIn = append(ss.connListIn, conna)
			ss.lock_connListIn.Unlock()
			ss.connList <- conna
		}
	}()

}

func (ss *RawServerSocket) Run() { ss.listen() }

func (ss *RawServerSocket) Close() {
	ss.close.Swap(true)
	ss.listener.Close()
	if ss.listen_close.Load() {
		close(ss.connList)
	}
	for {
		if ss.listen_close.Load() {
			ss.lock_connListIn.Lock()
			defer ss.lock_connListIn.Unlock()
			for _, v := range ss.connListIn {
				v.Close()
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func (ss *RawServerSocket) GetConnList() chan net.Conn { return ss.connList }
func (ss *RawServerSocket) GetServerAddr() string      { return ss.serverAddr }

type RawServerSocketEcho struct {
	RawServerSocket
	echo_running bool
}

func NewRawServerSocketEcho_RandAddr() *RawServerSocketEcho {
	return &RawServerSocketEcho{RawServerSocket: RawServerSocket{connList: make(chan net.Conn, 1024),
		serverAddr: GetRandAddr()}}
}

func (ss *RawServerSocketEcho) echoHandler() {
	if ss.echo_running {
		return
	}

	go func() {
		for !ss.close.Load() {
			for conn1 := range ss.connList {
				go func(conn net.Conn) {
					for {
						_, err := io.Copy(conn, conn)
						if errors.Is(err, net.ErrClosed) {
							break
						}
						checkError_exit(err)
					}
				}(conn1)
			}
		}
	}()
	ss.echo_running = true
}

func (ss *RawServerSocketEcho) Run() {
	ss.listen()
	ss.echoHandler()
}
