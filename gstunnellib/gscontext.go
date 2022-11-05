package gstunnellib

type GsContext interface {
	GetGsId() uint64
}

type gsContextImp struct {
	gsid uint64
}

func NewGsContextImp(gsid uint64) GsContext {
	return &gsContextImp{gsid}
}

func (gc *gsContextImp) GetGsId() uint64 {
	return gc.gsid
}
