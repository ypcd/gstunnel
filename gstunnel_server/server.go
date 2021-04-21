/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

import (
	"encoding/json"
	"fmt"
	"gstunnellib"
	"io/ioutil"
	"net"
	"os"
	"time"
	"timerm"
)

var version string = gstunnellib.Version

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string

var gsconfig gsConfig

type gsConfig struct {
	Listen string
	Server string
	Key    string
}

func CreateGsconfig() {
	f, err := os.Open("config.server.json")
	checkError(err)

	defer func() {
		f.Close()
	}()

	buf, err := ioutil.ReadAll(f)
	checkError(err)

	//fmt.Println(string(buf))

	json.Unmarshal(buf, &gsconfig)
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

	fmt.Println("VER:", version)
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

		go srcTOdstUn(acc, dst)
		go srcTOdstP(dst, acc)
		fmt.Println("go.")
	}
}

func find0(v1 []byte) (int, bool) {
	return gstunnellib.Find0(v1)
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
	checkError(err)
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
			//fmt.Println("goPackTotal:", goPackTotal)
			fmt.Println("ChangeCryKey_total:", ChangeCryKey_Total)
		}

	}
}

func srcTOdstUn(src net.Conn, dst net.Conn) {
	/*
	   defer func() {
	       if x := recover(); x != nil {
	           fmt.Println("Go exit.")
	       }
	   }()*/
	tmr := timerm.CreateTimer(time.Second * 60)
	tmrP := timerm.CreateTimer(time.Second * 1)
	tmrP2 := timerm.CreateTimer(time.Second * 1)

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
				buf, err = apack.Unpack(buf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
					return
				}
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
			fmt.Printf("trlen:%d  twlen:%d\n", rlent, wlent)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
