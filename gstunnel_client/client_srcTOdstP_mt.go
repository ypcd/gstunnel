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
	"sync"
	"sync/atomic"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstP_mt(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	tmr_changekey := timerm.CreateTimer(tmr_changekey_time)

	//recot_p_r := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer src.Close()

	rbuf := make([]byte, net_read_size)
	var buf []byte
	var wlent, rlent int64 = 0, 0

	ChangeCryKey_Total := 0

	dst_chan := make(chan []byte, netPUn_chan_cache_size)
	wg_w := new(sync.WaitGroup)

	dst_ok := gstunnellib.NewGorouStatus()
	defer func() {
		gstunnellib.CloseChan(dst_chan)
		wg_w.Wait()
	}()

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))
		log_List.GSNetIOLen.Printf("gorou exit.\n\t%s\tpack  trlen:%d  ChangeCryKey_total:%d\n",
			gstunnellib.GetNetConnAddrString("src", src), rlent, ChangeCryKey_Total)

		if debug_client {
			//Logger.Println("\tgoPackTotal:", atomic.LoadInt32(&goPackTotal))

			//Logger.Println("\tRecoTime_p_r All: ", recot_p_r.StringAll())
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

	wg_w.Add(1)
	go srcTOdstP_w(dst, dst_chan, dst_ok, wlent, wg_w)
	dst = nil

	for dst_ok.IsOk() {

		//recot_p_r.Run()
		rlen, err := src.Read(rbuf)
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
			var buf []byte = apack.ChangeCryKey()
			ChangeCryKey_Total += 1
			tmr_out.Boot()
			dst_chan <- buf
			buf = nil
		}
		if tmrP2.Run() && debug_client {
			Logger.Printf("pack  trlen:%d\n", rlent)
			Logger.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))
			Logger.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			//Logger.Println("RecoTime_p_r All: ", recot_p_r.StringAll())

		}

	}
	Logger.Println("Func exit.")
}

func srcTOdstP_w(dst net.Conn, dst_chan chan []byte, dst_ok gstunnellib.Gorou_status, wlentotal int64, wg_w *sync.WaitGroup) {
	defer wg_w.Done()
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//recot_p_w := timerm.CreateRecoTime()

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer dst.Close()

	var wlent int64 = wlentotal

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		log_List.GSNetIOLen.Printf("gorou exit.\n\t%s\tpack  twlen:%d\n",
			gstunnellib.GetNetConnAddrString("dst", dst), wlent)

		if debug_client {
			//Logger.Println("goPackTotal:", goPackTotal)

			//Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
	}()

	defer func() {
		dst_ok.SetClose()
		_, ok := <-dst_chan
		if ok {
			gstunnellib.CloseChan(dst_chan)
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

		if len(buf) > 0 {
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

		if tmrP2.Run() && debug_client {
			Logger.Printf("pack twlen:%d\n", wlent)
			//Logger.Println("goPackTotal:", goPackTotal)

			//Logger.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
		}
		if tmr_out.Run() {
			Logger.Println("Error: Time out, func exit.")
			return
		}
	}
}
