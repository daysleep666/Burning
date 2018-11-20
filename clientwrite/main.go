package main

import (
	"fmt"
	"net"
	"os"

	"github.com/daysleep666/Burning/tool"
)

func main() {
	ip := os.Args[1]
	conn, err := net.Dial("tcp", ip)
	defer conn.Close()
	tool.CheckErr(err)

	for {
		var content string
		fmt.Scanln(&content)
		conn.Write([]byte(content))
		fmt.Printf("\033[33m\033[1A\033[999D\033[K")
	}
	var w chan int
	<-w
}
