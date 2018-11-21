package main

import (
	"fmt"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var m sync.RWMutex

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		clearTerminal()
		termbox.Close()
	}()
	clearTerminal()
	go func() {
		for {
			shua(1, 1)
		}
	}()

	go func() {
		for {
			shua(2, 1)
		}
	}()

	go func() {
		for {
			shua(3, 1)
		}
	}()

	var content string
	var index int
	go func() {
		for {
			m.RLock()
			fmt.Printf("\033[35m\033[5;%vH%v", index, content)
			m.RUnlock()
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEnter:
				m.Lock()
				index = 0
				content = ""
				fmt.Printf("\033[33m\033[5;1H\033[K")
				m.Unlock()

			case termbox.KeyBackspace2:
				m.Lock()
				if index != 0 {
					index--
				}
				content = ""
				fmt.Printf("\033[34m\033[5;%vH\033[K", index+1)
				m.Unlock()

			case termbox.KeyCtrlC:
				fmt.Printf("You press ctrl c")
				return

			default:
				index++
				content = string(ev.Ch)

				// break
			}
		}
	}

	var w chan int
	<-w
}

func shua(_row, _column int) {
	fmt.Printf("\033[36m\033[%v;%vHi'm text", _row, _column)
	time.Sleep(time.Second)
	fmt.Printf("\033[%v;%vH\033[999D\033[K", _row, _column)
}

func clearTerminal() {
	fmt.Printf("\033[2J")
}
