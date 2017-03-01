/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
*/
package main

import (
    "fmt"
    "net"
    "os"
    "gstunnellib"
    "time"
    "timerm"
)

var p = gstunnellib.Nullprint
var pf = gstunnellib.Nullprintf

var fpnull = os.DevNull

var key string

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
    lstnaddr := os.Args[1]
    connaddr := os.Args[2]
    key = os.Args[3]
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
        
        go srcTOdstP( acc, dst)
        go srcTOdstUn( dst, acc)
        fmt.Println("go.")
    }
}

func find0(v1 []byte) (int, bool) {
    for i := 0; i<len(v1); i++ {
        if v1[i]==0 {
            return i, true
        }
    }
    return -1, false
}

func srcTOdstP(src net.Conn, dst net.Conn){
    defer func() {
        if x := recover(); x != nil {
            fmt.Println("Go exit.")
        }
    }()

    tmr := timerm.CreateTimer(time.Second*60)
    tmrP := timerm.CreateTimer(time.Second*1)
    tmrP2 := timerm.CreateTimer(time.Second*1)
    
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
        if err != nil{
            continue
        }
        
        outf.Write( buf[:rlen])
        tmr.Boot()
        rbuf = buf
        buf = buf[:rlen]
        wbuf = append(wbuf, buf...)
        fre := bool(len(wbuf)>0)
        if fre {
            buf = apack.Packing(wbuf)
            wbuf = wbuf[len(wbuf):]
            outf2.Write(buf)
            for{
                if len(buf)>0 {                
                    wlen, err := dst.Write(buf)
                    wlent = wlent + wlen
                    if wlen == 0 {
                        return
                    }
                    if err != nil && wlen <= 0{
                        continue
                    }
                    if len(buf)==wlen {
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
        }

    }
}

func srcTOdstUn(src net.Conn, dst net.Conn){
    defer func() {
        if x := recover(); x != nil {
            fmt.Println("Go exit.")
        }
    }()

    tmr := timerm.CreateTimer(time.Second*60)
    tmrP := timerm.CreateTimer(time.Second*1)
    tmrP2 := timerm.CreateTimer(time.Second*1)

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
            x1 :=1
            x1++
        }

        if tmr.Run() {
            return
        }
        if rlen == 0 {
            return
        }        
        if err != nil{
            continue
        }
        
        outf.Write( buf[:rlen])
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
                for{
                    if len(buf)>0 {                
                        wlen, err := dst.Write(buf)
                        wlent = wlent + wlen
                        pf("twlen:%d  wlen:%d\n", wlent, wlen)
                        if wlen == 0 {
                            return
                        }
                        if err != nil && wlen <= 0{
                            continue
                        }
                        if len(buf)==wlen {
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
        }
    }
}


func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}