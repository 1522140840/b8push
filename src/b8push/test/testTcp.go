package main

import (
	"net"
	"time"
	"encoding/binary"
	"io"
	"fmt"
)

func main()  {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:"+"9999")



	l, _ := net.ListenTCP("tcp", tcpAddr)



	handleAcpt(l)
}
func handleAcpt(l *net.TCPListener) {
	for {
		c, _ := l.AcceptTCP()

		go handleConn(c)
	}
}

func handleConn(conn *net.TCPConn){
	for {

		//设置两倍的心跳时间
		conn.SetDeadline(time.Now().Add(10000 * time.Second))

		var length int16

		binary.Read(conn, binary.BigEndian, &length)



		packetBytes := make([]byte, length)
		io.ReadFull(conn, packetBytes)

		fmt.Println(packetBytes)

		break;
	}
}
