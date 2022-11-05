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

func srcTOdstUn_w(dst net.Conn, dst_chan chan ([]byte), dst_ok gstunnellib.Gorou_status, wg_w *sync.WaitGroup, gctx gstunnellib.GsContext) {
	defer wg_w.Done()
	defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	defer func() {
		dst.Close()
		dst_ok.SetClose()
		gstunnellib.CloseChan(dst_chan)
	}()

	//tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timew1 time.Time
	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

	//	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPackNet(key)

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	_, _ = fp1, fp2

	var err error
	_ = err
	//outf, err := os.Create(fp1)
	//checkError_GsCtx(err,gctx)

	//outf2, err := os.Create(fp2)
	//checkError_GsCtx(err,gctx)

	//defer outf.Close()
	//defer outf2.Close()

	//buf := make([]byte, net_read_size)
	var wbuf []byte
	//var wbuff bytes.Buffer
	var wlent uint64 = 0

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tunpack  twlen:%d\n\t%s",
			gctx.GetGsId(), gstunnellib.GetNetConnAddrString("dst", dst), wlent,
			nt_write.PrintString(),
		)

		if GValues.GetDebug() {

			//Logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	Logger.Println("\tRecoTime_un_w All: ", recot_un_w.StringAll())
		}
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

		apack.WriteEncryData(buf)
		wbuf, err = apack.GetDecryData()
		checkError_panic_GsCtx(err, gctx)
		if len(wbuf) > 0 {
			dst.SetWriteDeadline(time.Now().Add(networkTimeout))
			timew1 = time.Now()
			rn, err := io.Copy(dst, bytes.NewBuffer(wbuf))
			wlent += uint64(rn)
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
		/*
			if tmr_out.Run() {
				Logger.Printf("Error: [%d] Time out, func exit.\n", gctx.GetGsId())
				return
			}
		*/
		if tmrP2.Run() && GValues.GetDebug() {
			Logger.Printf("unpack  twlen:%d\n", wlent)

			//	Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}
}
