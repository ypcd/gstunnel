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
	defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	//tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timew1 time.Time
	//recot_p_w := timerm.CreateRecoTime()

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer dst.Close()

	var wlent int64 = wlentotal

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tpack  twlen:%d\n\t%s",
			gctx.GetGsId(), gstunnellib.GetNetConnAddrString("dst", dst), wlent,
			nt_write.PrintString(),
		)

		if GValues.GetDebug() {
			//Logger.Println("goPackTotal:", goPackTotal)

			//Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	defer func() {
		dst_ok.SetClose()
		gstunnellib.CloseChan(dst_chan)
	}()

	for {

		buf, ok := <-dst_chan
		if !ok {
			Logger.Printf("Error: [%d] dst_chan is not ok, func exit.\n", gctx.GetGsId())
			return
		}
		if len(buf) <= 0 {
			continue
		}

		if len(buf) > 0 {
			dst.SetWriteDeadline(time.Now().Add(networkTimeout))
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
			nt_write.Add(time.Now().Sub(timew1))
			//tmr_out.Boot()
		}

		if tmrP2.Run() && GValues.GetDebug() {
			Logger.Printf("pack twlen:%d\n", wlent)
			//Logger.Println("goPackTotal:", goPackTotal)

			//Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
		/*
			if tmr_out.Run() {
				Logger.Printf("Error: [%d] Time out, func exit.\n", gctx.GetGsId())
				return
			}
		*/
	}
}
