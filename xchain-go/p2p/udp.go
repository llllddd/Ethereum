package p2p

import (
	"fmt"
	"net"
	"os"
	"time"
	"xchain-go/core/basic"
	"xchain-go/rlp"
)

var knowNodes = []string{"127.0.0.1", "192.168.82.156", "192.168.82.72", "192.168.82.87"}

/**
  UDP广播节点信息
*/
func UdpDial() {
	//通过广播指定端口实现对目标节点的广播
	conn, err := net.Dial("udp", "127.0.0.1:8001")
	fmt.Println("blocks:", "blockchain")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
		return
	}

	defer conn.Close()
	//序列化，通过rlp编码成字节类型进行传输
	b, err := rlp.EncodeToBytes("blockchain")
	fmt.Println("After Encode:", b)

	var bb []*basic.Block
	err = rlp.DecodeBytes(b, &bb)
	fmt.Printf("********bb********%v\n", bb)
	// fmt.Println("After bb number:", bb[0].header.Number)

	_, err = conn.Write(b) //[]byte("hello UDP")
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
		return
	}

	var buf [512]byte
	// 阻塞，直到接收到消息,设置阻塞时间1秒
	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
		return
	}

	fmt.Fprintf(os.Stdout, "dial recv: %s\n", buf[0:n])
}

/**
  UDP监听服务
*/
func UdpListen() {

	fmt.Println("UdpListen Start")
	packetConn, err := net.ListenPacket("udp", ":8001")

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
		var blocks []basic.Block
		err = rlp.DecodeBytes(buf[0:n], &blocks)
		//fmt.Println("number0:", blocks[0].header.Number)
		if err != nil {
			fmt.Println("decodeerr", err)
		}
		fmt.Printf("######blocks######=%v\n", blocks)

		_, err = packetConn.WriteTo(buf[0:n], addr)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error: %s", err.Error())
			return
		}
	}

	// if err != nil {
	// 	fmt.Fprintf(os.Stdout, "Listen Error: %s\n", err.Error())
	// 	return
	// }
	// defer packetConn.Close()

	// var buf []byte
	// for {
	// 	n, addr, err := packetConn.ReadFrom(buf[0:])
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stdout, "Error: %s\n", err.Error())
	// 		return
	// 	}

	// 	//fmt.Printf("buf[0:%v] = %v\n", n, buf[0:n])
	// 	fmt.Fprintf(os.Stdout, "listen recv: %s\n", string(buf[0:n]))
	// 	//2.将数组反序列化
	// 	// var intarray []string
	// 	// err = rlp.DecodeBytes(buf[0:n], &intarray)
	// 	// fmt.Printf("######intarray######=%v\n", intarray)

	// 	_, err = packetConn.WriteTo(buf[0:n], addr)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stdout, "Error: %s\n", err.Error())
	// 		return
	// 	}
	// }

}
