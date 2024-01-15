package main

import (
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/gstunnellib/gsobj"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_mt(obj *gsobj.GSTObj) {
	defer obj.Gctx.Close()
	defer gstunnellib.Panic_Recover_GSCtx(g_logger, obj.Gctx)
	defer obj.Close()

	//tmr_out := timerm.CreateTimer(obj.NetworkTimeout)
	//tmrP := timerm.CreateTimer(obj.Tmr_display_time)
	//tmrP2 := timerm.CreateTimer(obj.Tmr_display_time)

	tmr_changekey := timerm.CreateTimer(g_tmr_changekey_time)

	var timer1 time.Time
	//	recot_p_r := timerm.CreateRecoTime()

	//fp1, fp2 = g_fpnull, g_fpnull
	//_, _ = fp1, fp2

	var err error
	_ = err

	//outf, err := os.Create(fp1)
	//checkError_GSCtx(err,obj.Gctx)
	//outf2, err := os.Create(fp2)
	//checkError_GSCtx(err,obj.Gctx)
	//defer outf.Close()
	//defer outf2.Close()

	objw := obj.NewGstObjW(g_netPUn_chan_cache_size)

	defer func() {
		g_RuntimeStatistics.AddServerTotalNetData_recv(int(obj.Rlent))
		/*
			g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tpack  trlen:%d  ChangeCryKey_total:%d\n\t%s",
				obj.Gctx.GetGsId(), gstunnellib.GetNetConnAddrString("obj.Src", obj.Src), obj.Rlent, obj.ChangeCryKey_Total,
				obj.Nt_read.PrintString(),
			)*/
		g_log_List.GSNetIOLen.Println(obj.StringWithGOExit())

		if g_Values.GetDebug() {

			//g_logger.Println("goPackTotal:", goPackTotal)

			//	g_logger.Println("\tRecoTime_p_r All: ", recot_p_r.StringAll())
		}
	}()

	/*err = obj.VersionPack_send()
	checkError_panic_GSCtx(err, obj.Gctx)
	*/
	err = obj.ChangeCryKey_send()
	if err != nil {
		g_logger.Println("Error:", err, " func exit.")
		return
	}

	objw.Wg_w.Add(1)
	go srcTOdstP_w(objw)

	for objw.Dst_ok.IsOk() {
		//	recot_p_r.Run()
		//obj.Src.SetReadDeadline(time.Now().Add(obj.NetworkTimeout))
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
				g_logger.Println("Error: time out func exit.")
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

		//outf.Write(obj.Rbuf[:rlen])
		//tmr_out.Boot()

		if rlen > 0 {
			obj.Wbuf = obj.Packing(obj.Rbuf[:rlen])
		} else {
			continue
		}
		objw.Dst_chan <- obj.Wbuf
		obj.Wbuf = nil

		if !objw.Dst_ok.IsOk() {
			g_logger.Printf("Error: [%d] not objw.Dst_ok.isok() func exit.\n", obj.Gctx.GetGsId())
			return
		}

		if tmr_changekey.Run() {
			obj.Wbuf, _ = obj.ChangeCryKeyFromGSTServer()
			obj.ChangeCryKey_Total += 1
			objw.Dst_chan <- obj.Wbuf
			obj.Wbuf = nil
			//tmr_out.Boot()
		}

		/*if tmrP2.Run() && g_Values.GetDebug() {
			g_logger.Printf("pack  trlen:%d\n", obj.Rlent)
			//g_logger.Println("goPackTotal:", goPackTotal)
			g_logger.Println("ChangeCryKey_total:", obj.ChangeCryKey_Total)

			//	g_logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
		}*/
	}
	g_logger.Println("Func exit.")
}
