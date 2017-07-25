package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.0.199:12299")
	if err != nil {
		fmt.Printf("Cannot listen:%s\n", err)
		return
	}
	clientAddr, err := net.ResolveUDPAddr("udp", ":33444")
	if err != nil {
		fmt.Printf("Cannot listen:%s\n", err)
		return
	}
	conn, err := net.DialUDP("udp", clientAddr, serverAddr)
	if err != nil {
		fmt.Printf("Cannot create conn:%s\n", err)
	}
	defer conn.Close()
	buf := &bytes.Buffer{}
	var i int64
	for {
		buf.Reset()
		_ = binary.Write(buf, binary.BigEndian, i)
		conn.Write(buf.Bytes())
		i++
		if i%100000 == 0 {
			fmt.Println(time.Now())
			i = 0
		}
	}
}
