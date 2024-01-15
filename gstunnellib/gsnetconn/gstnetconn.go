package gsnetconn

import (
	"net"
	"sync/atomic"
	"time"
)

type GSTNetConn struct {
	conn                                         net.Conn
	rlent, wlent                                 int
	ntTimeOut                                    time.Duration
	setTimeOut, setTimeOutRead, setTimeOutWirite float64
	closed                                       atomic.Bool
	localAddr, remoteAddr                        string
}

// blocking
func NewGSTNetConn(conn net.Conn) *GSTNetConn {
	return &GSTNetConn{
		conn:       conn,
		ntTimeOut:  0,
		localAddr:  conn.LocalAddr().String(),
		remoteAddr: conn.RemoteAddr().String(),
	}

}

// non-blocking
func NewGSTNetConnNonBlcoking(conn net.Conn, timeout time.Duration) *GSTNetConn {
	return &GSTNetConn{conn: conn, ntTimeOut: timeout}
}

func (c *GSTNetConn) Read(buf []byte) (n int, err error) {
	var err2 error
	if c.ntTimeOut != 0 {
		err2 = c.conn.SetReadDeadline(time.Now().Add(c.ntTimeOut))
	} else {
		//err2 = c.conn.SetReadDeadline(time.Time{})
	}

	if err != nil {
		return 0, err2
	}
	rn, err := c.conn.Read(buf)
	c.rlent += rn
	return rn, err
}

func (c *GSTNetConn) Write(buf []byte) (n int, err error) {
	var err2 error
	if c.ntTimeOut != 0 {
		err2 = c.conn.SetReadDeadline(time.Now().Add(c.ntTimeOut))
	} else {
		//err2 = c.conn.SetReadDeadline(time.Time{})
	}

	if err != nil {
		return 0, err2
	}
	rn, err := c.conn.Write(buf)
	c.wlent += rn
	return rn, err
}

func (c *GSTNetConn) Close() error {
	c.closed.Store(true)
	return c.conn.Close()
}
func (c *GSTNetConn) RemoteAddr() net.Addr { return c.conn.RemoteAddr() }
func (c *GSTNetConn) LocalAddr() net.Addr  { return c.conn.LocalAddr() }

func (c *GSTNetConn) SetDeadline(t time.Time) error {
	c.setTimeOut = t.Sub(time.Now()).Seconds()
	return c.conn.SetDeadline(t)
}

func (c *GSTNetConn) SetReadDeadline(t time.Time) error {
	c.setTimeOutRead = t.Sub(time.Now()).Seconds()
	return c.conn.SetReadDeadline(t)
}
func (c *GSTNetConn) SetWriteDeadline(t time.Time) error {
	c.setTimeOutWirite = t.Sub(time.Now()).Seconds()
	return c.conn.SetWriteDeadline(t)
}

/*
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
	LocalAddr() Addr
	RemoteAddr() Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error

*/
