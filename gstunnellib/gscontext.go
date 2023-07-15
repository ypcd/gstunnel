package gstunnellib

type GsContext interface {
	GetGsId() uint64
	GetGsStatus() GsStatus
	Close()
}

type gsContextImp struct {
	gsid uint64
	gsSt GsStatus
}

func NewGsContextImp(gsid uint64, gsst GsStatus) GsContext {
	return &gsContextImp{gsid, gsst}
}

func (gc *gsContextImp) GetGsId() uint64 { return gc.gsid }

func (gc *gsContextImp) GetGsStatus() GsStatus { return gc.gsSt }

func (gc *gsContextImp) Close() {
	gc.gsSt.GetStatusConnList().Delete(gc.gsid)
}
