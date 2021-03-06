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
	"gstunnellib"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
	"timerm"
)

const version string = gstunnellib.Version

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string

var gsconfig *gstunnellib.GsConfig

var Logger *log.Logger

var debug_server bool = false

const net_read_size = 4 * 1024
const netPUn_chan_cache_size = 64

var Mt_model bool = true
var tmr_display_time = time.Second * 5
var tmr_changekey_time = time.Second * 60

const networkTimeout time.Duration = time.Minute * 1

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func init() {
	Logger = gstunnellib.CreateFileLogger("gstunnel_server.err.log")

	fmt.Println("VER:", version)
	Logger.Println("VER:", version)

	gsconfig = gstunnellib.CreateGsconfig("config.server.json")

	debug_server = gsconfig.Debug

	Mt_model = gsconfig.Mt_model

	tmr_display_time = time.Second * time.Duration(gsconfig.Tmr_display_time)
	tmr_changekey_time = time.Second * time.Duration(gsconfig.Tmr_changekey_time)

	fmt.Println("debug:", debug_server)
	Logger.Println("debug:", debug_server)

	fmt.Println("Mt_model:", Mt_model)
	Logger.Println("Mt_model:", Mt_model)
	fmt.Println("tmr_display_time:", tmr_display_time)
	Logger.Println("tmr_display_time:", tmr_display_time)
	fmt.Println("tmr_changekey_time:", tmr_changekey_time)
	Logger.Println("tmr_changekey_time:", tmr_changekey_time)

	fmt.Println("info_protobuf:", gstunnellib.Info_protobuf)
	Logger.Println("info_protobuf:", gstunnellib.Info_protobuf)

	if debug_server {
		go func() {
			Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6070", nil))
		}()
		fmt.Println("Debug server listen: localhost:6070")
		Logger.Println("Debug server listen: localhost:6070")
	}

}

func main() {
	for {
		run()
	}
}

func run() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("App restart.")
			Logger.Println(x, "  App restart.")
		}
	}()

	var lstnaddr, connaddr string

	if len(os.Args) == 4 {
		lstnaddr = os.Args[1]
		connaddr = os.Args[2]
		key = os.Args[3]
	} else {
		lstnaddr = gsconfig.Listen
		connaddr = gsconfig.GetServer()
		key = gsconfig.Key
	}

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
			continue
		}

		service := connaddr
		//tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		//fmt.Println(tcpAddr)
		//checkError(err)
		dst, err := net.Dial("tcp", service)
		checkError(err)
		fmt.Println("conn.")

		go srcTOdstUn(acc, dst)
		go srcTOdstP(dst, acc)
		fmt.Println("go.")
	}
}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
}

func srcTOdstP_old(src net.Conn, dst net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
			err := fmt.Errorf("Error:%s", x)
			fmt.Println(err)
			Logger.Println(x, "  Go exit.")
		}
	}()

	tmr := timerm.CreateTimer(time.Second * 60)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	tmr_changekey := timerm.CreateTimer(tmr_changekey_time)

	recot_p_r := timerm.CreateRecoTime()
	recot_p_w := timerm.CreateRecoTime()

	apack := gstunnellib.CreateAesPack(key)

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
	var rbuf []byte
	var wbuf bytes.Buffer

	var wlent, rlent uint64 = 0, 0

	ChangeCryKey_Total := 0

	defer func() {
		fmt.Println("gorou exit.")
		fmt.Printf("\tpack  trlen:%d  twlen:%d\n", rlent, wlent)
		//fmt.Println("goPackTotal:", goPackTotal)
		fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

		fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
		fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
	}()

	if true {
		buf := apack.IsTheVersionConsistent()
		//tmr.Boot()
		//ChangeCryKey_Total += 1
		//outf2.Write(buf)
		for {
			if len(buf) > 0 {
				wlen, err := dst.Write(buf)
				wlent = wlent + uint64(wlen)
				if wlen == 0 {
					return
				}
				if err != nil && wlen <= 0 {
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

	if true {
		buf := apack.ChangeCryKey()
		//tmr.Boot()
		ChangeCryKey_Total += 1
		//outf2.Write(buf)
		for {
			if len(buf) > 0 {
				wlen, err := dst.Write(buf)
				wlent = wlent + uint64(wlen)
				if wlen == 0 {
					return
				}
				if err != nil && wlen <= 0 {
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

	for {

		recot_p_r.Run()
		rlen, err := src.Read(buf)
		recot_p_r.Run()

		rlent = rlent + uint64(rlen)

		if tmr.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
			continue
		}

		//outf.Write(buf[:rlen])
		tmr.Boot()
		rbuf = buf
		buf = buf[:rlen]

		wbuf.Reset()
		wbuf.Write(buf)
		//wbuf = append(wbuf, buf...)
		//fre := bool(len(wbuf) > 0)

		if wbuf.Len() > 0 {
			buf = apack.Packing(wbuf.Bytes())
			//wbuf = wbuf[len(wbuf):]
			//outf2.Write(buf)
			for {
				if len(buf) > 0 {
					recot_p_w.Run()
					wlen, err := dst.Write(buf)
					recot_p_w.Run()

					wlent = wlent + uint64(wlen)
					if wlen == 0 {
						return
					}
					if err != nil && wlen <= 0 {
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
			buf := apack.ChangeCryKey()
			ChangeCryKey_Total += 1
			tmr.Boot()
			//outf2.Write(buf)
			for {
				if len(buf) > 0 {
					wlen, err := dst.Write(buf)
					wlent = wlent + uint64(wlen)
					if wlen == 0 {
						return
					}
					if err != nil && wlen <= 0 {
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
		buf = rbuf
		if tmrP2.Run() {
			fmt.Printf("pack  trlen:%d  twlen:%d\n", rlent, wlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}

	}
}

func srcTOdstP_mt(src net.Conn, dst net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
			err := fmt.Errorf("Error:%s", x)
			fmt.Println(err)
			Logger.Println(x, "  Go exit.")
		}
	}()

	tmr_out := timerm.CreateTimer(networkTimeout)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	tmr_changekey := timerm.CreateTimer(tmr_changekey_time)

	recot_p_r := timerm.CreateRecoTime()

	apack := gstunnellib.CreateAesPack(key)

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
	buf := make([]byte, net_read_size)
	//var rbuf []byte

	var rlent, wlent uint64 = 0, 0

	ChangeCryKey_Total := 0

	dst_chan := make(chan ([]byte), netPUn_chan_cache_size)
	defer close(dst_chan)

	dst_ok := gstunnellib.CreateGorouStatus()

	defer func() {
		fmt.Println("\tgorou exit.")
		fmt.Printf("\t\tpack  trlen:%d\n", rlent)
		//fmt.Println("goPackTotal:", goPackTotal)
		fmt.Println("\tChangeCryKey_total:", ChangeCryKey_Total)

		fmt.Println("\tRecoTime_p_r All: ", recot_p_r.StringAll())
	}()

	if true {
		buf := apack.IsTheVersionConsistent()
		//tmr.Boot()
		//ChangeCryKey_Total += 1
		//outf2.Write(buf)
		for {
			if len(buf) > 0 {
				wlen, err := dst.Write(buf)
				wlent = wlent + uint64(wlen)
				if wlen == 0 {
					return
				}
				if err != nil && wlen <= 0 {
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

	if true {
		buf := apack.ChangeCryKey()
		//tmr.Boot()
		ChangeCryKey_Total += 1
		//outf2.Write(buf)
		for {
			if len(buf) > 0 {
				wlen, err := dst.Write(buf)
				wlent = wlent + uint64(wlen)
				if wlen == 0 {
					return
				}
				if err != nil && wlen <= 0 {
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

	go srcTOdstP_w(dst, dst_chan, dst_ok, wlent)
	dst = nil

	for dst_ok.IsOk() {

		recot_p_r.Run()
		rlen, err := src.Read(buf)
		recot_p_r.Run()

		rlent = rlent + uint64(rlen)

		if tmr_out.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
			continue
		}

		//outf.Write(buf[:rlen])
		tmr_out.Boot()
		buf = buf[:rlen]

		if len(buf) > 0 {
			buf = apack.Packing(buf)
		}
		dst_chan <- buf

		buf = make([]byte, net_read_size)

		if !dst_ok.IsOk() {
			return
		}

		if tmr_changekey.Run() {
			buf := apack.ChangeCryKey()
			ChangeCryKey_Total += 1
			tmr_out.Boot()
			//outf2.Write(buf)

			dst_chan <- buf
		}

		if tmrP2.Run() {
			fmt.Printf("pack  trlen:%d\n", rlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
		}

	}
}

func srcTOdstP_w(dst net.Conn, dst_chan chan ([]byte), dst_ok *gstunnellib.Gorou_status, wlentotal uint64) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
			err := fmt.Errorf("Error:%s", x)
			fmt.Println(err)
			Logger.Println(x, "  Go exit.")
		}
	}()

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	recot_p_w := timerm.CreateRecoTime()

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

	var wlent uint64 = wlentotal

	ChangeCryKey_Total := 0

	defer func() {
		fmt.Println("gorou exit.")
		fmt.Printf("\tpack  twlen:%d\n", wlent)
		//fmt.Println("goPackTotal:", goPackTotal)
		fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

		fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
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
			return
		}

		//wbuf.Reset()
		//wbuf.Write(buf)
		//wbuf = append(wbuf, buf...)
		//fre := bool(len(wbuf) > 0)

		if len(buf) > 0 {

			//wbuf = wbuf[len(wbuf):]
			//outf2.Write(buf)
			for {
				if len(buf) > 0 {
					recot_p_w.Run()
					wlen, err := dst.Write(buf)
					recot_p_w.Run()

					wlent = wlent + uint64(wlen)

					tmr_out.Boot()
					if wlen == 0 {
						return
					}
					if err != nil && wlen <= 0 {
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

		if tmrP2.Run() {
			fmt.Printf("pack twlen:%d\n", wlent)
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}
		if tmr_out.Run() {
			return
		}

	}
}

func srcTOdstP(src net.Conn, dst net.Conn) {
	srcTOdstP_mt(src, dst)
}

func srcTOdstUn_st(src net.Conn, dst net.Conn) {
	/*
	   defer func() {
	       if x := recover(); x != nil {
	           fmt.Println("Go exit.")
	       }
	   }()*/
	tmr := timerm.CreateTimer(time.Second * 60)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	recot_un_r := timerm.CreateRecoTime()
	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.CreateAesPack(key)

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
	buf := make([]byte, net_read_size)
	var rbuf, wbuf []byte
	var wlent, rlent uint64 = 0, 0

	defer func() {
		fmt.Println("gorou exit.")
		fmt.Printf("\tunpack trlen:%d  twlen:%d\n", rlent, wlent)

		fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
	}()

	for {
		recot_un_r.Run()
		rlen, err := src.Read(buf)
		recot_un_r.Run()

		rlent = rlent + uint64(rlen)

		if tmr.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
			continue
		}

		//outf.Write(buf[:rlen])
		tmr.Boot()
		rbuf = buf
		buf = buf[:rlen]
		wbuf = append(wbuf, buf...)
		for {
			ix, fre := find0(wbuf)
			p(ix, fre)
			if fre {
				buf = wbuf[:ix+1]
				wbuf = wbuf[ix+1:]
				pf("buf b:%d\n", len(buf))
				buf, err = apack.Unpack(buf)
				if err != nil {
					fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
					Logger.Println(err.Error())
					return
				}
				pf("buf a:%d\n", len(buf))
				//outf2.Write(buf)
				for {
					if len(buf) > 0 {
						recot_un_w.Run()
						wlen, err := dst.Write(buf)
						recot_un_w.Run()

						wlent = wlent + uint64(wlen)
						pf("twlen:%d  wlen:%d\n", wlent, wlen)
						if wlen == 0 {
							return
						}
						if err != nil && wlen <= 0 {
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
		if tmrP2.Run() {
			fmt.Printf("unpack trlen:%d  twlen:%d\n", rlent, wlent)

			fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}
}

func srcTOdstUn_mt(src net.Conn, dst net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
			Logger.Println(x, "  Go exit.")
		}
	}()

	tmr_out := timerm.CreateTimer(networkTimeout)
	//tmrP := timerm.CreateTimer(tmr_display_time)
	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

	recot_un_r := timerm.CreateRecoTime()

	//apack := gstunnellib.CreateAesPack(key)

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
	buf := make([]byte, net_read_size)
	//var rbuf []byte
	//var wbuff bytes.Buffer
	rlent := int64(0)

	dst_chan := make(chan ([]byte), netPUn_chan_cache_size)
	defer close(dst_chan)

	dst_ok := gstunnellib.CreateGorouStatus()

	go srcTOdstUn_w(dst, dst_chan, dst_ok)
	dst = nil

	defer func() {
		fmt.Println("\tgorou exit.")
		fmt.Printf("\t\tunpack  trlen:%d\n", rlent)
		//fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

		fmt.Println("\tRecoTime_un_r All: ", recot_un_r.StringAll())
	}()

	for dst_ok.IsOk() {
		recot_un_r.Run()
		rlen, err := src.Read(buf)
		recot_un_r.Run()
		//recot_un_r.RunDisplay("---recot_un_r:")
		//fmt.Println("---rlen:", rlen)

		pf("trlen:%d  rlen:%d\n", rlent, rlen)

		rlent = rlent + int64(rlen)
		//rlen = 0

		if tmr_out.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
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

		buf = make([]byte, net_read_size)

		if tmrP2.Run() {
			fmt.Printf("unpack  trlen:%d\n", rlent)
			//fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		}
	}
}

func srcTOdstUn_w(dst net.Conn, dst_chan chan ([]byte), dst_ok *gstunnellib.Gorou_status) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
			Logger.Println(x, "  Go exit.")
		}
	}()

	tmr_out := timerm.CreateTimer(networkTimeout)

	tmrP2 := timerm.CreateTimer(tmr_display_time)

	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.CreateAesPack(key)

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

	defer dst.Close()
	//defer outf.Close()
	//defer outf2.Close()

	//buf := make([]byte, net_read_size)
	var wbuf []byte
	//var wbuff bytes.Buffer
	var wlent uint64 = 0

	defer func() {
		fmt.Println("\tgorou exit.")
		fmt.Printf("\t\tunpack  twlen:%d\n", wlent)
		//fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

		fmt.Println("\tRecoTime_un_w All: ", recot_un_w.StringAll())
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

		//wbuff.Reset()
		//wbuff.Write(buf)

		//wbuf = wbuff.Bytes()

		buf, ok := <-dst_chan
		if !ok {
			return
		}

		wbuf = append(wbuf, buf...)
		for {
			ix, fre := find0(wbuf)
			pf("ix: %d, fre: %t\n", ix, fre)
			if fre {
				buf = wbuf[:ix+1]
				wbuf = wbuf[ix+1:]
				pf("buf not unpack:%d\n", len(buf))
				buf, err = apack.Unpack(buf)
				if err != nil {
					fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
					Logger.Println(err.Error())
					return
				}
				pf("buf unpack:%d\n", len(buf))
				//outf2.Write(buf)
				for {
					if len(buf) > 0 {
						recot_un_w.Run()
						wlen, err := dst.Write(buf)
						recot_un_w.Run()
						//recot_un_w.RunDisplay("---recot_un_w:")
						//fmt.Println("---wlen:", wlen)

						wlent = wlent + uint64(wlen)
						//wlen = 0
						tmr_out.Boot()

						pf("twlen:%d  wlen:%d\n", wlent, wlen)
						if wlen == 0 {
							return
						}
						if err != nil && wlen <= 0 {
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
			return
		}

		if tmrP2.Run() {
			fmt.Printf("unpack  twlen:%d\n", wlent)

			fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
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
	gstunnellib.CheckError(err)
}
