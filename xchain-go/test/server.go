package main

import (
	"fmt"
	"net"
	"os"
	"xchain-go/rlp"
	//	"time"
)

func main() {
	fmt.Println("UdpListen Start")
	packetConn, err := net.ListenPacket("udp", ":8333")

	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
		return
	}
	defer packetConn.Close()

	var buf [512]byte
	for {
		n, addr, err := packetConn.ReadFrom(buf[0:])
		fmt.Println("buf:", buf[0:n])
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
			return
		}
		fmt.Fprintf(os.Stdout, "listen recv: %s\n", string(buf[0:n]))

		// 将数组反序列化
		var cc string
		err = rlp.DecodeBytes(buf[0:n], &cc)
		//fmt.Println("number0:", blocks[0].header.Number)
		if err != nil {
			fmt.Println("decodeerr", err)
		}
		fmt.Printf("######blocks######=%v\n", cc)

		_, err = packetConn.WriteTo(buf[0:n], addr)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
			return
		}
	}
}
