package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/daysleep666/Burning/tool"
)

func main() {
	ip := os.Args[1]
	conn, err := net.Dial("tcp", ip)
	defer conn.Close()
	tool.CheckErr(err)
	go func() {
		for {
			bs := make([]byte, 10000)
			_, err := conn.Read(bs)
			if err != nil {
				return
			}
			userContent := string(bs)
			fmt.Printf("\033[34m%v", string(userContent))
			time.Sleep(time.Second * 1)
			fmt.Printf("\033[999D\033[K")

		}
	}()

	for {
		var content string
		fmt.Scanln(&content)
		conn.Write([]byte(content))
		fmt.Printf("\033[33m\033[1A\033[999D\033[K")
	}

	conn.Write([]byte("reader"))
	var w chan int
	<-w
}
