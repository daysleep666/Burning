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
	fmt.Println(ip)
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
		// fmt.Printf("\033[%dA\r\033[%dC", 1, 0)
		fmt.Printf("\033[1A\033[999D\033[K")
	}
	var w chan int
	<-w
}
