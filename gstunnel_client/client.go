/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	//	"bytes"
	"encoding/json"
	"flag"
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

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string

var gsconfig gsConfig

var goPackTotal, goUnpackTotal int32 = 0, 0

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {

	go func() {
		log.Println("http server: ", http.ListenAndServe("localhost:6060", nil))
	}()
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

type gsConfig struct {
	Listen string
	Server string
	Key    string
}

func test1() {
	f, _ := os.Open("config.json")
	buf := make([]byte, 10000)
	n, _ := f.Read(buf)
	buf = buf[:n]
	fmt.Println(string(buf))

	gsc1 := gsConfig{}
	json.Unmarshal(buf, &gsc1)
}

func CreateGsconfig() {
	f, err := os.Open("config.json")
	checkError(err)

	defer func() {
		f.Close()
	}()

	buf := make([]byte, 10000)
	n, _ := f.Read(buf)
	buf = buf[:n]
	//fmt.Println(string(buf))

	//gsc1 := gsConfig{}
	json.Unmarshal(buf, &gsconfig)
}

func run() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Error:", x)
			fmt.Println("App restart.")

		}
	}()

	var lstnaddr, connaddr string

	if len(os.Args) == 4 {
		lstnaddr = os.Args[1]
		connaddr = os.Args[2]
		key = os.Args[3]
	} else {
		CreateGsconfig()
		lstnaddr = gsconfig.Listen
		connaddr = gsconfig.Server
		key = gsconfig.Key
	}
	fmt.Println(lstnaddr)
	fmt.Println(connaddr)
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

func find0_1(v1 []byte) (int, bool) {
	for i := 0; i < len(v1); i++ {
		if v1[i] == 0 {
			return i, true
		}
	}
	return -1, false
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
		}
	}()

	tmr := timerm.CreateTimer(time.Second * 60)
	tmrP := timerm.CreateTimer(time.Second * 1)
	tmrP2 := timerm.CreateTimer(time.Second * 1)

	tmr_changekey := timerm.CreateTimer(time.Minute * 1)

	apack := gstunnellib.CreateAesPack(key)

	fp1 := "CPrecv.data"
	fp2 := "CPsend.data"

	fp1, fp2 = fpnull, fpnull

	outf, err := os.Create(fp1)
	outf2, err := os.Create(fp2)
	checkError(err)
	defer src.Close()
	defer dst.Close()
	defer outf.Close()
	buf := make([]byte, 1024*64)
	var rbuf, wbuf []byte

	wlent, rlent := 0, 0

	ChangeCryKey_Total := 0

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

		rlen, err := src.Read(buf)
		rlent = rlent + rlen
		if tmrP.Run() {
			fmt.Fprintf(os.Stderr, "%d read end...", rlen)
		}

		if tmr.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
			continue
		}

		outf.Write(buf[:rlen])
		tmr.Boot()
		rbuf = buf
		buf = buf[:rlen]
		wbuf = append(wbuf, buf...)
		fre := bool(len(wbuf) > 0)
		if fre {
			buf = apack.Packing(wbuf)
			wbuf = wbuf[len(wbuf):]
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
		if tmr_changekey.Run() {
			buf := apack.ChangeCryKey()
			ChangeCryKey_Total += 1
			tmr.Boot()
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
			fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)
		}

	}
}

func srcTOdstUn(src net.Conn, dst net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Go exit.")
		}
	}()

	tmr := timerm.CreateTimer(time.Second * 60)
	tmrP := timerm.CreateTimer(time.Second * 1)
	tmrP2 := timerm.CreateTimer(time.Second * 1)

	//tmr_changekey := timerm.CreateTimer(time.Minute * 10)

	apack := gstunnellib.CreateAesPack(key)

	fp1 := "SUrecv.data"
	fp2 := "SUsend.data"
	fp1 = fpnull
	fp2 = fpnull
	outf, err := os.Create(fp1)
	outf2, err := os.Create(fp2)

	checkError(err)
	defer src.Close()
	defer dst.Close()
	defer outf.Close()
	defer outf2.Close()
	buf := make([]byte, 1024*64)
	var rbuf, wbuf []byte
	wlent, rlent := 0, 0

	for {

		rlen, err := src.Read(buf)
		rlent = rlent + rlen
		if tmrP.Run() {
			fmt.Fprintf(os.Stderr, "%d read end...", rlen)
			x1 := 1
			x1++
		}

		if tmr.Run() {
			return
		}
		if rlen == 0 {
			return
		}
		if err != nil {
			continue
		}

		outf.Write(buf[:rlen])
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
				buf = apack.Unpack(buf)
				pf("buf a:%d\n", len(buf))
				outf2.Write(buf)
				for {
					if len(buf) > 0 {
						wlen, err := dst.Write(buf)
						wlent = wlent + wlen
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
			fmt.Println("goUnpackTotal:", goUnpackTotal)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(-11)
	}
}
