package gstestpipe

import (
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib/gstestpipe/gsmap"
)

// no echo
type RawServerSocket struct {
	listen_close    atomic.Bool
	connList        chan net.Conn
	close           atomic.Bool
	listening       bool
	serverAddr      string
	connListIn      []net.Conn
	lock_connListIn sync.Mutex
	listener        net.Listener
	lock_listener   sync.Mutex

	count_Accept atomic.Int32
}

func NewRawServerSocket_RandAddr() *RawServerSocket {
	return &RawServerSocket{connList: make(chan net.Conn, 10240),
		serverAddr: GetRandAddr()}
}

func (ss *RawServerSocket) listen() {
	if ss.listening {
		return
	}
	wgrun := sync.WaitGroup{}
	wgrun.Add(1)
	go func() {
		defer ss.listen_close.Swap(true)
		defer close(ss.connList)

		lst, err := net.Listen("tcp4", ss.serverAddr)
		checkError_exit(err)
		defer lst.Close()
		ss.lock_listener.Lock()
		ss.listener = lst
		ss.lock_listener.Unlock()
		wg := sync.WaitGroup{}
		for i := 0; i < 32; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
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
			ss.count_Accept.Add(1)
		}
		wgrun.Done()
		wg.Wait()
	}()
	ss.listening = true
	wgrun.Wait()
}

func (ss *RawServerSocket) Run() { ss.listen() }

func (ss *RawServerSocket) Close() {
	ss.close.Swap(true)
	ss.lock_listener.Lock()
	ss.listener.Close()
	ss.lock_listener.Unlock()

	for conn := range ss.connList {
		conn.Close()
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
	echo_running      int
	lock_echo_running sync.Mutex

	count_echo_go         atomic.Int32
	count_list_echo_bytes *gsmap.Gsmap
}

func NewRawServerSocketEcho_RandAddr() *RawServerSocketEcho {
	return &RawServerSocketEcho{RawServerSocket: RawServerSocket{connList: make(chan net.Conn, 10240),
		serverAddr: GetRandAddr()},
		count_list_echo_bytes: gsmap.NewGSMap(),
	}
}

func (ss *RawServerSocketEcho) echoHandler() {
	ss.lock_echo_running.Lock()
	defer ss.lock_echo_running.Unlock()
	if ss.echo_running == 1 {
		return
	}

	go func() {
		if !ss.close.Load() {
			for conn1 := range ss.connList {
				go func(conn net.Conn) {
					cpSize, err := io.Copy(conn, conn)
					ss.count_list_echo_bytes.Add(conn.RemoteAddr().String(), strconv.Itoa(int(cpSize)))
					if errors.Is(err, net.ErrClosed) {
						checkError_info(err)
						return
					}
					checkError_exit(err)

				}(conn1)
				ss.count_echo_go.Add(1)
			}
		}
	}()
	ss.echo_running++
}

func (ss *RawServerSocketEcho) Run() {
	ss.listen()
	ss.echoHandler()
}
