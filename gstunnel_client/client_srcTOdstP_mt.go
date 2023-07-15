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

func srcTOdstP_mt(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	tmr_changekey := timerm.CreateTimer(g_tmr_changekey_time)

	nt_read := gstunnellib.NewNetTimeImpName("read")
	var timer1 time.Time

	//recot_p_r := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(g_key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = g_fpnull, g_fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	rbuf := make([]byte, g_net_read_size)
	var buf []byte
	var wlent, rlent int64 = 0, 0

	ChangeCryKey_Total := 0

	dst_chan := make(chan []byte, g_netPUn_chan_cache_size)
	wg_w := new(sync.WaitGroup)

	dst_ok := gstunnellib.NewGorouStatusNetConn([]net.Conn{src, dst})
	defer func() {
		dst.Close()
		src.Close()
		gstunnellib.ChanClose(dst_chan)
		wg_w.Wait()
	}()

	defer func() {
		g_RuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))
		g_log_List.GSNetIOLen.Printf(
			"[%d] gorou exit.\n\t%s\tpack  trlen:%d  ChangeCryKey_total:%d\n\t%s",
			gctx.GetGsId(), gstunnellib.GetNetConnAddrString("src", src), rlent, ChangeCryKey_Total,
			nt_read.PrintString(),
		)

		if g_Values.GetDebug() {
			//g_Logger.Println("\tgoPackTotal:", atomic.LoadInt32(&g_goPackTotal))

			//g_Logger.Println("\tRecoTime_p_r All: ", recot_p_r.StringAll())
		}
	}()

	err = IsTheVersionConsistent_send(dst, apack, &wlent)
	if err != nil {
		g_Logger.Println("Error:", err, " func exit.")
		return
	}
	err = ChangeCryKey_send(dst, apack, &ChangeCryKey_Total, &wlent)
	if err != nil {
		g_Logger.Println("Error:", err, " func exit.")
		return
	}

	wg_w.Add(1)
	go srcTOdstP_w(dst, dst_chan, dst_ok, wlent, wg_w, gctx)

	for dst_ok.IsOk() {

		//recot_p_r.Run()
		src.SetReadDeadline(time.Now().Add(g_networkTimeout))
		timer1 = time.Now()
		rlen, err := src.Read(rbuf)
		rlent += int64(rlen)
		//recot_p_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info_GsCtx(err, gctx)
			return
		} else {
			checkError_panic_GsCtx(err, gctx)
		}
		nt_read.Add(time.Since(timer1))
		/*
			if tmr_out.Run() {
				g_Logger.Println("Error: time out func exit.")
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

		if rlen > 0 {
			buf = apack.Packing(rbuf[:rlen])
		} else {
			continue
		}

		dst_chan <- buf
		buf = nil

		if !dst_ok.IsOk() {
			g_Logger.Printf("Error: [%d] not dst_ok.isok() func exit.\n", gctx.GetGsId())
			return
		}

		if tmr_changekey.Run() {
			buf = apack.ChangeCryKey()
			ChangeCryKey_Total += 1
			dst_chan <- buf
			buf = nil
			//tmr_out.Boot()
		}
		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("pack  trlen:%d\n", rlent)
			g_Logger.Println("g_goPackTotal:", atomic.LoadInt32(&g_goPackTotal))
			g_Logger.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			//g_Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())

		}

	}
	g_Logger.Println("Func exit.")
}
