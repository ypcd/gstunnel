package main

import (
	"errors"
	"io"
	"net"
	"os"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_st(obj *gsobj.GstObj) {
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, obj.Gctx)

	defer obj.Close()

	//tmr_out := timerm.CreateTimer(obj.NetworkTimeout)
	//tmrP := timerm.CreateTimer(obj.Tmr_display_time)
	tmrP2 := timerm.CreateTimer(obj.Tmr_display_time)

	tmr_changekey := timerm.CreateTimer(g_tmr_changekey_time)

	var timer1, timew1 time.Time

	var err error
	_ = err

	defer func() {
		g_RuntimeStatistics.AddServerTotalNetData_recv(int(obj.Rlent))
		g_RuntimeStatistics.AddSrcTotalNetData_send(int(obj.Wlent))
		g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\t%s\tpack  trlen:%d  twlen:%d  ChangeCryKey_total:%d\n\t%s\t%s",
			obj.Gctx.GetGsId(),
			gstunnellib.GetNetConnAddrString("obj.Src", obj.Src),
			gstunnellib.GetNetConnAddrString("obj.Dst", obj.Dst),
			obj.Rlent, obj.Wlent, obj.ChangeCryKey_Total,
			obj.Nt_read.PrintString(),
			obj.Nt_write.PrintString(),
		)

		if g_Values.GetDebug() {
		}
	}()

	err = IsTheVersionConsistent_send(obj.Dst, obj.Apack, &obj.Wlent)
	checkError_panic_GsCtx(err, obj.Gctx)

	err = ChangeCryKey_send(obj.Dst, obj.Apack, &obj.ChangeCryKey_Total, &obj.Wlent)
	checkError_panic_GsCtx(err, obj.Gctx)

	for {
		//	recot_p_r.Run()
		obj.Src.SetReadDeadline(time.Now().Add(obj.NetworkTimeout))
		timer1 = time.Now()
		rlen, err := obj.Src.Read(obj.Rbuf)
		obj.Rlent += int64(rlen)
		//	recot_p_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info_GsCtx(err, obj.Gctx)
			return
		} else {
			checkError_panic_GsCtx(err, obj.Gctx)
		}
		obj.Nt_read.Add(time.Since(timer1))
		/*
			if tmr_out.Run() {
				g_Logger.Printf("Error: [%d] Time out, func exit.\n", obj.Gctx.GetGsId())
				return
			}
		*/
		if rlen == 0 {
			g_Logger.Println("Error: obj.Src.read() rlen==0 func exit.")
			return
		}

		if rlen > 0 {
			if gstunnellib.G_RunTime_Debug {
				gstunnellib.G_RunTimeDebugInfo1.AddPackingPackSizeList("server_srcToDstP_st_packing len", rlen)
			}
			obj.Wbuf = obj.Apack.Packing(obj.Rbuf[:rlen])

			if len(obj.Wbuf) <= 0 {
				g_Logger.Println("Error: gspack.packing is error.")
				return
			}
			obj.Dst.SetWriteDeadline(time.Now().Add(obj.NetworkTimeout))
			timew1 = time.Now()
			//wlen, err := gstunnellib.NetConnWriteAll(obj.Dst, obj.Wbuf)
			wlen, err := gstunnellib.NetConnWriteAll(obj.Dst, obj.Wbuf)
			obj.Wlent += int64(wlen)
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				checkError_info_GsCtx(err, obj.Gctx)
				return
			} else {
				checkError_panic_GsCtx(err, obj.Gctx)
			}
			obj.Nt_write.Add(time.Since(timew1))
			//tmr_out.Boot()
		}
		if tmr_changekey.Run() {
			err = ChangeCryKey_send(obj.Dst, obj.Apack, &obj.ChangeCryKey_Total, &obj.Wlent)
			if err != nil {
				g_Logger.Println("Error:", err)
				return
			}
		}
		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("pack  trlen:%d  twlen:%d\n", obj.Rlent, obj.Wlent)
			//g_Logger.Println("goPackTotal:", goPackTotal)
			g_Logger.Println("ChangeCryKey_total:", obj.ChangeCryKey_Total)

			//	g_Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			//	g_Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}

	}
}
