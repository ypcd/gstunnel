package gsobj

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gserror"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsrsa"
)

type GSTObjW struct {
	dst              net.Conn
	Wbuf             []byte
	Dst_chan         chan []byte
	Dst_ok           gstunnellib.IGorou_status
	Wlent            int64
	Wg_w             *sync.WaitGroup
	Gctx             gstunnellib.IGSContext
	Nt_write         gstunnellib.INetTime
	Tmr_display_time time.Duration
	NetworkTimeout   time.Duration
	//Key              string
	src                net.Conn
	apack              gstunnellib.IGSRSAPackNet
	ChangeCryKey_Total int
}

func (o *GSTObjW) Close() {
	o.src.Close()
	o.dst.Close()

	o.Dst_ok.SetClose()
	gstunnellib.ChanClean(o.Dst_chan)

}

func (o *GSTObjW) WriteEncryData(data []byte) error {
	return o.apack.WriteEncryData(data)
}

func (o *GSTObjW) GetDecryData() ([]byte, error) {
	return o.apack.GetDecryData()
}

func (o *GSTObjW) Packing(data []byte) []byte {
	return o.apack.Packing(data)
}

func (o *GSTObjW) ClientPublicKeyPack() []byte {
	return o.apack.ClientPublicKeyPack()
}

/*
	func (o *GSTObjW) ChangeCryKeyFromGSTServer() []byte {
		return o.apack.ChangeCryKeyFromGSTServer()
	}

	func (o *GSTObjW) ChangeCryKeyFromGSTClient() []byte {
		return o.apack.ChangeCryKeyFromGSTClient()
	}
*/
func (o *GSTObjW) ChangeCryKeyFromGSTServer() (packData []byte, outkey string) {
	return o.apack.ChangeCryKeyFromGSTServer()
}

func (o *GSTObjW) ChangeCryKeyFromGSTClient() (packData []byte, outkey string) {
	return o.apack.ChangeCryKeyFromGSTClient()
}

func (o *GSTObjW) IsExistsClientKey() bool {
	return o.apack.IsExistsClientKey()
}

func (o *GSTObjW) GetClientRSAKey() *gsrsa.RSA {
	return o.apack.GetClientRSAKey()
}

func (o *GSTObjW) PackVersion() []byte {
	return o.apack.PackVersion()
}

func (o *GSTObjW) VersionPack_send() error {
	return VersionPack_sendEx(o.dst, o, &o.Wlent, nil)
}

func (o *GSTObjW) ChangeCryKey_send() error {
	return changeCryKey_sendEX_fromServer(o.dst, o, &o.ChangeCryKey_Total, &o.Wlent, nil)
}

func (o *GSTObjW) ReadNetSrc(buf []byte) (n int, err error) {
	err2 := o.src.SetReadDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.src.Read(buf)
}
func (o *GSTObjW) WriteNetSrc(buf []byte) (n int, err error) {
	err2 := o.src.SetWriteDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.src.Write(buf)
}

func (o *GSTObjW) ReadNetDst(buf []byte) (n int, err error) {
	err2 := o.dst.SetReadDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.dst.Read(buf)
}
func (o *GSTObjW) WriteNetDst(buf []byte) (n int, err error) {
	err2 := o.dst.SetWriteDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.dst.Write(buf)
}

func (o *GSTObjW) NetConnWriteAll(buf []byte) (int64, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	var wlen int = 0
	for {
		wsz, err := o.dst.Write(buf)
		wlen += wsz
		if gserror.IsErrorNetUsually(err) {
			gserror.CheckError_info(err)
			return int64(wlen), nil
		} else {
			gserror.CheckError_panic(err)
		}
		if wlen == len(buf) {
			return int64(wlen), err
		} else if wsz > len(buf) {
			panic("error wlen>len(buf)")
		}
		buf = buf[wsz:]
	}
}

func (obj *GSTObjW) StringWithGOExit() string {
	return fmt.Sprintf("[%d] gorou exit.\n\t%s\t%s\tunpack  twlen:%d\n\t%s",
		obj.Gctx.GetGsId(),
		gstunnellib.GetNetConnAddrString("obj.Src", obj.src),
		gstunnellib.GetNetConnAddrString("obj.Dst", obj.dst),
		//obj.Rlent,
		obj.Wlent,
		//obj.Nt_read.PrintString(),
		obj.Nt_write.PrintString(),
	)
}
