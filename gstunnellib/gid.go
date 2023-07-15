package gstunnellib

import "sync/atomic"

type GId interface {
	GenerateId() uint64
	GetId() uint64
}

type gidImp struct {
	Id uint64
}

func NewGIdImp() GId {
	return &gidImp{0}
}

func (g *gidImp) GenerateId() uint64 {
	return atomic.AddUint64(&g.Id, 1)
}

func (g *gidImp) GetId() uint64 {
	return atomic.LoadUint64(&g.Id)
}
