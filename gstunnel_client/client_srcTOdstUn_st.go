/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstUn_st(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)
	defer src.Close()
	defer dst.Close()

	//tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	nt_read := gstunnellib.NewNetTimeImpName("read")
	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timer1, timew1 time.Time

	//recot_un_r := timerm.CreateRecoTime()
	//recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPackNet(key)

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	var rbuf []byte = make([]byte, net_read_size)
	var wbuf []byte
	var wlent, rlent int64

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_recv(int(rlent))
		GRuntimeStatistics.AddSrcTotalNetData_send(int(wlent))
		log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\t%s\tunpack  trlen:%d  twlen:%d\n\t%s\t%s",
			gctx.GetGsId(),
			gstunnellib.GetNetConnAddrString("src", src),
			gstunnellib.GetNetConnAddrString("dst", dst),
			rlent, wlent,
			nt_read.PrintString(),
			nt_write.PrintString(),
		)

		if GValues.GetDebug() {

			//	Logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}()

	for {
		//	recot_un_r.Run()
		src.SetReadDeadline(time.Now().Add(networkTimeout))
		timer1 = time.Now()
		rlen, err := src.Read(rbuf)
		rlent += int64(rlen)
		//	recot_un_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info_GsCtx(err, gctx)
			return
		} else {
			checkError_panic_GsCtx(err, gctx)
		}
		nt_read.Add(time.Now().Sub(timer1))

		pf("trlen:%d  rlen:%d\n", rlent, rlen)
		/*
			if tmr_out.Run() {
				Logger.Printf("Error: [%d] Time out, func exit.\n", gctx.GetGsId())
				return
			}
		*/
		if rlen == 0 {
			Logger.Println("Error: src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			Logger.Println("Error:", err)
			continue
		}

		//tmr_out.Boot()

		apack.WriteEncryData(rbuf[:rlen])
		wbuf, err = apack.GetDecryData()
		checkError_panic_GsCtx(err, gctx)
		if len(wbuf) > 0 {
			dst.SetWriteDeadline(time.Now().Add(networkTimeout))
			timew1 = time.Now()
			rn, err := io.Copy(dst, bytes.NewBuffer(wbuf))
			wlent += rn
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
			Logger.Printf("unpack  trlen:%d  twlen:%d\n", rlent, wlent)
			Logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}

}
