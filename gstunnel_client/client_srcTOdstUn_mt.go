/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
)

func srcTOdstUn_mt(obj *gsobj.GSTObj) {
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_logger, obj.Gctx)
	defer obj.Close()

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	//tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	//obj.Nt_read := gstunnellib.NewNetTimeImpName("read")
	var timer1 time.Time

	//recot_un_r := timerm.CreateRecoTime()

	//	fp1 := "SUrecv.data"
	//	fp2 := "SUsend.data"
	//	fp1 = g_fpnull
	//	fp2 = g_fpnull
	//	_, _ = fp1, fp2

	//	var err error
	//	_ = err

	objw := obj.NewGstObjW(g_netPUn_chan_cache_size)

	objw.Wg_w.Add(1)
	go srcTOdstUn_w(objw)

	defer func() {
		g_RuntimeStatistics.AddServerTotalNetData_recv(int(obj.Rlent))
		/*
			g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tunpack  trlen:%d\n\t%s",
				obj.Gctx.GetGsId(), gstunnellib.GetNetConnAddrString("obj.Src", obj.Src), obj.Rlent,
				obj.Nt_read.PrintString(),
			)*/
		g_log_List.GSNetIOLen.Println(obj.StringWithGOExit())

		if g_Values.GetDebug() {
			//g_logger.Println("\tgoUnpackTotal:", atomic.LoadInt32(&g_goUnpackTotal))

			//	g_logger.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}()

	for objw.Dst_ok.IsOk() {
		obj.Rbuf = make([]byte, g_net_read_size)

		//recot_un_r.Run()
		//obj.Src.SetReadDeadline(time.Now().Add(g_networkTimeout))
		timer1 = time.Now()
		rlen, err := obj.ReadNetSrc(obj.Rbuf)
		obj.Rlent += int64(rlen)
		//recot_un_r.Run()
		if gstunnellib.IsErrorNetUsually(err) {
			checkError_info_GSCtx(err, obj.Gctx)
			return
		} else {
			checkError_panic_GSCtx(err, obj.Gctx)
		}
		obj.Nt_read.Add(time.Since(timer1))

		//pf("trlen:%d  rlen:%d\n", obj.Rlent, rlen)
		/*
			if tmr_out.Run() {
				g_logger.Printf("Error: [%d] Time out, func exit.\n", obj.Gctx.GetGsId())
				return
			}
		*/
		if rlen <= 0 {
			g_logger.Println("Error: obj.Src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			g_logger.Println("Error:", err)
			continue
		}

		//tmr_out.Boot()

		objw.Dst_chan <- obj.Rbuf[:rlen]
		obj.Rbuf = nil

		/*/*if tmrP2.Run() && g_Values.GetDebug() {
			g_logger.Printf("unpack  trlen:%d\n", obj.Rlent)
			g_logger.Println("g_goUnpackTotal:", atomic.LoadInt32(&g_goUnpackTotal))

			//g_logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}
		*/
	}
	g_logger.Println("Func exit.")
}
