package main

import (
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

/*
func NewLocalNode() (*enode.LocalNode, *enode.DB) {
	db, _ := enode.OpenDB("")
	//生成本地节点的密钥
	key, _ := crypto.GenerateKey()
	return enode.NewLocalNode(db, key), db
}
*/

//func NewTransport(remotePk *ecdsa.PublicKey,fd net.Conn)transport{
//	wrapped :=
//}
func main() {
	db, _ := enode.OpenDB("")
	//生成本地节点的密钥
	key, _ := crypto.GenerateKey()
	localNode := enode.NewLocalNode(db, key)
	//本地的ip
	ip := net.ParseIP("192.168.240.129")
	fmt.Println(ip)

	localNode.SetStaticIP(ip)
	defer db.Close()

	//设置端口
	//TCP 连接的本地接口
	x := enr.TCP(8080)
	localNode.Set(enr.WithEntry("tcp", x))

	if err := localNode.Node().Load(enr.WithEntry("tcp", &x)); err != nil {
		fmt.Println("不能加载输入 'tcp':", err)
	} else if x != 3 {
		fmt.Println("错误的输入:", x)
	}

	fmt.Println("直接生成的节点信息:", localNode.Node().String())
	fmt.Println("生成的本地节点的IP:", localNode.Node().IP())
	fmt.Println("生成的本地节点的TCP端口:", localNode.Node().TCP())
	fmt.Println("生成的本地节点的ID(64位公钥):", localNode.Node().ID())

	//由公钥,ip,tcp端口生成节点
	node1 := enode.NewV4(&key.PublicKey, ip, 8080, 0)
	fmt.Println("直接生成的节点信息:", node1.String())

	node2, _ := enode.ParseV4(node1.String())

}
