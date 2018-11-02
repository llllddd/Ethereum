package p2p

import (
	"fmt"
	"log"
	"net"
	"os"
)

/**
  TCP Client服务
*/
func TcpDial() {
	conn, err := net.Dial("tcp", "127.0.0.1:7070")
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, 512)

	n, err2 := conn.Read(buffer)
	if err2 != nil {
		fmt.Println("Read failed:", err2)
		return
	}

	fmt.Println("count:", n, "msg:", string(buffer))
}

/**
  TCP监听服务
*/
func TcpListen() {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":2002")

	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", ":2002")
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", ":2002", err)
	} else {

		log.Println("rpc listening", ":2002")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}
		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	var buffer []byte = []byte("You are welcome. I'm server.")
	n, err := conn.Write(buffer)

	if err != nil {
		fmt.Println("Write error:", err)
	}
	fmt.Println("send:", n)
	fmt.Println("connetion end")

}
