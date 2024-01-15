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

const (
	GOServer uint8 = 1<<8 - 1     //255
	GOClient uint8 = 1<<8 - 1 - 3 //252
)

type GSTObj struct {
	Rlent int64
	Wlent int64

	Rbuf []byte
	Wbuf []byte

	ChangeCryKey_Total int

	apack gstunnellib.IGSRSAPackNet

	Nt_read  gstunnellib.INetTime
	Nt_write gstunnellib.INetTime

	src  net.Conn
	dst  net.Conn
	Gctx gstunnellib.IGSContext

	Key string
	//serverRSAKey     *gsrsa.RSA
	//clientRSAKey     *gsrsa.RSA
	Tmr_display_time time.Duration
	NetworkTimeout   time.Duration
	objw             *GSTObjW

	gstType uint8
	//Net_read_size    int

	//	apack := gstunnellib.NewGsPack(obj.Key)

	//	var wbuf []byte
	//var rbuf []byte = make([]byte, obj.Net_read_size)
	//

	//还没有完成写入的缓存数据

	//g_Values
	//g_gstst
	//g_tmr_Changekey_time
}

func newGstObj(gstType uint8, src, dst net.Conn, gctx gstunnellib.IGSContext, tmr_display_time, networkTimeout time.Duration, key string, net_read_size int, apack gstunnellib.IGSRSAPackNet) *GSTObj {

	return &GSTObj{
		gstType:          gstType,
		src:              src,
		dst:              dst,
		Gctx:             gctx,
		Tmr_display_time: tmr_display_time,
		NetworkTimeout:   networkTimeout,
		Key:              key,
		//Net_read_size:    net_read_size,
		Nt_read:  gstunnellib.NewNetTimeImpName("read"),
		Nt_write: gstunnellib.NewNetTimeImpName("write"),
		Rbuf:     make([]byte, net_read_size),
		apack:    apack,
		//serverRSAKey: serverRSAKey,
		//clientRSAKey: clientRSAKey,
	}

}
func NewGstObjWithServer(src, dst net.Conn, gctx gstunnellib.IGSContext, tmr_display_time, networkTimeout time.Duration, key string, net_read_size int, serverRSAKey *gsrsa.RSA) *GSTObj {
	return newGstObj(
		GOServer,
		src,
		dst,
		gctx,
		tmr_display_time,
		networkTimeout,
		key,
		net_read_size,
		//serverRSAKey.NewRSA(),
		//nil,
		gstunnellib.NewGSRSAPackNetImpWithGSTServer(key, serverRSAKey),
	)
}

func NewGstObjWithClient(src, dst net.Conn, gctx gstunnellib.IGSContext, tmr_display_time, networkTimeout time.Duration, key string, net_read_size int, serverRSAKey, clientRSAKey *gsrsa.RSA) *GSTObj {
	return newGstObj(
		GOClient,
		src,
		dst,
		gctx,
		tmr_display_time,
		networkTimeout,
		key,
		net_read_size,
		//serverRSAKey.NewRSA(),
		//clientRSAKey.NewRSA(),
		gstunnellib.NewGSRSAPackNetImpWithGSTClient(key, serverRSAKey, clientRSAKey),
	)
}

func (o *GSTObj) Close() {
	o.src.Close()
	o.dst.Close()

	if o.objw != nil {
		//gstunnellib.ChanClose(o.objw.Dst_chan)
		close(o.objw.Dst_chan)
		o.objw.Wg_w.Wait()
	}
}

func (o *GSTObj) NewGstObjW(netPUn_chan_cache_size int) *GSTObjW {

	objw := &GSTObjW{
		src:              o.src,
		dst:              o.dst,
		Dst_chan:         make(chan []byte, netPUn_chan_cache_size),
		Dst_ok:           gstunnellib.NewGorouStatusNetConn([]net.Conn{o.src, o.dst}),
		Wlent:            o.Wlent,
		Wg_w:             &sync.WaitGroup{},
		Gctx:             o.Gctx,
		Nt_write:         gstunnellib.NewNetTimeImpName("write"),
		Tmr_display_time: o.Tmr_display_time,
		//Key:              o.Key,
		NetworkTimeout: o.NetworkTimeout,
		apack:          gstunnellib.NewGSRSAPackNetImp(o.Key, o.apack.GetServerRSAKey(), o.GetClientRSAKey()),
	}
	//o.objw.apack.SetClientRSAKey(o.clientRSAKey)
	o.objw = objw
	return objw
}

func (o *GSTObj) SetClientRSAKey(ckey *gsrsa.RSA) {
	//o.clientRSAKey = ckey.NewRSA()
	o.apack.SetClientRSAKey(ckey)
}

func (o *GSTObj) WriteEncryData(data []byte) error {
	return o.apack.WriteEncryData(data)
}

func (o *GSTObj) GetDecryData() ([]byte, error) {
	return o.apack.GetDecryData()
}

func (o *GSTObj) ClientPublicKeyPack() []byte {
	return o.apack.ClientPublicKeyPack()
}

func (o *GSTObj) ChangeCryKeyFromGSTServer() (packData []byte, outkey string) {
	return o.apack.ChangeCryKeyFromGSTServer()
}

func (o *GSTObj) ChangeCryKeyFromGSTClient() (packData []byte, outkey string) {
	return o.apack.ChangeCryKeyFromGSTClient()
}

func (o *GSTObj) IsExistsClientKey() bool {
	return o.apack.IsExistsClientKey()
}

func (o *GSTObj) GetClientRSAKey() *gsrsa.RSA {
	return o.apack.GetClientRSAKey()
}

func (o *GSTObj) PackVersion() []byte {
	return o.apack.PackVersion()
}

func (o *GSTObj) Packing(data []byte) []byte {
	return o.apack.Packing(data)
}

func (o *GSTObj) VersionPack_send() error {
	return VersionPack_send(o)
}

func (o *GSTObj) ChangeCryKey_send() error {
	switch o.gstType {
	case GOServer:
		return changeCryKey_send_fromServer(o)
	case GOClient:
		return changeCryKey_send_fromClient(o)
	}
	panic("obj.gstType is error")
}

func (o *GSTObj) ReadNetSrc(buf []byte) (n int, err error) {
	err2 := o.src.SetReadDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.src.Read(buf)
}
func (o *GSTObj) WriteNetSrc(buf []byte) (n int, err error) {
	err2 := o.src.SetWriteDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.src.Write(buf)
}

func (o *GSTObj) ReadNetDst(buf []byte) (n int, err error) {
	err2 := o.dst.SetReadDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.dst.Read(buf)
}
func (o *GSTObj) WriteNetDst(buf []byte) (n int, err error) {
	err2 := o.dst.SetWriteDeadline(time.Now().Add(o.NetworkTimeout))
	if err != nil {
		return 0, err2
	}
	return o.dst.Write(buf)
}

func (o *GSTObj) NetConnWriteAll(buf []byte) (int64, error) {
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

func (o *GSTObj) GetNetConnAddrString(str string, conn net.Conn) string {
	return fmt.Sprintf("%s: [localIp:%s  remoteIp:%s]\n",
		str, conn.LocalAddr().String(), conn.RemoteAddr().String())
}

func (obj *GSTObj) StringWithGOExit() string {
	return fmt.Sprintf("[%d] gorou exit.\n\t%s\t%s\tunpack trlen:%d  twlen:%d\n\t%s\t%s",
		obj.Gctx.GetGsId(),
		gstunnellib.GetNetConnAddrString("obj.Src", obj.src),
		gstunnellib.GetNetConnAddrString("obj.Dst", obj.dst),
		obj.Rlent, obj.Wlent,
		obj.Nt_read.PrintString(),
		obj.Nt_write.PrintString(),
	)
}
