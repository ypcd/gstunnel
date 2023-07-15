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
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	nt_read := gstunnellib.NewNetTimeImpName("read")
	var timer1 time.Time

	//recot_un_r := timerm.CreateRecoTime()

	//	fp1 := "SUrecv.data"
	//	fp2 := "SUsend.data"
	//	fp1 = g_fpnull
	//	fp2 = g_fpnull
	//	_, _ = fp1, fp2

	//	var err error
	//	_ = err

	var buf []byte
	rlent := int64(0)

	dst_chan := make(chan []byte, g_netPUn_chan_cache_size)
	wg_w := new(sync.WaitGroup)

	dst_ok := gstunnellib.NewGorouStatusNetConn([]net.Conn{src, dst})
	defer func() {
		dst.Close()
		src.Close()
		gstunnellib.ChanClose(dst_chan)
		wg_w.Wait()
	}()

	wg_w.Add(1)
	go srcTOdstUn_w(dst, dst_chan, dst_ok, wg_w, gctx)

	defer func() {
		g_RuntimeStatistics.AddServerTotalNetData_recv(int(rlent))
		g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tunpack  trlen:%d\n\t%s",
			gctx.GetGsId(), gstunnellib.GetNetConnAddrString("src", src), rlent,
			nt_read.PrintString(),
		)

		if g_Values.GetDebug() {
			//g_Logger.Println("\tgoUnpackTotal:", atomic.LoadInt32(&g_goUnpackTotal))

			//	g_Logger.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}()

	for dst_ok.IsOk() {
		buf = make([]byte, g_net_read_size)

		//recot_un_r.Run()
		src.SetReadDeadline(time.Now().Add(g_networkTimeout))
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
		nt_read.Add(time.Since(timer1))

		//pf("trlen:%d  rlen:%d\n", rlent, rlen)
		/*
			if tmr_out.Run() {
				g_Logger.Printf("Error: [%d] Time out, func exit.\n", gctx.GetGsId())
				return
			}
		*/
		if rlen <= 0 {
			g_Logger.Println("Error: src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			g_Logger.Println("Error:", err)
			continue
		}

		//tmr_out.Boot()

		dst_chan <- buf[:rlen]
		buf = nil

		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("unpack  trlen:%d\n", rlent)
			g_Logger.Println("g_goUnpackTotal:", atomic.LoadInt32(&g_goUnpackTotal))

			//g_Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}
	g_Logger.Println("Func exit.")
}
