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

	//tmr_out := timerm.CreateTimer(obj.NetworkTimeout)
	//tmrP := timerm.CreateTimer(obj.Tmr_display_time)
	//tmrP2 := timerm.CreateTimer(obj.Tmr_display_time)

	var timer1 time.Time
	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

	//	recot_un_r := timerm.CreateRecoTime()

	//obj := gstunnellib.NewGsPack(obj.Key)

	//	fp1 := "SUrecv.data"
	//	fp2 := "SUsend.data"
	//	fp1 = g_fpnull
	//	fp2 = g_fpnull
	//	_, _ = fp1, fp2

	//	var err error
	//	_ = err
	//outf, err := os.Create(fp1)
	//checkError_GSCtx(err,obj.Gctx)

	//outf2, err := os.Create(fp2)
	//checkError_GSCtx(err,obj.Gctx)

	//defer outf.Close()
	//defer outf2.Close()
	//var obj.Rbuf []byte
	//var wbuff bytes.Buffer

	objw := obj.NewGstObjW(g_netPUn_chan_cache_size)

	objw.Wg_w.Add(1)
	go srcTOdstUn_w(objw)

	defer func() {
		g_RuntimeStatistics.AddSrcTotalNetData_recv(int(obj.Rlent))
		/*
			g_log_List.GSNetIOLen.Printf("[%d] gorou exit.\n\t%s\tunpack  trlen:%d\n\t%s",
				obj.Gctx.GetGsId(), gstunnellib.GetNetConnAddrString("obj.Src", obj.Src), obj.Rlent,
				obj.Nt_read.PrintString(),
			)*/
		g_log_List.GSNetIOLen.Println(obj.StringWithGOExit())

		if g_Values.GetDebug() {
			//g_logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	g_logger.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}()

	for objw.Dst_ok.IsOk() {
		obj.Rbuf = make([]byte, g_net_read_size)

		//	recot_un_r.Run()
		//obj.Src.SetReadDeadline(time.Now().Add(obj.NetworkTimeout))
		timer1 = time.Now()
		rlen, err := obj.ReadNetSrc(obj.Rbuf)
		obj.Rlent += int64(rlen)
		//	recot_un_r.Run()
		//recot_un_r.RunDisplay("---recot_un_r:")
		//g_logger.Println("---rlen:", rlen)
		if gstunnellib.IsErrorNetUsually(err) {
			checkError_info_GSCtx(err, obj.Gctx)
			return
		} else {
			checkError_panic_GSCtx(err, obj.Gctx)
		}
		obj.Nt_read.Add(time.Since(timer1))

		//pf("trlen:%d  rlen:%d\n", obj.Rlent, rlen)

		//rlen = 0
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

		objw.Dst_chan <- obj.Rbuf[:rlen]
		obj.Rbuf = nil

		/*if tmrP2.Run() && g_Values.GetDebug() {
			g_logger.Printf("unpack  trlen:%d\n", obj.Rlent)
			//g_logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	g_logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}*/
	}
	g_logger.Println("Func exit.")
}
