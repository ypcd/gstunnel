package main

import (
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_st(obj *gsobj.GSTObj) {
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_logger, obj.Gctx)

	defer obj.Close()

	//tmr_out := timerm.CreateTimer(obj.NetworkTimeout)
	//tmrP := timerm.CreateTimer(obj.Tmr_display_time)
	//tmrP2 := timerm.CreateTimer(obj.Tmr_display_time)

	tmr_changekey := timerm.CreateTimer(g_tmr_changekey_time)

	var timer1, timew1 time.Time

	var err error
	_ = err

	defer func() {
		g_RuntimeStatistics.AddServerTotalNetData_recv(int(obj.Rlent))
		g_RuntimeStatistics.AddSrcTotalNetData_send(int(obj.Wlent))
		/*
			g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\t%s\tpack  trlen:%d  twlen:%d  ChangeCryKey_total:%d\n\t%s\t%s",
				obj.Gctx.GetGsId(),
				gstunnellib.GetNetConnAddrString("obj.Src", obj.Src),
				gstunnellib.GetNetConnAddrString("obj.Dst", obj.Dst),
				obj.Rlent, obj.Wlent, obj.ChangeCryKey_Total,
				obj.Nt_read.PrintString(),
				obj.Nt_write.PrintString(),
			)
		*/
		g_log_List.GSNetIOLen.Println(obj.StringWithGOExit())

		if g_Values.GetDebug() {
		}
	}()

	//err = obj.VersionPack_send()
	//checkError_panic_GSCtx(err, obj.Gctx)

	err = obj.ChangeCryKey_send()
	checkError_panic_GSCtx(err, obj.Gctx)

	for {
		//	recot_p_r.Run()
		//err = //obj.Src.SetReadDeadline(time.Now().Add(obj.NetworkTimeout))
		timer1 = time.Now()
		rlen, err := obj.ReadNetSrc(obj.Rbuf)
		obj.Rlent += int64(rlen)
		//	recot_p_r.Run()
		if gstunnellib.IsErrorNetUsually(err) {
			checkError_info_GSCtx(err, obj.Gctx)
			return
		} else {
			checkError_panic_GSCtx(err, obj.Gctx)
		}
		obj.Nt_read.Add(time.Since(timer1))
		/*
			if tmr_out.Run() {
				g_logger.Printf("Error: [%d] Time out, func exit.\n", obj.Gctx.GetGsId())
				return
			}
		*/
		if rlen == 0 {
			g_logger.Println("Error: obj.Src.read() rlen==0 func exit.")
			return
		}

		if rlen > 0 {
			if gstunnellib.G_RunTime_Debug {
				gstunnellib.G_RunTimeDebugInfo1.AddPackingPackSizeList("server_srcToDstP_st_packing len", rlen)
			}
			obj.Wbuf = obj.Packing(obj.Rbuf[:rlen])

			if len(obj.Wbuf) <= 0 {
				g_logger.Println("Error: gspack.packing is error.")
				return
			}
			//obj.SetWriteDeadline(time.Now().Add(obj.NetworkTimeout))
			timew1 = time.Now()
			//wlen, err := obj.NetConnWriteAll(obj.Wbuf)
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
		if tmr_changekey.Run() {
			err = obj.ChangeCryKey_send()
			if err != nil {
				g_logger.Println("Error:", err)
				return
			}
		}
		/*if tmrP2.Run() && g_Values.GetDebug() {
			g_logger.Printf("pack  trlen:%d  twlen:%d\n", obj.Rlent, obj.Wlent)
			//g_logger.Println("goPackTotal:", goPackTotal)
			g_logger.Println("ChangeCryKey_total:", obj.ChangeCryKey_Total)

			//	g_logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			//	g_logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}*/

	}
}
