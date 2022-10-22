package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstUn_st(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(time.Second * 60)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//	recot_un_r := timerm.CreateRecoTime()
	//	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPackNet(key)

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	_, _ = fp1, fp2

	var err error
	_ = err
	//outf, err := os.Create(fp1)
	//checkError(err)
	//outf2, err := os.Create(fp2)
	//checkError(err)

	defer src.Close()
	defer dst.Close()
	//defer outf.Close()
	//defer outf2.Close()
	var rbuf, wbuf []byte
	rbuf = make([]byte, net_read_size)

	var wlent, rlent uint64 = 0, 0

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))
		GRuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		Logger.Printf("gorou exit.\n%s%s\tunpack trlen:%d  twlen:%d\n",
			gstunnellib.GetNetConnAddrString("src", src),
			gstunnellib.GetNetConnAddrString("dst", dst),
			rlent, wlent)

		if debug_server {

			//	fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}()

	for {
		//	recot_un_r.Run()
		rlen, err := src.Read(rbuf)
		//	recot_un_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info(err)
		} else {
			checkError(err)
		}
		rlent += uint64(rlen)

		if tmr_out.Run() {
			Logger.Println("Error: Time out, func exit.")
			return
		}
		if rlen == 0 {
			Logger.Println("Error: src.read() rlen==0 func exit.")
			return
		}
		if err != nil {
			Logger.Println("Error:", err)
			continue
		}
		tmr_out.Boot()

		apack.WriteEncryData(rbuf[:rlen])
		wbuf, err = apack.GetDecryData()
		checkError(err)
		if len(wbuf) > 0 {
			rn, err := io.Copy(dst, bytes.NewBuffer(wbuf))
			wlent += uint64(rn)
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, io.EOF) {
				checkError_NoExit(err)
				return
			} else {
				checkError(err)
			}
		}

		if tmrP2.Run() && debug_server {
			fmt.Printf("unpack trlen:%d  twlen:%d\n", rlent, wlent)

			//	fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			//	fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}
}
