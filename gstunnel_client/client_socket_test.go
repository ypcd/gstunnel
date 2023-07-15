package main

import "testing"

func Test_client_socket(t *testing.T) {
	inTest_client_socket(t, false)
}

func Test_client_socket_mtg(t *testing.T) {
	inTest_client_socket(t, true)
}

func Test_client_socket_mt(t *testing.T) {
	inTest_client_socket_mt(t, false)
}

func Test_client_socket_mt_mtg(t *testing.T) {
	inTest_client_socket_mt(t, true)
}