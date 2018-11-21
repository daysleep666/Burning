package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/daysleep666/Burning/tool"
	termbox "github.com/nsf/termbox-go"
)

var m sync.RWMutex

func clearTerminal() {
	fmt.Printf("\033[2J")
}

func main() {
	var (
		myMsg    string
		index    int
		contents chan string = make(chan string, 1000)
	)
	// init
	ip := os.Args[1]
	conn, err := net.Dial("tcp", ip)
	tool.CheckErr(err)
	err = termbox.Init()
	tool.CheckErr(err)
	clearTerminal()
	defer func() {
		clearTerminal()
		termbox.Close()
		conn.Close()
	}()

	// read
	go func() {
		for {
			bs := make([]byte, 10000)
			_, err := conn.Read(bs)
			if err != nil {
				return
			}
			contents <- string(bs)
		}
	}()

	go func() {
		for {
			content := <-contents
			fmt.Printf("\033[36m\033[1;1H%v", content)
			time.Sleep(time.Second)
			fmt.Printf("\033[1;1H\033[999D\033[K")
		}
	}()

	// write
	go func() {
		for {
			m.RLock()
			fmt.Printf("\033[35m\033[5;1H%v", myMsg)
			m.RUnlock()
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				m.Lock()
				conn.Write([]byte(myMsg))
				index = 0
				myMsg = ""
				fmt.Printf("\033[33m\033[5;1H\033[K")
				m.Unlock()

			case termbox.KeyBackspace2:
				m.Lock()
				if index != 0 {
					index--
					myMsg = myMsg[:len(myMsg)-1]
				}
				fmt.Printf("\033[34m\033[5;%vH\033[K", index+1)
				m.Unlock()

			case termbox.KeySpace:
				index++
				myMsg += string(ev.Ch)

			case termbox.KeyCtrlC:
				fmt.Printf("You press ctrl c")
				return

			default:
				index++
				myMsg += string(ev.Ch)
				// break
			}
		}
	}

	var w chan int
	<-w
}

func display(_row, _column int, _content string) {
	fmt.Printf("\033[36m\033[%v;%vH%v", _row, _column, _content)
	time.Sleep(time.Second * 3)
	fmt.Printf("\033[1;1H\033[999D\033[K")
}
