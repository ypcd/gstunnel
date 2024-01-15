package main

import (
	"testing"
)

func Test_client_socket(t *testing.T) {
	inTest_client_socket(t, false, 200)
	//time.Sleep(time.Second * 1000)
	logger_test.Println("------Test done.------")
}

func Test_client_socket_mtg(t *testing.T) {
	inTest_client_socket(t, true, 200)
	logger_test.Println("------Test done.------")
}

func Test_client_socket_mt(t *testing.T) {
	inTest_client_socket_mt(t, false)
	logger_test.Println("------Test done.------")
}

func Test_client_socket_mt_mtg(t *testing.T) {
	inTest_client_socket_mt(t, true)
	logger_test.Println("------Test done.------")
}
