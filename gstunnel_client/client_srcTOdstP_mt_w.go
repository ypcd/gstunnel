package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_w(dst net.Conn, dst_chan chan []byte, dst_ok gstunnellib.Gorou_status, wlentotal int64, wg_w *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg_w.Done()
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timew1 time.Time
	//recot_p_w := timerm.CreateRecoTime()

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = g_fpnull, g_fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer dst.Close()

	var wlent int64 = wlentotal

	defer func() {
		g_RuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tpack  twlen:%d\n\t%s",
			gctx.GetGsId(), gstunnellib.GetNetConnAddrString("dst", dst), wlent,
			nt_write.PrintString(),
		)

		if g_Values.GetDebug() {
			//g_Logger.Println("g_goPackTotal:", g_goPackTotal)

			//g_Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	defer func() {
		dst_ok.SetClose()
		gstunnellib.ChanClean(dst_chan)
	}()

	for {

		buf, ok := <-dst_chan
		if !ok {
			g_Logger.Printf("Info: [%d] dst_chan is not ok, func exit.\n", gctx.GetGsId())
			return
		}
		if len(buf) <= 0 {
			continue
		}

		if len(buf) > 0 {
			dst.SetWriteDeadline(time.Now().Add(g_networkTimeout))
			timew1 = time.Now()
			wlen, err := io.Copy(dst, bytes.NewBuffer(buf))
			wlent += int64(wlen)
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				checkError_info_GsCtx(err, gctx)
				return
			} else {
				checkError_panic_GsCtx(err, gctx)
			}
			nt_write.Add(time.Since(timew1))
			//tmr_out.Boot()
		}

		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("pack twlen:%d\n", wlent)
			//g_Logger.Println("g_goPackTotal:", g_goPackTotal)

			//g_Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
		/*
			if tmr_out.Run() {
				g_Logger.Printf("Error: [%d] Time out, func exit.\n", gctx.GetGsId())
				return
			}
		*/
	}
}
