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
	"fmt"
	"io"
	"net"
	"os"
	"sync/atomic"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstUn_mt(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//recot_un_r := timerm.CreateRecoTime()

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer src.Close()
	defer dst.Close()

	var buf []byte
	rlent := int64(0)

	dst_chan := make(chan []byte, netPUn_chan_cache_size)
	defer close(dst_chan)

	dst_ok := gstunnellib.CreateGorouStatus()

	go srcTOdstUn_w(dst, dst_chan, dst_ok)
	dst = nil

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_recv(int(rlent))

		if debug_client {
			fmt.Println("\tgorou exit.")
			fmt.Printf("\t\tunpack  trlen:%d\n", rlent)
			fmt.Println("\tgoUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	fmt.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}()

	for dst_ok.IsOk() {
		buf = make([]byte, net_read_size)

		//recot_un_r.Run()
		rlen, err := src.Read(buf)
		//recot_un_r.Run()
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info(err)
		} else {
			checkError(err)
		}

		pf("trlen:%d  rlen:%d\n", rlent, rlen)

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
		buf = buf[:rlen]
		if len(buf) <= 0 {
			continue
		}

		dst_chan <- buf
		buf = nil

		if tmrP2.Run() && debug_client {
			fmt.Printf("unpack  trlen:%d\n", rlent)
			fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}
	Logger.Println("Func exit.")
}

func srcTOdstUn_w(dst net.Conn, dst_chan chan []byte, dst_ok *gstunnellib.Gorou_status) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)

	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPackNet(key)

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull

	_, _ = fp1, fp2

	var err error
	_ = err

	defer dst.Close()

	var wbuf []byte
	var wlent uint64 = 0

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_send(int(wlent))

		if debug_client {
			fmt.Println("\tgorou exit.")
			fmt.Printf("\t\tunpack  twlen:%d\n", wlent)
			fmt.Println("\tgoUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//fmt.Println("\tRecoTime_un_w All: ", recot_un_w.StringAll())
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

		buf, ok := <-dst_chan
		if !ok {
			Logger.Println("Error: dst_chan is not ok, func exit.")
			return
		}
		if len(buf) <= 0 {
			continue
		}

		apack.WriteEncryData(buf)
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

		if tmr_out.Run() {
			Logger.Println("Error: Time out, func exit.")
			return
		}

		if tmrP2.Run() && debug_client {
			fmt.Printf("unpack  twlen:%d\n", wlent)

			//fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}

	}
}
