package main

import "testing"

func Test_server_socket(t *testing.T) {
	inTest_server_socket(t, false)
}

func Test_server_socket_mtg(t *testing.T) {
	inTest_server_socket(t, true)
}

func Test_server_socket_mt(t *testing.T) {
	inTest_server_socket_mt(t, false)
}

func Test_server_socket_mt_mtg(t *testing.T) {
	inTest_server_socket_mt(t, true)
}

func Test_server_socket_mt_old(t *testing.T) {
	inTest_server_socket_mt_old(t, false)
}
