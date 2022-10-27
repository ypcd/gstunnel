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

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_st(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	tmr_changekey := timerm.CreateTimer(tmr_changekey_time)

	//recot_p_r := timerm.CreateRecoTime()
	//recot_p_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer src.Close()
	defer dst.Close()

	var buf []byte = make([]byte, net_read_size)
	var rbuf []byte = buf
	var wbuf bytes.Buffer

	var wlent, rlent int64 = 0, 0

	ChangeCryKey_Total := 0

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))
		GRuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		log_List.GSNetIOLen.Printf("gorou exit.\n\t%s\t%s\tpack  trlen:%d  twlen:%d  ChangeCryKey_total:%d\n",
			gstunnellib.GetNetConnAddrString("src", src),
			gstunnellib.GetNetConnAddrString("dst", dst),
			rlent, wlent, ChangeCryKey_Total)

		if debug_client {
			//Logger.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))

			//Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			//Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
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
		//recot_p_r.Run()
		rlen, err := src.Read(buf)
		//recot_p_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info(err)
			return
		} else {
			checkError_panic(err)

		}

		rlent += int64(rlen)

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
		rbuf = buf
		buf = buf[:rlen]

		wbuf.Reset()
		_, err = wbuf.Write(buf)
		if err != nil {
			Logger.Println("Error:", err)
			return
		}
		buf = nil

		if wbuf.Len() > 0 {
			buf = apack.Packing(wbuf.Bytes())
			//wbuf = wbuf[len(wbuf):]
			//outf2.Write(buf)
			if len(buf) <= 0 {
				Logger.Println("Error: gspack.packing is error.")
				return
			}
			wlen, err := io.Copy(dst, bytes.NewBuffer(buf))
			wlent += int64(wlen)
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
				checkError_info(err)
				return
			} else {
				checkError_panic(err)
			}
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
		if tmrP2.Run() && debug_client {

			Logger.Printf("pack  trlen:%d  twlen:%d\n", rlent, wlent)
			Logger.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))
			Logger.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			//Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			//Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}

	}
}
