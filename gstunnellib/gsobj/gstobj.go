package gsobj

import (
	"net"
	"sync"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
)

type GstObj struct {
	Rlent int64
	Wlent int64

	Rbuf []byte
	Wbuf []byte

	ChangeCryKey_Total int

	Apack gstunnellib.GsPackNet

	Nt_read  gstunnellib.NetTime
	Nt_write gstunnellib.NetTime

	Src  net.Conn
	Dst  net.Conn
	Gctx gstunnellib.GsContext

	Key              string
	Tmr_display_time time.Duration
	NetworkTimeout   time.Duration
	objw             *GstObjW
	//Net_read_size    int

	//	apack := gstunnellib.NewGsPack(obj.Key)

	//	var wbuf []byte
	//var rbuf []byte = make([]byte, obj.Net_read_size)
	//

	//还没有完成写入的缓存数据

	//g_Values
	//g_gstst
	//g_tmr_changekey_time
}

func NewGstObj(src, dst net.Conn, gctx gstunnellib.GsContext, tmr_display_time, networkTimeout time.Duration, key string, net_read_size int) *GstObj {
	return &GstObj{
		Src:              src,
		Dst:              dst,
		Gctx:             gctx,
		Tmr_display_time: tmr_display_time,
		NetworkTimeout:   networkTimeout,
		Key:              key,
		//Net_read_size:    net_read_size,
		Nt_read:  gstunnellib.NewNetTimeImpName("read"),
		Nt_write: gstunnellib.NewNetTimeImpName("write"),
		Rbuf:     make([]byte, net_read_size),
		Apack:    gstunnellib.NewGsPackNet(key),
	}
}

func (o *GstObj) Close() {
	if o.objw != nil {
		gstunnellib.ChanClose(o.objw.Dst_chan)
		o.objw.Wg_w.Wait()
	}
	o.Src.Close()
	o.Dst.Close()
}

type GstObjW struct {
	Dst              net.Conn
	Wbuf             []byte
	Dst_chan         chan []byte
	Dst_ok           gstunnellib.Gorou_status
	Wlent            int64
	Wg_w             *sync.WaitGroup
	Gctx             gstunnellib.GsContext
	Nt_write         gstunnellib.NetTime
	Tmr_display_time time.Duration
	NetworkTimeout   time.Duration
	Key              string
	Src              net.Conn
	Apack            gstunnellib.GsPackNet
}

func (o *GstObj) NewGstObjW(netPUn_chan_cache_size int) *GstObjW {
	o.objw = &GstObjW{
		Src:              o.Src,
		Dst:              o.Dst,
		Dst_chan:         make(chan []byte, netPUn_chan_cache_size),
		Dst_ok:           gstunnellib.NewGorouStatusNetConn([]net.Conn{o.Src, o.Dst}),
		Wlent:            o.Wlent,
		Wg_w:             &sync.WaitGroup{},
		Gctx:             o.Gctx,
		Nt_write:         gstunnellib.NewNetTimeImpName("write"),
		Tmr_display_time: o.Tmr_display_time,
		Key:              o.Key,
		NetworkTimeout:   o.NetworkTimeout,
		Apack:            gstunnellib.NewGsPackNet(o.Key),
	}
	return o.objw
}

func (o *GstObjW) Close() {
	o.Dst_ok.SetClose()
	gstunnellib.ChanClean(o.Dst_chan)

	o.Src.Close()
	o.Dst.Close()
}
