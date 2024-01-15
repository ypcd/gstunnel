package gstunnellib

type IGSContext interface {
	GetGsId() uint64
	GetGsStatus() IGsStatus
	Close()
}

type gsContextImp struct {
	gsid uint64
	gsSt IGsStatus
}

func NewGSContextImp(gsid uint64, gsst IGsStatus) IGSContext {
	return &gsContextImp{gsid, gsst}
}

func (gc *gsContextImp) GetGsId() uint64 { return gc.gsid }

func (gc *gsContextImp) GetGsStatus() IGsStatus { return gc.gsSt }

func (gc *gsContextImp) Close() {
	gc.gsSt.GetStatusConnList().Delete(gc.gsid)
}
