package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/1681-dodo-bird/nettest/msg"
)

/*
受信ごるーちん
*/
func recv(conn *net.UDPConn, ch chan time.Duration) {
	buf := make([]byte, 1024)
	m := msg.Message{}
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("error at Read.%s\n", err)
			continue
		}
		binary.Read(bytes.NewReader(buf[:n]), binary.BigEndian, &m)
		now := time.Now()
		startAt := time.Unix(m.StartAtSec, m.StartAtNSec)
		tat := now.Sub(startAt)
		ch <- tat
	}
}

/*
統計情報を表示するゴルーチン
*/
func printer(ch chan time.Duration) {

	i := 0
	var sum, min, max time.Duration
	min = time.Minute
	max = 0 * time.Millisecond
	for {
		v := <-ch
		if v > time.Minute {
			continue
		}
		if max < v {
			if v > time.Minute {
				continue
			}
			max = v
		}
		if min > v {
			min = v
		}
		sum += v
		i++
		if i >= 100 {
			fmt.Printf("avg: %v, max: %v, min :%v\n", sum/100, max, min)
			i = 0
			sum = 0
			min = time.Minute
			max = 0 * time.Millisecond
		}

	}
}

func main() {

	var addr string

	flag.StringVar(&addr, "address", "127.0.0.1:12299", "IP:Port")
	flag.Parse()

	// UDP準備
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("Cannot listen:%s\n", err)
		return
	}
	clientAddr, err := net.ResolveUDPAddr("udp", ":")
	if err != nil {
		fmt.Printf("Cannot listen:%s\n", err)
		return
	}
	conn, err := net.DialUDP("udp", clientAddr, serverAddr)
	if err != nil {
		fmt.Printf("Cannot create conn:%s\n", err)
	}
	defer conn.Close()

	// 受信、統計ゴルーチン開始
	ch := make(chan time.Duration)
	go recv(conn, ch)
	go printer(ch)

	// 送信処理
	buf := &bytes.Buffer{}
	m := msg.Message{}
	var now time.Time
	var i uint64 = 0
	for {
		// バッファ初期化
		buf.Reset()

		// バッファに書き込み
		now = time.Now()
		m.Index = i
		i++
		m.StartAtSec = now.Unix()
		m.StartAtNSec = now.UnixNano() % (1000 * 1000 * 1000)
		err = binary.Write(buf, binary.BigEndian, &m)
		if err != nil {
			fmt.Printf("Cannot create conn:%s\n", err)
		}

		// 送信
		conn.Write(buf.Bytes())

		// 休み
		time.Sleep(10 * time.Millisecond)
	}
}
