/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ypcd/gstunnel/v6/gstunnellib"
	"github.com/ypcd/gstunnel/v6/timerm"
)

const version string = gstunnellib.Version

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string

var gsconfig *gstunnellib.GsConfig

var goPackTotal, goUnpackTotal int32 = 0, 0

var Logger *log.Logger

const networkTimeout time.Duration = time.Minute * 1

var debug_client bool = false

const net_read_size = 4 * 1024
const netPUn_chan_cache_size = 64

var Mt_model bool = true
var tmr_display_time = time.Second * 5
var tmr_changekey_time = time.Second * 60

var bufPool sync.Pool

var GRuntimeStatistics gstunnellib.Runtime_statistics

type buf_chan struct {
	Gbuf *[]byte
	Len  int
}

func newBufChan(inGbuf *[]byte, inLen int) *buf_chan {
	return &buf_chan{inGbuf, inLen}
}

func (bc *buf_chan) Bytes() []byte {
	return (*bc.Gbuf)[:bc.Len]
}

func init() {

	GRuntimeStatistics = gstunnellib.NewRuntimeStatistics()

	bufPool = sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			re := make([]byte, net_read_size)
			return &re
		},
	}

	Logger = gstunnellib.CreateFileLogger("gstunnel_client.err.log")

	Logger.Println("gstunnel client.")
	Logger.Println("VER:", version)

	gsconfig = gstunnellib.CreateGsconfig("config.client.json")
	debug_client = gsconfig.Debug

	Mt_model = gsconfig.Mt_model

	tmr_display_time = time.Second * time.Duration(gsconfig.Tmr_display_time)
	tmr_changekey_time = time.Second * time.Duration(gsconfig.Tmr_changekey_time)

	Logger.Println("debug:", debug_client)

	Logger.Println("Mt_model:", Mt_model)
	Logger.Println("tmr_display_time:", tmr_display_time)
	Logger.Println("tmr_changekey_time:", tmr_changekey_time)

	Logger.Println("info_protobuf:", gstunnellib.Info_protobuf)

	if debug_client {
		go func() {
			Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6060", nil))
		}()
		Logger.Println("Debug server listen: localhost:6060")
	}
	debug_client = false
	go gstunnellib.RunGRuntimeStatistics_print(Logger, GRuntimeStatistics)
}

func main() {
	for {
		run()
	}
}

func run() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprintln(os.Stderr, x)
			Logger.Println("Panic:", x, "  App restart.")
			tmp := make([]byte, 6000)
			nlen := runtime.Stack(tmp, true)
			Logger.Println("Panic stack:", string(tmp[:nlen]))

		}
	}()

	var lstnaddr string
	var connaddr []string

	lstnaddr = gsconfig.Listen
	connaddr = gsconfig.GetServers()
	key = gsconfig.Key

	fmt.Println("Listen_Addr:", lstnaddr)
	fmt.Println("Conn_Addr:", connaddr)
	fmt.Println("Begin......")

	service := lstnaddr
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		acc, err := listener.Accept()
		if err != nil {
			Logger.Println("Error:", err)
			continue
		}
		server_conn_error_total := 0
		tmr := timerm.CreateTimer(time.Second * 10)
		for {
			if tmr.Run() {
				Logger.Println("Error: server_conn_error_total: ", server_conn_error_total)
				server_conn_error_total = 0
				tmr.Boot()
			}
			service := gsconfig.GetServers()[0]

			dst, err := net.Dial("tcp", service)
			//checkError(err)
			if err != nil {
				if server_conn_error_total > 10000 {
					fmt.Fprintln(os.Stderr, "Error: server_conn_error_total > 10000")
					Logger.Fatalln("Error: server_conn_error_total > 10000")
				}
				server_conn_error_total++
				Logger.Println("Error: [net.Dial('tcp', service)]:", err.Error())
				fmt.Fprintln(os.Stderr, "Error: [net.Dial('tcp', service)]:", err.Error())

				continue
			}
			fmt.Println("conn.", service)

			//acc: 		src---client
			//dst: 		client---serever
			//pack: 	acc read recv, dst wirte send.
			//unpack:	dst read recv, acc wirte send.
			go srcTOdstP_count(acc, dst)
			go srcTOdstUn_count(dst, acc)
			break
		}
		fmt.Println("go.")
	}
}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
}

func srcTOdstP_count(src net.Conn, dst net.Conn) {
	atomic.AddInt32(&goPackTotal, 1)
	srcTOdstP(src, dst)
	atomic.AddInt32(&goPackTotal, -1)

}

func srcTOdstUn_count(src net.Conn, dst net.Conn) {
	atomic.AddInt32(&goUnpackTotal, 1)
	srcTOdstUn(src, dst)
	atomic.AddInt32(&goUnpackTotal, -1)
}

func IsTheVersionConsistent_send(dst net.Conn, apack gstunnellib.GsPack, wlent *int64) error {
	return gstunnellib.IsTheVersionConsistent_send(dst, apack, wlent)
}

func ChangeCryKey_send(dst net.Conn, apack gstunnellib.GsPack, ChangeCryKey_Total *int, wlent *int64) error {
	return gstunnellib.ChangeCryKey_send(dst, apack, ChangeCryKey_Total, wlent)
}

func srcTOdstP_st(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_exit(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
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

		if debug_client {
			fmt.Println("gorou exit.")
			fmt.Printf("\tpack  trlen:%d  twlen:%d\n", rlent, wlent)
			fmt.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))
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
		rlen, err := src.Read(buf)
		recot_p_r.Run()

		rlent = rlent + int64(rlen)

		if tmr_out.Run() {
			Logger.Println("Error: Time out func exit.")
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
			for len(buf) > 0 {
				if len(buf) > 0 {
					recot_p_w.Run()
					wlen, err := dst.Write(buf)
					recot_p_w.Run()
					if err != nil {
						Logger.Println("Error:", err)
					}

					wlent = wlent + int64(wlen)
					if wlen == 0 {
						Logger.Println("Error: dst.write() wlen==0 func exit.")
						return
					}
					if err != nil && wlen <= 0 {
						Logger.Println("Error: dst.write() wlen<=0 func exit.  ", err)
						continue
					}
					if len(buf) == wlen {
						break
					}
					buf = buf[wlen:]
				} else {
					break
				}
			}
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

			fmt.Printf("pack  trlen:%d  twlen:%d\n", rlent, wlent)
			fmt.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}

	}
}

func srcTOdstP_mt(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_exit(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	tmr_changekey := timerm.CreateTimer(tmr_changekey_time)

	recot_p_r := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer src.Close()

	var buf []byte
	var wlent, rlent int64 = 0, 0

	ChangeCryKey_Total := 0

	dst_chan := make(chan *buf_chan, netPUn_chan_cache_size)
	defer close(dst_chan)

	dst_ok := gstunnellib.CreateGorouStatus()

	defer func() {
		GRuntimeStatistics.AddSrcTotalNetData_recv(int(rlent))

		if debug_client {
			fmt.Println("\tgorou exit.")
			fmt.Printf("\t\tpack  trlen:%d\n", rlent)
			fmt.Println("\tgoPackTotal:", atomic.LoadInt32(&goPackTotal))
			fmt.Println("\tChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("\tRecoTime_p_r All: ", recot_p_r.StringAll())
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
		//buf = make([]byte, net_read_size)
		gbuf, ok := bufPool.Get().(*[]byte)
		if !ok {
			Logger.Println("Error: bufPool.Get().(*[]byte).")
			return
		}
		buf = *gbuf
		recot_p_r.Run()
		rlen, err := src.Read(buf)
		recot_p_r.Run()

		rlent = rlent + int64(rlen)

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
		buf = buf[:rlen]

		if len(buf) > 0 {
			buf = apack.Packing(buf)
			if len(*gbuf) == net_read_size {
				bufPool.Put(gbuf)
				gbuf = nil
			}
		} else {
			continue
		}
		tmpbuf := buf
		dst_chan <- newBufChan(&tmpbuf, len(buf))
		buf = nil

		if !dst_ok.IsOk() {
			Logger.Println("Error: not dst_ok.isok() func exit.")
			return
		}

		if tmr_changekey.Run() {
			var buf []byte = apack.ChangeCryKey()
			ChangeCryKey_Total += 1
			tmr_out.Boot()
			dst_chan <- newBufChan(&buf, len(buf))

		}
		if tmrP2.Run() && debug_client {
			fmt.Printf("pack  trlen:%d\n", rlent)
			fmt.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())

		}

	}
	Logger.Println("Func exit.")
}

func srcTOdstP_w(dst net.Conn, dst_chan chan *buf_chan, dst_ok *gstunnellib.Gorou_status, wlentotal int64) {
	defer gstunnellib.Panic_exit(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	recot_p_w := timerm.CreateRecoTime()

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	defer dst.Close()

	var wlent int64 = wlentotal
	ChangeCryKey_Total := 0

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_send(int(wlent))

		if debug_client {
			fmt.Println("gorou exit.")
			fmt.Printf("\tpack  twlen:%d\n", wlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
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

		chanbuf, ok := <-dst_chan
		if !ok {
			Logger.Println("Error: dst_chan is not ok, func exit.")
			return
		}
		var buf []byte = chanbuf.Bytes()
		if len(buf) <= 0 {
			continue
		}

		for len(buf) > 0 {
			if len(buf) > 0 {
				recot_p_w.Run()
				wlen, err := dst.Write(buf)
				recot_p_w.Run()

				wlent = wlent + int64(wlen)

				tmr_out.Boot()
				if wlen == 0 {
					Logger.Println("Error: dst.write() wlen==0 func exit.")
					return
				}
				if err != nil && wlen <= 0 {
					Logger.Println("Error: dst.write() wlen<=0 func exit.  ", err)
					continue
				}
				if len(buf) == wlen {
					break
				}
				buf = buf[wlen:]
			} else {
				break
			}
		}
		if len(*chanbuf.Gbuf) == net_read_size {
			bufPool.Put(chanbuf.Gbuf)
			chanbuf = nil
			buf = nil
		}
		if tmrP2.Run() && debug_client {
			fmt.Printf("pack twlen:%d\n", wlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}
		if tmr_out.Run() {
			Logger.Println("Error: Time out func exit.")
			return
		}
	}
}

func srcTOdstP(src net.Conn, dst net.Conn) {
	if Mt_model {
		srcTOdstP_mt(src, dst)
	} else {
		srcTOdstP_st(src, dst)
	}
}

func srcTOdstUn_st(src net.Conn, dst net.Conn) {
	defer src.Close()
	defer dst.Close()

	defer gstunnellib.Panic_exit(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	recot_un_r := timerm.CreateRecoTime()
	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(key)

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	_, _ = fp1, fp2

	var err error
	_ = err

	var buf []byte = make([]byte, net_read_size)
	var rbuf, wbuf []byte
	rbuf = buf
	wlent, rlent := 0, 0

	defer func() {
		GRuntimeStatistics.AddServerTotalNetData_recv(int(rlent))
		GRuntimeStatistics.AddSrcTotalNetData_send(int(wlent))

		if debug_client {
			fmt.Println("gorou exit.")
			fmt.Printf("\tunpack  trlen:%d  twlen:%d\n", rlent, wlent)
			fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}()

	for {
		buf = rbuf
		recot_un_r.Run()
		rlen, err := src.Read(buf)
		recot_un_r.Run()

		pf("trlen:%d  rlen:%d\n", rlent, rlen)

		rlent = rlent + rlen

		if tmr_out.Run() {
			Logger.Println("Error: Time out func exit.")
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

		wbuf = append(wbuf, buf...)
		buf = nil
		for {
			ix, fre := find0(wbuf)
			pf("ix: %d, fre: %t\n", ix, fre)
			if fre {
				buf = wbuf[:ix+1]
				wbuf = wbuf[ix+1:]
				pf("buf not unpack:%d\n", len(buf))
				buf, err = apack.Unpack(buf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
					Logger.Println(err.Error())
					return
				}
				pf("buf unpack:%d\n", len(buf))
				for {
					if len(buf) > 0 {
						recot_un_w.Run()
						wlen, err := dst.Write(buf)
						recot_un_w.Run()

						wlent = wlent + wlen

						pf("twlen:%d  wlen:%d\n", wlent, wlen)
						if wlen == 0 {
							Logger.Println("Error: dst.write() wlen==0 func exit.")
							return
						}
						if err != nil && wlen <= 0 {
							Logger.Println("Error: dst.write() wlen<=0 func exit.  ", err)
							continue
						}
						if len(buf) == wlen {
							break
						}
						buf = buf[wlen:]
					} else {
						break
					}
				}

			} else {
				break
			}
		}
		buf = rbuf
		if tmrP2.Run() && debug_client {
			fmt.Printf("unpack  trlen:%d  twlen:%d\n", rlent, wlent)
			fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}

}

func srcTOdstUn_mt(src net.Conn, dst net.Conn) {
	defer gstunnellib.Panic_exit(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	recot_un_r := timerm.CreateRecoTime()

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

	dst_chan := make(chan *buf_chan, netPUn_chan_cache_size)
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

			fmt.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}()

	for dst_ok.IsOk() {
		//buf = make([]byte, net_read_size)
		gbuf, ok := bufPool.Get().(*[]byte)
		if !ok {
			Logger.Println("Error: bufPool.Get().(*[]byte).")
			return
		}
		buf = *gbuf

		recot_un_r.Run()
		rlen, err := src.Read(buf)
		recot_un_r.Run()

		pf("trlen:%d  rlen:%d\n", rlent, rlen)

		rlent = rlent + int64(rlen)

		if tmr_out.Run() {
			Logger.Println("Error: Time out func exit.")
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

		dst_chan <- newBufChan(gbuf, len(buf))
		buf = nil
		gbuf = nil

		if tmrP2.Run() && debug_client {
			fmt.Printf("unpack  trlen:%d\n", rlent)
			fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}
	Logger.Println("Func exit.")
}

func srcTOdstUn_w(dst net.Conn, dst_chan chan *buf_chan, dst_ok *gstunnellib.Gorou_status) {
	defer gstunnellib.Panic_exit(Logger)

	tmr_out := timerm.CreateTimer(networkTimeout)

	tmrP2 := timerm.CreateTimer(tmr_display_time)

	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.NewGsPack(key)

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

			fmt.Println("\tRecoTime_un_w All: ", recot_un_w.StringAll())
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

		chanbuf, ok := <-dst_chan
		if !ok {
			Logger.Println("Error: dst_chan is not ok, func exit.")
			return
		}
		var buf []byte = chanbuf.Bytes()
		if len(buf) <= 0 {
			continue
		}

		wbuf = append(wbuf, buf...)
		buf = nil
		for {
			ix, fre := find0(wbuf)
			pf("ix: %d, fre: %t\n", ix, fre)
			if fre {
				buf = wbuf[:ix+1]
				wbuf = wbuf[ix+1:]
				pf("buf not unpack:%d\n", len(buf))
				buf, err = apack.Unpack(buf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
					Logger.Println(err.Error())
					return
				}
				pf("buf unpack:%d\n", len(buf))
				for {
					if len(buf) > 0 {
						recot_un_w.Run()
						wlen, err := dst.Write(buf)
						recot_un_w.Run()

						wlent = wlent + uint64(wlen)
						tmr_out.Boot()

						pf("twlen:%d  wlen:%d\n", wlent, wlen)
						if wlen == 0 {
							Logger.Println("Error: dst.write() wlen==0 func exit.")
							return
						}
						if err != nil && wlen <= 0 {
							Logger.Println("Error: dst.write() wlen<=0 func exit.  ", err)
							continue
						}
						if len(buf) == wlen {
							break
						}
						buf = buf[wlen:]
					} else {
						break
					}
				}

			} else {
				break
			}
		}
		if tmr_out.Run() {
			Logger.Println("Error: Time out func exit.")
			return
		}

		if tmrP2.Run() && debug_client {
			fmt.Printf("unpack  twlen:%d\n", wlent)

			fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
		if len(*chanbuf.Gbuf) == net_read_size {
			bufPool.Put(chanbuf.Gbuf)
			chanbuf = nil
			buf = nil
		}
	}
}

func srcTOdstUn(src net.Conn, dst net.Conn) {
	if Mt_model {
		srcTOdstUn_mt(src, dst)
	} else {
		srcTOdstUn_st(src, dst)
	}
}

func checkError(err error) {
	gstunnellib.CheckErrorEx(err, Logger)
}
