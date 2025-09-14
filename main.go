package main

import (
	reader "zakaria/mist-vpn/client/client-reader"
	"zakaria/mist-vpn/client/connection"
)

func main() {

	go connection.SendDummyContent()
	connection.InitClient()

	reader.ReadPacketsFromTun0()

}
