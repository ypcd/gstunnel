/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	//	"bytes"

	"bytes"
	"fmt"
	"gstunnellib"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	//	"runtime"
	//	"runtime/pprof"
	"sync/atomic"
	"time"
	"timerm"
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

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func init() {
	if debug_client {
		//pf = fmt.Printf
	}

	Logger = gstunnellib.CreateFileLogger("gstunnel_client.err.log")

	fmt.Println("VER:", version)
	Logger.Println("VER:", version)

	gsconfig = gstunnellib.CreateGsconfig("config.client.json")
	debug_client = gsconfig.Debug

	fmt.Println("debug:", debug_client)
	Logger.Println("debug:", debug_client)

	if debug_client {
		/*
			go func() {

				mux := http.NewServeMux()
				mux.HandleFunc("/custom_debug_path/profile", pprof.Profile)
				log.Fatal(http.ListenAndServe("127.0.0.1:7777", mux))

			}()
		*/
		go func() {
			Logger.Fatalln("http server: ", http.ListenAndServe("localhost:6060", nil))
		}()
		fmt.Println("Debug server listen: localhost:6060")
		Logger.Println("Debug server listen: localhost:6060")
	}

}

func main() {

	/*
		go func() {
			log.Println("http server: ", http.ListenAndServe("localhost:6060", nil))
		}()
	*/
	/*
		flag.Parse()
		//fmt.Println
		scpu := string("cpu.prof")
		cpuprofile = &scpu
		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			fmt.Println("cpufile is create.")
			defer f.Close() // error handling omitted for example
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			fmt.Println("cpufile start.")
			defer pprof.StopCPUProfile()
		}

		// ... rest of the program ...

		if *memprofile != "" {
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			defer f.Close() // error handling omitted for example
			runtime.GC()    // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}
	*/
	for {
		run()
	}

	//fmt.Println("test...")

	//test1()
	//	a1 := 123
	//	_ = a1
	//return
}

func run() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Error:", x)
			Logger.Println(x, "  App restart.")
			fmt.Println("App restart.")

		}
	}()

	var lstnaddr, connaddr string

	if len(os.Args) == 4 {
		lstnaddr = os.Args[1]
		connaddr = os.Args[2]
		key = os.Args[3]
	} else {
		lstnaddr = gsconfig.Listen
		connaddr = gsconfig.Server
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
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		fmt.Println(tcpAddr)
		checkError(err)
		dst, err := net.Dial("tcp", service)
		checkError(err)
		fmt.Println("conn.")

		go srcTOdstP_count(acc, dst)
		go srcTOdstUn_count(dst, acc)

		//time.Sleep(time.Second * 3)
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

func srcTOdstP(src net.Conn, dst net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
			err := fmt.Errorf("Error:%s", x)
			fmt.Println(err)
			Logger.Println(x, "  Go exit.")

		}
	}()

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP := timerm.CreateTimer(time.Second * 1)
	tmrP2 := timerm.CreateTimer(time.Second * 1)

	tmr_changekey := timerm.CreateTimer(time.Minute * 1)

	recot_p_r := timerm.CreateRecoTime()
	recot_p_w := timerm.CreateRecoTime()

	apack := gstunnellib.CreateAesPack(key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull

	outf, err := os.Create(fp1)
	checkError(err)

	outf2, err := os.Create(fp2)
	checkError(err)
	defer src.Close()
	defer dst.Close()
	defer outf.Close()
	buf := make([]byte, net_read_size)
	var rbuf []byte
	var wbuf bytes.Buffer

	wlent, rlent := 0, 0

	ChangeCryKey_Total := 0

	defer func() {
		fmt.Printf("pack  trlen:%d  twlen:%d\n", rlent, wlent)
		fmt.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))
		fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

		fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
		fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())
	}()

	if true {
		buf := apack.IsTheVersionConsistent()
		//tmr.Boot()
		//ChangeCryKey_Total += 1
		outf2.Write(buf)
		for {
			if len(buf) > 0 {
				wlen, err := dst.Write(buf)
				wlent = wlent + wlen
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
		outf2.Write(buf)
		for {
			if len(buf) > 0 {
				wlen, err := dst.Write(buf)
				wlent = wlent + wlen
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
		//recot_p_r.RunDisplay("---recot_p_r:")
		//fmt.Println("---rlen:", rlen)

		rlent = rlent + rlen
		if tmrP.Run() {
			fmt.Fprintf(os.Stderr, "%d read end...\n", rlen)
		}

		if tmr_out.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
			continue
		}

		outf.Write(buf[:rlen])
		tmr_out.Boot()
		rbuf = buf
		buf = buf[:rlen]

		wbuf.Reset()
		wbuf.Write(buf)
		//wbuf = append(wbuf, buf...)
		// := bool(len(wbuf) > 0)
		if wbuf.Len() > 0 {
			buf = apack.Packing(wbuf.Bytes())
			//wbuf = wbuf[len(wbuf):]
			outf2.Write(buf)
			for {
				if len(buf) > 0 {
					recot_p_w.Run()
					wlen, err := dst.Write(buf)
					recot_p_w.Run()
					//recot_p_w.RunDisplay("---recot_p_w:")
					//fmt.Println("---wlen:", wlen)

					wlent = wlent + wlen
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
			tmr_out.Boot()
			outf2.Write(buf)
			for {
				if len(buf) > 0 {

					wlen, err := dst.Write(buf)

					wlent = wlent + wlen
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
			fmt.Println("goPackTotal:", atomic.LoadInt32(&goPackTotal))
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)

			fmt.Println("RecoTime_p_r All: ", recot_p_r.StringAll())
			fmt.Println("RecoTime_p_w All: ", recot_p_w.StringAll())

		}

	}
}

func srcTOdstUn(src net.Conn, dst net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
			Logger.Println(x, "  Go exit.")
		}
	}()

	tmr_out := timerm.CreateTimer(networkTimeout)
	tmrP := timerm.CreateTimer(time.Second * 1)
	tmrP2 := timerm.CreateTimer(time.Second * 1)

	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

	recot_un_r := timerm.CreateRecoTime()
	recot_un_w := timerm.CreateRecoTime()

	apack := gstunnellib.CreateAesPack(key)

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	outf, err := os.Create(fp1)
	checkError(err)

	outf2, err := os.Create(fp2)
	checkError(err)
	defer src.Close()
	defer dst.Close()
	defer outf.Close()
	defer outf2.Close()
	buf := make([]byte, net_read_size)
	var rbuf, wbuf []byte
	//var wbuff bytes.Buffer
	wlent, rlent := 0, 0

	defer func() {
		fmt.Printf("unpack  trlen:%d  twlen:%d\n", rlent, wlent)
		fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

		fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
		fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
	}()

	for {
		recot_un_r.Run()
		rlen, err := src.Read(buf)
		recot_un_r.Run()
		//recot_un_r.RunDisplay("---recot_un_r:")
		//fmt.Println("---rlen:", rlen)

		pf("trlen:%d  rlen:%d\n", rlent, rlen)

		rlent = rlent + rlen
		//rlen = 0

		if tmrP.Run() {
			fmt.Fprintf(os.Stderr, "%d read end...\n", rlen)
			x1 := 1
			x1++
		}

		if tmr_out.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
			continue
		}

		outf.Write(buf[:rlen])
		tmr_out.Boot()
		rbuf = buf
		buf = buf[:rlen]

		//wbuff.Reset()
		//wbuff.Write(buf)

		//wbuf = wbuff.Bytes()

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
					fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
					Logger.Println(err.Error())
					return
				}
				pf("buf unpack:%d\n", len(buf))
				outf2.Write(buf)
				for {
					if len(buf) > 0 {
						recot_un_w.Run()
						wlen, err := dst.Write(buf)
						recot_un_w.Run()
						//recot_un_w.RunDisplay("---recot_un_w:")
						//fmt.Println("---wlen:", wlen)

						wlent = wlent + wlen
						//wlen = 0

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
			fmt.Printf("unpack  trlen:%d  twlen:%d\n", rlent, wlent)
			fmt.Println("goUnpackTotal:", atomic.LoadInt32(&goUnpackTotal))

			fmt.Println("RecoTime_un_r All: ", recot_un_r.StringAll())
			fmt.Println("RecoTime_un_w All: ", recot_un_w.StringAll())
		}
	}
}

func checkError(err error) {
	gstunnellib.CheckError(err)
}
