package gstunnellib

import (
	"net"
	"sync/atomic"
)

type Gorou_status interface {
	IsOk() bool
	SetOk()
	SetClose()
}

type gorou_statusImp struct {
	status int32
}

const (
	gorou_s_begin = 0
	gorou_s_ok    = 1
	gorou_s_close = 2
	gorou_s_end   = 3
)

func (g *gorou_statusImp) IsOk() bool {
	return atomic.LoadInt32(&g.status) == gorou_s_ok
}
func (g *gorou_statusImp) SetOk()    { atomic.SwapInt32(&g.status, gorou_s_ok) }
func (g *gorou_statusImp) SetClose() { atomic.SwapInt32(&g.status, gorou_s_close) }

func NewGorouStatus() Gorou_status {
	g1 := new(gorou_statusImp)
	g1.SetOk()
	return g1
}

type gorou_status_netConnImp struct {
	gorou_statusImp
	connList []net.Conn
}

func (g *gorou_status_netConnImp) SetClose() {
	for _, item := range g.connList {
		item.Close()
	}
	g.gorou_statusImp.SetClose()
}

func NewGorouStatusNetConn(inconnList []net.Conn) Gorou_status {
	g1 := new(gorou_status_netConnImp)
	g1.connList = inconnList
	g1.SetOk()
	return g1
}
