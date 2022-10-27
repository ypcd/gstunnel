package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"sync"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

func srcTOdstUn_mt(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_Recover(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

	//	recot_un_r := timerm.CreateRecoTime()

	//apack := gstunnellib.NewGsPack(key)

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

	//defer outf.Close()
	//defer outf2.Close()
	buf := make([]byte, net_read_size)
	//var rbuf []byte
	//var wbuff bytes.Buffer
	rlent := int64(0)

	dst_chan := make(chan ([]byte), netPUn_chan_cache_size)
	wg_w := new(sync.WaitGroup)

	dst_ok := gstunnellib.NewGorouStatus()
	defer func() {
		defer src.Close()
		gstunnellib.CloseChan(dst_chan)
		wg_w.Wait()
	}()

	wg_w.Add(1)
	go srcTOdstUn_w(dst, dst_chan, dst_ok, wg_w)
	dst = nil

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))
		log_List.GSNetIOLen.Printf("gorou exit.\n\t%s\tunpack  trlen:%d\n",
			gstunnellib.GetNetConnAddrString("src", src), rlent)

		if debug_server {
			//Logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	Logger.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}()

	for dst_ok.IsOk() {
		buf = make([]byte, net_read_size)
		//	recot_un_r.Run()
		rlen, err := src.Read(buf)
		//	recot_un_r.Run()
		//recot_un_r.RunDisplay("---recot_un_r:")
		//Logger.Println("---rlen:", rlen)
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) ||
			errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) {
			checkError_info(err)
			return
		} else {
			checkError_panic(err)
		}

		pf("trlen:%d  rlen:%d\n", rlent, rlen)

		rlent += int64(rlen)
		//rlen = 0

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
		buf = buf[:rlen]

		//wbuff.Reset()
		//wbuff.Write(buf)

		//wbuf = wbuff.Bytes()
		dst_chan <- buf
		buf = nil

		if tmrP2.Run() && debug_server {
			Logger.Printf("unpack  trlen:%d\n", rlent)
			//Logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	Logger.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}
	Logger.Println("Func exit.")
}

func srcTOdstUn_w(dst net.Conn, dst_chan chan ([]byte), dst_ok gstunnellib.Gorou_status, wg_w *sync.WaitGroup) {
	defer wg_w.Done()
	defer gstunnellib.Panic_Recover(Logger)

	defer func() {
		dst.Close()
		dst_ok.SetClose()
		gstunnellib.CloseChan(dst_chan)
	}()

	tmr_out := timerm.CreateTimer(networkTimeout)

	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

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

	//defer outf.Close()
	//defer outf2.Close()

	//buf := make([]byte, net_read_size)
	var wbuf []byte
	//var wbuff bytes.Buffer
	var wlent uint64 = 0

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_send(int(wlent))
		log_List.GSNetIOLen.Printf("gorou exit.\n\t%s\tunpack  twlen:%d\n",
			gstunnellib.GetNetConnAddrString("dst", dst), wlent)

		if debug_server {

			//Logger.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			//	Logger.Println("\tRecoTime_un_w All: ", recot_un_w.StringAll())
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
		checkError_panic(err)
		if len(wbuf) > 0 {
			rn, err := io.Copy(dst, bytes.NewBuffer(wbuf))
			wlent += uint64(rn)
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, io.EOF) {
				checkError_info(err)
				return
			} else {
				checkError_panic(err)
			}
		}
		if tmr_out.Run() {
			Logger.Println("Error: Time out, func exit.")
			return
		}

		if tmrP2.Run() && debug_server {
			Logger.Printf("unpack  twlen:%d\n", wlent)

			//	Logger.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}
}
