/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"errors"
	"io"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_mt(obj *gsobj.GstObj) {
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, obj.Gctx)
	defer obj.Close()

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	tmr_changekey := timerm.CreateTimer(g_tmr_changekey_time)

	nt_read := gstunnellib.NewNetTimeImpName("read")
	var timer1 time.Time

	//recot_p_r := timerm.CreateRecoTime()

	//fp1, fp2 = g_fpnull, g_fpnull
	//_, _ = fp1, fp2

	var err error
	_ = err

	objw := obj.NewGstObjW(g_netPUn_chan_cache_size)

	defer func() {
		g_RuntimeStatistics.AddSrcTotalNetData_recv(int(obj.Rlent))
		g_log_List.GSNetIOLen.Printf(
			"[%d] gorou exit.\n\t%s\tpack  trlen:%d  ChangeCryKey_total:%d\n\t%s",
			obj.Gctx.GetGsId(), gstunnellib.GetNetConnAddrString("obj.Src", obj.Src), obj.Rlent, obj.ChangeCryKey_Total,
			nt_read.PrintString(),
		)

		if g_Values.GetDebug() {
			//g_Logger.Println("\tgoPackTotal:", atomic.LoadInt32(&g_goPackTotal))

			//g_Logger.Println("\tRecoTime_p_r All: ", recot_p_r.StringAll())
		}
	}()

	err = IsTheVersionConsistent_send(obj.Dst, obj.Apack, &obj.Wlent)
	if err != nil {
		g_Logger.Println("Error:", err, " func exit.")
		return
	}
	err = ChangeCryKey_send(obj.Dst, obj.Apack, &obj.ChangeCryKey_Total, &obj.Wlent)
	if err != nil {
		g_Logger.Println("Error:", err, " func exit.")
		return
	}

	objw.Wg_w.Add(1)
	go srcTOdstP_w(objw)

	for objw.Dst_ok.IsOk() {

		//recot_p_r.Run()
		obj.Src.SetReadDeadline(time.Now().Add(g_networkTimeout))
		timer1 = time.Now()
		rlen, err := obj.Src.Read(obj.Rbuf)
		obj.Rlent += int64(rlen)
		//recot_p_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info_GsCtx(err, obj.Gctx)
			return
		} else {
			checkError_panic_GsCtx(err, obj.Gctx)
		}
		nt_read.Add(time.Since(timer1))
		/*
			if tmr_out.Run() {
				g_Logger.Println("Error: time out func exit.")
				return
			}
		*/
		if rlen <= 0 {
			g_Logger.Println("Error: obj.Src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			g_Logger.Println("Error:", err)
			continue
		}

		//tmr_out.Boot()

		if rlen > 0 {
			obj.Wbuf = obj.Apack.Packing(obj.Rbuf[:rlen])
		} else {
			continue
		}

		objw.Dst_chan <- obj.Wbuf
		obj.Wbuf = nil

		if !objw.Dst_ok.IsOk() {
			g_Logger.Printf("Error: [%d] not objw.Dst_ok.isok() func exit.\n", obj.Gctx.GetGsId())
			return
		}

		if tmr_changekey.Run() {
			obj.Wbuf = obj.Apack.ChangeCryKey()
			obj.ChangeCryKey_Total += 1
			objw.Dst_chan <- obj.Wbuf
			obj.Wbuf = nil
			//tmr_out.Boot()
		}
		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("pack  trlen:%d\n", obj.Rlent)
			g_Logger.Println("g_goPackTotal:", atomic.LoadInt32(&g_goPackTotal))
			g_Logger.Println("ChangeCryKey_total:", obj.ChangeCryKey_Total)

			//g_Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())

		}

	}
	g_Logger.Println("Func exit.")
}
