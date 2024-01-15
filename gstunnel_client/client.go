/*
*
*Open source agreement:
*   The project is based on the GPLv3 protocol.
*
 */
package main

func main() {
	init_client_run()
	gstClient := newGstClient(g_gsconfig.Listen, g_gsconfig.GetServer_rand())
	defer gstClient.close()
	for {
		gstClient.run()
	}
}

/*
func createToGstServerConnHandler(chanconnlist <-chan net.Conn, gstServerAddr string) {

	const maxTryConnService int = 5000
	//var err error
	for {

		rawClient, ok := <-chanconnlist
		if !ok {
			checkError_info(errors.New("chanconnlist is closed"))
			return
		}
		g_log_List.GSIpLogger.Printf("Raw client ip: %s\n", rawClient.RemoteAddr().String())

		server_conn_error_total := 0
		tmr := timerm.CreateTimer(time.Second * 10)
		for {
			if tmr.Run() {
				g_logger.Println("Error: server_conn_error_total: ", server_conn_error_total)
				tmr.Boot()
			}

			gstServer, err := net.Dial("tcp", gstServerAddr)
			//checkError(err)
			if err != nil {
				if server_conn_error_total > 1000 {
					g_logger.Fatalln("Error: server_conn_error_total > 1000")
				}
				server_conn_error_total++
				g_logger.Println("Error: [net.Dial('tcp', service)]:", err.Error())

				continue
			}
			//g_logger.Println("conn.", service)

			//rawClient: 		src---client
			//dst: 		client---serever
			//pack: 	rawClient read recv, dst wirte send.
			//unpack:	dst read recv, rawClient wirte send.

			gctx := gstunnellib.NewGSContextImp(g_gid.GenerateId(), g_gstst)
			g_gstst.GetStatusConnList().Add(gctx.GetGsId(), rawClient, gstServer)

			go srcTOdstP_count(rawClient, gstServer, gctx)
			go srcTOdstUn_count(gstServer, rawClient, gctx)
			g_logger.Printf("go [%d].\n", gctx.GetGsId())
			break
		}
	}
}


	func run() {
		//defer gstunnellib.Panic_Recover_GSCtx(g_logger, gctx)

		var lstnaddr string
		var connaddr []string

		lstnaddr = g_gsconfig.Listen
		connaddr = g_gsconfig.GetServers()

		g_logger.Println("Listen_Addr:", lstnaddr)
		g_logger.Println("Conn_Addr:", connaddr)
		g_logger.Println("Begin......")

		service := lstnaddr
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		checkError(err)
		listener, err := net.ListenTCP("tcp", tcpAddr)
		checkError(err)

		chanConnList := make(chan net.Conn, 10240)
		defer close(chanConnList)

		for i := 0; i < g_connHandleGoNum; i++ {
			go createToGstServerConnHandler(chanConnList, connaddr[0])
		}
		for {
			rawClient, err := listener.Accept()
			if err != nil {
				g_logger.Println("Error:", err)
				continue
			}
			chanConnList <- rawClient
		}
	}
*/
/*
func VersionPack_send(obj *gsobj.GSTObj) error {
	return gsobj.VersionPack_send(obj)
}

func ChangeCryKey_send(obj *gsobj.GSTObj) error {
	return gsobj.ChangeCryKey_send_fromClient(obj)
}
*/
