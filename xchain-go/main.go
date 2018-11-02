package main

const dbFile = "xchain.db"

var NodeArr = make([]Node, 10)

func main() {

	cli := CLI{}
	cli.Run()

}
