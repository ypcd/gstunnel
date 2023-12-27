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

func srcTOdstUn_st(obj *gsobj.GstObj) {
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, obj.Gctx)

	defer obj.Close()

	//tmr_out := timerm.CreateTimer(obj.NetworkTimeout)
	//tmrP := timerm.CreateTimer(obj.Tmr_display_time)
	tmrP2 := timerm.CreateTimer(obj.Tmr_display_time)

	var timer1, timew1 time.Time

	//	recot_un_r := timerm.CreateRecoTime()
	//	recot_un_w := timerm.CreateRecoTime()

	//	fp1 := "SUrecv.data"
	//	fp2 := "SUsend.data"
	//	fp1 = g_fpnull
	//	fp2 = g_fpnull
	//	_, _ = fp1, fp2

	//var err error
	//outf, err := os.Create(fp1)
	//checkError_GsCtx(err,obj.Gctx)
	//outf2, err := os.Create(fp2)
	//checkError_GsCtx(err,obj.Gctx)

	//defer outf.Close()
	//defer outf2.Close()

	defer func() {
		g_RuntimeStatistics.AddSrcTotalNetData_recv(int(obj.Rlent))
		g_RuntimeStatistics.AddServerTotalNetData_send(int(obj.Wlent))
		g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\t%s\tunpack trlen:%d  twlen:%d\n\t%s\t%s",
			obj.Gctx.GetGsId(),
			gstunnellib.GetNetConnAddrString("obj.Src", obj.Src),
			gstunnellib.GetNetConnAddrString("obj.Dst", obj.Dst),
			obj.Rlent, obj.Wlent,
			obj.Nt_read.PrintString(),
			obj.Nt_write.PrintString(),
		)

		if g_Values.GetDebug() {

			//	g_Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	g_Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}()

	for {
		//	recot_un_r.Run()
		obj.Src.SetReadDeadline(time.Now().Add(obj.NetworkTimeout))
		timer1 = time.Now()
		rlen, err := obj.Src.Read(obj.Rbuf)
		obj.Rlent += int64(rlen)
		//	recot_un_r.Run()
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
		if err != nil {
			g_Logger.Println("Error:", err)
			continue
		}
		//tmr_out.Boot()

		obj.Apack.WriteEncryData(obj.Rbuf[:rlen])
		obj.Wbuf, err = obj.Apack.GetDecryData()
		checkError_panic_GsCtx(err, obj.Gctx)
		if len(obj.Wbuf) > 0 {
			obj.Dst.SetWriteDeadline(time.Now().Add(obj.NetworkTimeout))
			timew1 = time.Now()
			//rn, err := gstunnellib.NetConnWriteAll(obj.Dst, obj.Wbuf)
			rn, err := gstunnellib.NetConnWriteAll(obj.Dst, obj.Wbuf)
			obj.Wlent += int64(rn)
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

		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("unpack trlen:%d  twlen:%d\n", obj.Rlent, obj.Wlent)

			//	g_Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	g_Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}
}
