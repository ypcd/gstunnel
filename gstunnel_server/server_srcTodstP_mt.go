package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_mt(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	tmr_changekey := timerm.CreateTimer(tmr_changekey_time)

	//	recot_p_r := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	//outf, err := os.Create(fp1)
	//checkError(err)
	//outf2, err := os.Create(fp2)
	//checkError(err)
	defer src.Close()
	//defer outf.Close()
	//defer outf2.Close()
	rbuf := make([]byte, net_read_size)
	var buf []byte

	var rlent, wlent int64 = 0, 0

	ChangeCryKey_Total := 0

	dst_chan := make(chan ([]byte), netPUn_chan_cache_size)

	dst_ok := gstunnellib.CreateGorouStatus()
	defer func() {
		close(dst_chan)
		for {
			if !dst_ok.IsOk() {
				break
			}
		}
	}()

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_recv(int(rlent))
		Logger.Printf("gorou exit.\n%s\tpack  trlen:%d  ChangeCryKey_total:%d\n",
			gstunnellib.GetNetConnAddrString("src", src), rlent, ChangeCryKey_Total)

		if debug_server {

			//fmt.Println("goPackTotal:", goPackTotal)

			//	fmt.Println("\tRecoTime_p_r All: ", recot_p_r.StringAll())
		}
	}()

	err = IsTheVersionConsistent_send(dst, apack, &wlent)
	if err != nil {
		Logger.Println("Error:", err, " func exit.")
		return
	}
	err = ChangeCryKey_send(dst, apack, &ChangeCryKey_Total, &wlent)
	if err != nil {
		Logger.Println("Error:", err, " func exit.")
		return
	}

	go srcTOdstP_w(dst, dst_chan, dst_ok, wlent)
	dst = nil

	for dst_ok.IsOk() {
		//	recot_p_r.Run()
		rlen, err := src.Read(rbuf)
		//	recot_p_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info(err)
		} else {
			checkError(err)
		}

		rlent += int64(rlen)

		if tmr_out.Run() {
			Logger.Println("Error: time out func exit.")
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

		//outf.Write(buf[:rlen])
		tmr_out.Boot()
		buf = rbuf[:rlen]

		if len(buf) > 0 {
			buf = apack.Packing(buf)
		} else {
			continue
		}
		dst_chan <- buf
		buf = nil

		if !dst_ok.IsOk() {
			Logger.Println("Error: not dst_ok.isok() func exit.")
			return
		}

		if tmr_changekey.Run() {
			buf := apack.ChangeCryKey()
			ChangeCryKey_Total += 1
			tmr_out.Boot()
			//outf2.Write(buf)

			dst_chan <- buf
			buf = nil
		}

		if tmrP2.Run() && debug_server {
			fmt.Printf("pack  trlen:%d\n", rlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			//	fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
		}
	}
	Logger.Println("Func exit.")
}

func srcTOdstP_w(dst net.Conn, dst_chan chan ([]byte), dst_ok *gstunnellib.Gorou_status, wlentotal int64) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//	recot_p_w := timerm.CreateRecoTime()

	//apack := aespack

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	//outf, err := os.Create(fp1)
	//checkError(err)
	//outf2, err := os.Create(fp2)
	//checkError(err)
	defer dst.Close()
	//defer outf.Close()
	//defer outf2.Close()

	//var buf []byte
	//var wbuf bytes.Buffer

	var wlent int64 = wlentotal

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_send(int(wlent))
		Logger.Printf("gorou exit.\n%s\tpack  twlen:%d\n",
			gstunnellib.GetNetConnAddrString("dst", dst), wlent)

		if debug_server {
			//fmt.Println("goPackTotal:", goPackTotal)

			//	fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	defer func() {
		dst_ok.SetClose()
		for {
			_, ok := <-dst_chan
			if !ok {
				break
			}
		}
	}()

	for {
		if tmr_out.Run() {
			Logger.Println("Error: Time out, func exit.")
			return
		}

		buf, ok := <-dst_chan
		if !ok {
			Logger.Println("Error: dst_chan is not ok, func exit.")
			return
		}
		if len(buf) <= 0 {
			continue
		}

		//wbuf.Reset()
		//wbuf.Write(buf)
		//wbuf = append(wbuf, buf...)
		//fre := bool(len(wbuf) > 0)

		//wbuf = wbuf[len(wbuf):]
		//outf2.Write(buf)

		if len(buf) > 0 {
			wlen, err := io.Copy(dst, bytes.NewBuffer(buf))
			wlent += int64(wlen)
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				checkError_NoExit(err)
				return
			} else {
				checkError(err)
			}
			tmr_out.Boot()
		}

		if tmrP2.Run() && debug_server {
			fmt.Printf("pack twlen:%d\n", wlent)
			//fmt.Println("goPackTotal:", goPackTotal)

			//	fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}

	}
}
