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
	"sync"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstUn_mt(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	defer gstunnellib.Panic_Recover_GSCtx(Logger, gctx)

	//tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	nt_read := gstunnellib.NewNetTimeImpName("read")
	var timer1 time.Time

	//recot_un_r := timerm.CreateRecoTime()

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer src.Close()
	defer dst.Close()

	var buf []byte
	rlent := int64(0)

	dst_chan := make(chan []byte, netPUn_chan_cache_size)
	wg_w := new(sync.WaitGroup)

	dst_ok := gstunnellib.NewGorouStatus()
	defer func() {
		gstunnellib.CloseChan(dst_chan)
		wg_w.Wait()
	}()

	wg_w.Add(1)
	go srcTOdstUn_w(dst, dst_chan, dst_ok, wg_w, gctx)
	dst = nil

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_recv(int(rlent))
		log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tunpack  trlen:%d\n\t%s",
			gctx.GetGsId(), gstunnellib.GetNetConnAddrString("src", src), rlent,
			nt_read.PrintString(),
		)

		if GValues.GetDebug() {
			//Logger.Println("\tgoUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	Logger.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}()

	for dst_ok.IsOk() {
		buf = make([]byte, net_read_size)

		//recot_un_r.Run()
		src.SetReadDeadline(time.Now().Add(networkTimeout))
		timer1 = time.Now()
		rlen, err := src.Read(buf)
		rlent += int64(rlen)
		//recot_un_r.Run()
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
		buf = buf[:rlen]
		if len(buf) <= 0 {
			continue
		}

		dst_chan <- buf
		buf = nil

		if tmrP2.Run() && GValues.GetDebug() {
			Logger.Printf("unpack  trlen:%d\n", rlent)
			Logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}
	Logger.Println("Func exit.")
}
