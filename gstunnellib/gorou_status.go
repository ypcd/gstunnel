package gstunnellib

import "sync/atomic"

type Gorou_status struct {
	status int32
}

const (
	gorou_s_begin = 0
	gorou_s_ok    = 1
	gorou_s_close = 2
	gorou_s_end   = 3
)

func (g *Gorou_status) IsOk() bool {
	v1 := atomic.LoadInt32(&g.status)
	return v1 == gorou_s_ok
}
func (g *Gorou_status) SetOk()    { atomic.SwapInt32(&g.status, gorou_s_ok) }
func (g *Gorou_status) SetClose() { atomic.SwapInt32(&g.status, gorou_s_close) }

func CreateGorouStatus() *Gorou_status {
	g1 := new(Gorou_status)
	g1.status = gorou_s_ok
	return g1
}
