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

func srcTOdstUn_w(obj *gsobj.GstObjW) {
	defer obj.Wg_w.Done()
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, obj.Gctx)
	defer obj.Close()

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timew1 time.Time

	//recot_un_w := timerm.CreateRecoTime()

	//	fp1 := "SUrecv.data"
	//	fp2 := "SUsend.data"
	//	fp1 = g_fpnull
	//	fp2 = g_fpnull

	//	_, _ = fp1, fp2

	var err error
	//	_ = err

	defer func() {
		g_RuntimeStatistics.AddSrcTotalNetData_send(int(obj.Wlent))
		g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tunpack  twlen:%d\n\t%s",
			obj.Gctx.GetGsId(), gstunnellib.GetNetConnAddrString("obj.Dst", obj.Dst), obj.Wlent,
			nt_write.PrintString(),
		)

		if g_Values.GetDebug() {
			//g_Logger.Println("\tgoUnpackTotal:", atomic.LoadInt32(&g_goUnpackTotal))

			//g_Logger.Println("\tRecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}()

	var ok bool

	for {

		obj.Wbuf, ok = <-obj.Dst_chan
		if !ok {
			g_Logger.Printf("Info: [%d] obj.Dst_chan is not ok, func exit.\n", obj.Gctx.GetGsId())
			return
		}
		if len(obj.Wbuf) <= 0 {
			continue
		}

		obj.Apack.WriteEncryData(obj.Wbuf)
		obj.Wbuf, err = obj.Apack.GetDecryData()
		checkError_panic_GsCtx(err, obj.Gctx)
		if len(obj.Wbuf) > 0 {
			obj.Dst.SetWriteDeadline(time.Now().Add(g_networkTimeout))
			timew1 = time.Now()
			rn, err := gstunnellib.NetConnWriteAll(obj.Dst, obj.Wbuf)
			obj.Wlent += int64(rn)
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				checkError_info_GsCtx(err, obj.Gctx)
				return
			} else {
				checkError_panic_GsCtx(err, obj.Gctx)
			}
			nt_write.Add(time.Since(timew1))
			//tmr_out.Boot()
		}
		/*
			if tmr_out.Run() {
				g_Logger.Printf("Error: [%d] Time out, func exit.\n", obj.Gctx.GetGsId())
				return
			}
		*/
		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("unpack  twlen:%d\n", obj.Wlent)

			//g_Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}

	}
}
