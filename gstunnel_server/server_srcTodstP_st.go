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

func srcTOdstP_st(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(time.Second * 60)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	tmr_changekey := timerm.CreateTimer(tmr_changekey_time)

	recot_p_r := timerm.CreateRecoTime()
	recot_p_w := timerm.CreateRecoTime()

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
	defer dst.Close()
	//defer outf.Close()
	//defer outf2.Close()
	buf := make([]byte, net_read_size)
	var rbuf []byte = buf
	var wbuf bytes.Buffer

	var wlent, rlent int64 = 0, 0

	ChangeCryKey_Total := 0

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_recv(int(rlent))
		GRuntimeStatistics.AddSrcTotalNetData_send(int(wlent))

		if debug_server {
			fmt.Println("gorou exit.")
			fmt.Printf("\tpack  trlen:%d  twlen:%d\n", rlent, wlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	err = IsTheVersionConsistent_send(dst, apack, &wlent)
	if err != nil {
		Logger.Println("Error:", err)
		return
	}
	err = ChangeCryKey_send(dst, apack, &ChangeCryKey_Total, &wlent)
	if err != nil {
		Logger.Println("Error:", err)
		return
	}

	for {
		buf = rbuf
		recot_p_r.Run()
		rlen, err := src.Read(rbuf)
		recot_p_r.Run()

		rlent = rlent + int64(rlen)

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

		//outf.Write(buf[:rlen])
		tmr_out.Boot()
		//rbuf = buf
		buf = rbuf[:rlen]

		wbuf.Reset()
		_, err = wbuf.Write(buf)
		if err != nil {
			Logger.Println("Error:", err)
			return
		}
		buf = nil
		//wbuf = append(wbuf, buf...)
		//fre := bool(len(wbuf) > 0)

		if wbuf.Len() > 0 {
			if debug_server {
				gstunnellib.RunTimeDebugInfoV1.AddPackingPackSizeList("server_srcToDstP_st_packing", wbuf.Len())
			}
			buf = apack.Packing(wbuf.Bytes())
			//wbuf = wbuf[len(wbuf):]
			//outf2.Write(buf)
			if len(buf) <= 0 {
				Logger.Println("Error: gspack.packing is error.")
				return
			}
			wlen, err := io.Copy(dst, bytes.NewBuffer(buf))
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				checkError_NoExit(err)
				return
			} else {
				checkError(err)
			}
			wlent = wlent + int64(wlen)
			tmr_out.Boot()

		}
		if tmr_changekey.Run() {
			err = ChangeCryKey_send(dst, apack, &ChangeCryKey_Total, &wlent)
			if err != nil {
				Logger.Println("Error:", err)
				return
			}
		}
		buf = rbuf
		if tmrP2.Run() && debug_server {
			fmt.Printf("pack  trlen:%d  twlen:%d\n", rlent, wlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}

	}
}
