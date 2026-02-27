package main

import "github.com/rouzbehsbz/manticore/server/pkg/network"

func main() {
	server := network.NewServer()

	server.Listen("0.0.0.0:3000")
}
