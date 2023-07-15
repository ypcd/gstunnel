package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstUn_st(src net.Conn, dst net.Conn, gctx gstunnellib.GsContext) {
	defer gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_Logger, gctx)

	defer src.Close()
	defer dst.Close()

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	//tmrP := timerm.CreateTimer(g_tmr_display_time)
	tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	nt_read := gstunnellib.NewNetTimeImpName("read")
	nt_write := gstunnellib.NewNetTimeImpName("write")
	var timer1, timew1 time.Time

	//	recot_un_r := timerm.CreateRecoTime()
	//	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPackNet(g_key)

	//	fp1 := "SUrecv.data"
	//	fp2 := "SUsend.data"
	//	fp1 = g_fpnull
	//	fp2 = g_fpnull
	//	_, _ = fp1, fp2

	//var err error
	//outf, err := os.Create(fp1)
	//checkError_GsCtx(err,gctx)
	//outf2, err := os.Create(fp2)
	//checkError_GsCtx(err,gctx)

	//defer outf.Close()
	//defer outf2.Close()
	var wbuf []byte
	var rbuf = make([]byte, g_net_read_size)

	var wlent, rlent uint64 = 0, 0

	defer func() {
		g_RuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))
		g_RuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\t%s\tunpack trlen:%d  twlen:%d\n\t%s\t%s",
			gctx.GetGsId(),
			gstunnellib.GetNetConnAddrString("src", src),
			gstunnellib.GetNetConnAddrString("dst", dst),
			rlent, wlent,
			nt_read.PrintString(),
			nt_write.PrintString(),
		)

		if g_Values.GetDebug() {

			//	g_Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	g_Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}()

	for {
		//	recot_un_r.Run()
		src.SetReadDeadline(time.Now().Add(g_networkTimeout))
		timer1 = time.Now()
		rlen, err := src.Read(rbuf)
		rlent += uint64(rlen)
		//	recot_un_r.Run()
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
		if rlen == 0 {
			g_Logger.Println("Error: src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			g_Logger.Println("Error:", err)
			continue
		}
		//tmr_out.Boot()

		apack.WriteEncryData(rbuf[:rlen])
		wbuf, err = apack.GetDecryData()
		checkError_panic_GsCtx(err, gctx)
		if len(wbuf) > 0 {
			dst.SetWriteDeadline(time.Now().Add(g_networkTimeout))
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
			nt_write.Add(time.Since(timew1))
			//tmr_out.Boot()
		}

		if tmrP2.Run() && g_Values.GetDebug() {
			g_Logger.Printf("unpack trlen:%d  twlen:%d\n", rlent, wlent)

			//	g_Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	g_Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}
}
