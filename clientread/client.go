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
	var stopWrite bool
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
			fmt.Printf(string(userContent))
			if userContent == "welcome" {
				stopWrite = true
			}
			if !stopWrite {
				continue
			}
			time.Sleep(time.Second * 1)
			var bn string = "\b"
			for _, _ = range userContent {
				bn += "\b"
			}
			fmt.Printf(bn)
			for _, v := range string(userContent) {
				if v == 0 {
					break
				}
				fmt.Printf("..")
			}
			fmt.Println()
		}
	}()
	for {
		var content string
		fmt.Scanln(&content)
		if stopWrite {
			break
		}
		conn.Write([]byte(content))
	}
	var w chan int
	<-w
}
