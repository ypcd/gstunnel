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

func srcTOdstP_w(dst net.Conn, dst_chan chan ([]byte), dst_ok gstunnellib.Gorou_status, wlentotal int64, wg_w *sync.WaitGroup, gctx gstunnellib.GsContext) {
	wg_w.Done()
	defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	//tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timew1 time.Time

	//	recot_p_w := timerm.CreateRecoTime()

	//apack := aespack

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	//outf, err := os.Create(fp1)
	//checkError_GsCtx(err,gctx)
	//outf2, err := os.Create(fp2)
	//checkError_GsCtx(err,gctx)
	defer dst.Close()
	//defer outf.Close()
	//defer outf2.Close()

	//var buf []byte
	//var wbuf bytes.Buffer

	var wlent int64 = wlentotal

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_send(int(wlent))
		log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tpack  twlen:%d\n\t%s",
			gctx.GetGsId(), gstunnellib.GetNetConnAddrString("dst", dst), wlent,
			nt_write.PrintString(),
		)

		if GValues.GetDebug() {
			//Logger.Println("goPackTotal:", goPackTotal)

			//	Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	defer func() {
		dst_ok.SetClose()
		gstunnellib.CloseChan(dst_chan)
	}()

	for {
		/*
			if tmr_out.Run() {
				Logger.Printf("Error: [%d] Time out, func exit.\n", gctx.GetGsId())
				return
			}
		*/
		buf, ok := <-dst_chan
		if !ok {
			Logger.Printf("Error: [%d] dst_chan is not ok, func exit.\n", gctx.GetGsId())
			return
		}
		if len(buf) <= 0 {
			continue
		}

		//wbuf.Reset()
		//wbuf.Write(buf)
		//wbuf = append(wbuf, buf...)
		//fre := bool(len(wbuf) > 0)

		//wbuf = wbuf[len(wbuf):]
		//outf2.Write(buf)
		//tmr_out.Boot()
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
		}

		if tmrP2.Run() && GValues.GetDebug() {
			Logger.Printf("pack twlen:%d\n", wlent)
			//Logger.Println("goPackTotal:", goPackTotal)

			//	Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}

	}
}
