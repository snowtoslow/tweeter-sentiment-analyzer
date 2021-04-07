package main

import (
	"log"
	"net"
)

func main() {
	var conn net.Conn
	srv := NewServer(conn)
	log.Fatal(srv.RunServer(":8088"))
}
