package gstunnellib

import "sync/atomic"

type GId interface {
	GetId() uint64
}

type gidImp struct {
	id uint64
}

func NewGIdImp() GId {
	return &gidImp{0}
}

func (g *gidImp) GetId() uint64 {
	return atomic.AddUint64(&g.id, 1)
}
