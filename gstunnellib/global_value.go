package gstunnellib

import "sync/atomic"

type Global_values interface {
	GetDebug() bool
	SetDebug(bool)
}

type global_valuesImp struct {
	debug uint32
}

func NewGlobalValuesImp() Global_values {
	g1 := global_valuesImp{}
	g1.init()
	return &g1
}

func (g *global_valuesImp) init() {
	atomic.SwapUint32(&g.debug, 0)
}

func (g *global_valuesImp) GetDebug() bool {
	d1 := atomic.LoadUint32(&g.debug)
	return d1 == 1
}
func (g *global_valuesImp) SetDebug(bl bool) {
	if bl {
		atomic.SwapUint32(&g.debug, 1)
	} else {
		atomic.SwapUint32(&g.debug, 0)
	}
}
