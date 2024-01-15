/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

func main() {
	init_server_run()

	gstServer := newGstServer(g_gsconfig.Listen, g_gsconfig.GetServer_rand())
	defer gstServer.close()
	for {
		gstServer.run()
	}
}

/*
func run() {
	//defer gstunnellib.Panic_Recover_GSCtx(g_logger, gctx)

	var listenAddr, rawServiceAddr string

	listenAddr = g_gsconfig.Listen
	rawServiceAddr = g_gsconfig.GetServer_rand()

	g_logger.Println("Listen_Addr:", listenAddr)
	g_logger.Println("Conn_Addr:", rawServiceAddr)
	g_logger.Println("Begin......")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", listenAddr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	//var rawService net.Conn

	chanConnList := make(chan net.Conn, 10240)
	defer close(chanConnList)

	for i := 0; i < g_connHandleGoNum; i++ {
		go createToRawServiceConnHandler(chanConnList, rawServiceAddr)
	}

	for {
		gstClient, err := listener.Accept()
		if err != nil {
			checkError_NoExit(err)
			continue
		}
		chanConnList <- gstClient
	}
}
*/
/*
func VersionPack_send(dst net.Conn, apack gstunnellib.IGSPack, wlent *int64) error {
	return gstunnellib.VersionPack_send(dst, apack, wlent)
}

func ChangeCryKey_send(dst net.Conn, apack gstunnellib.IGSRSAPackNet, ChangeCryKey_Total *int, wlent *int64) error {
	return gstunnellib.ChangeCryKey_send_fromServer(dst, apack, ChangeCryKey_Total, wlent)
}
*/
