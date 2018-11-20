package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ip := os.Args[1]
	fmt.Println(ip)
	conn, err := net.Dial("tcp", ip)
	defer conn.Close()
	CheckErr(err)

	for {
		var content string
		fmt.Scanln(&content)
		conn.Write([]byte(content))
		fmt.Printf("\033[%dA\r\033[%dC\033[K", 1, 0)
	}
	var w chan int
	<-w
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
