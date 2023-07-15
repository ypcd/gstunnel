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

func srcTOdstP_st(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)
	defer dst.Close()
	defer src.Close()
	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	tmr_changekey := timerm.CreateTimer(g_tmr_changekey_time)

	nt_read := gstunnellib.NewNetTimeImpName("read")
	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timer1, timew1 time.Time

	//recot_p_r := timerm.CreateRecoTime()
	//recot_p_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(g_key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = g_fpnull, g_fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	var wbuf []byte
	var rbuf []byte = make([]byte, g_net_read_size)

	var wlent, rlent int64 = 0, 0

	ChangeCryKey_Total := 0

	defer func() {
		g_RuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))
		g_RuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		g_log_List.GSNetIOLen.Printf(
			"[%d] gorou exit.\n\t%s\t%s\tpack  trlen:%d  twlen:%d  ChangeCryKey_total:%d\n\t%s\t%s",
			gctx.GetGsId(),
			gstunnellib.GetNetConnAddrString("src", src),
			gstunnellib.GetNetConnAddrString("dst", dst),
			rlent, wlent, ChangeCryKey_Total,
			nt_read.PrintString(), nt_write.PrintString(),
		)

		if g_Values.GetDebug() {
			//g_Logger.Println("g_goPackTotal:", atomic.LoadInt32(&g_goPackTotal))

			//g_Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			//g_Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	err = IsTheVersionConsistent_send(dst, apack, &wlent)
	if err != nil {
		g_Logger.Println("Error:", err)
		return
	}
	err = ChangeCryKey_send(dst, apack, &ChangeCryKey_Total, &wlent)
	if err != nil {
		g_Logger.Println("Error:", err)
		return
	}
	for {
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
				g_Logger.Printf("Error: [%d] Time out, func exit.\n", gctx.GetGsId())
				return
			}
		*/
		if rlen <= 0 {
			g_Logger.Println("Error: src.read() rlen==0 func exit.")
			return
		}

		//tmr_out.Boot()

		if rlen > 0 {
			wbuf = apack.Packing(rbuf[:rlen])
			if len(wbuf) <= 0 {
				g_Logger.Println("Error: gspack.packing is error.")
				return
			}

			dst.SetWriteDeadline(time.Now().Add(g_networkTimeout))
			timew1 = time.Now()
			wlen, err := io.Copy(dst, bytes.NewBuffer(wbuf))
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
		if tmr_changekey.Run() {
			err = ChangeCryKey_send(dst, apack, &ChangeCryKey_Total, &wlent)
			if err != nil {
				g_Logger.Println("Error:", err)
				return
			}
		}
		if tmrP2.Run() && g_Values.GetDebug() {

			g_Logger.Printf("pack  trlen:%d  twlen:%d\n", rlent, wlent)
			g_Logger.Println("g_goPackTotal:", atomic.LoadInt32(&g_goPackTotal))
			g_Logger.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			//g_Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			//g_Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}

	}
}
