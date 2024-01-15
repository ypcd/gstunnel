package main

import (
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
)

func srcTOdstP_w(obj *gsobj.GSTObjW) {
	defer obj.Wg_w.Done()
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_logger, obj.Gctx)
	defer obj.Close()

	//tmr_out := timerm.CreateTimer(g_networkTimeout)
	//tmrP2 := timerm.CreateTimer(g_tmr_display_time)

	//obj.Nt_write := gstunnellib.NewNetTimeImpName("write")
	//obj.Nt_write
	var timew1 time.Time
	//recot_p_w := timerm.CreateRecoTime()

	//fp1, fp2 := g_fpnull, g_fpnull
	//_, _ = fp1, fp2

	var err error
	_ = err

	//	var obj.Wlent int64 = wlentotal

	defer func() {
		g_RuntimeStatistics.AddServerTotalNetData_send(int(obj.Wlent))
		/*
			g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tpack  twlen:%d\n\t%s",
				obj.Gctx.GetGsId(), gstunnellib.GetNetConnAddrString("obj.Dst", obj.Dst), obj.Wlent,
				obj.Nt_write.PrintString(),
			)
		*/
		g_log_List.GSNetIOLen.Println(obj.StringWithGOExit())

		if g_Values.GetDebug() {
			//g_logger.Println("g_goPackTotal:", g_goPackTotal)

			//g_logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	var ok bool

	for {

		obj.Wbuf, ok = <-obj.Dst_chan
		if !ok {
			g_logger.Printf("Info: [%d] obj.Dst_chan is not ok, func exit.\n", obj.Gctx.GetGsId())
			return
		}
		if len(obj.Wbuf) <= 0 {
			continue
		}

		if len(obj.Wbuf) > 0 {
			//obj.Dst.SetWriteDeadline(time.Now().Add(g_networkTimeout))
			timew1 = time.Now()
			wlen, err := obj.NetConnWriteAll(obj.Wbuf)
			obj.Wlent += int64(wlen)
			if gstunnellib.IsErrorNetUsually(err) {
				checkError_info_GSCtx(err, obj.Gctx)
				return
			} else {
				checkError_panic_GSCtx(err, obj.Gctx)
			}
			obj.Nt_write.Add(time.Since(timew1))
			//tmr_out.Boot()
		}
		/*
			/*if tmrP2.Run() && g_Values.GetDebug() {
				g_logger.Printf("pack twlen:%d\n", obj.Wlent)
				//g_logger.Println("g_goPackTotal:", g_goPackTotal)

				//g_logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
			}
		*/
		/*
			if tmr_out.Run() {
				g_logger.Printf("Error: [%d] Time out, func exit.\n", obj.Gctx.GetGsId())
				return
			}
		*/
	}
}
