package main

import (
	"fmt"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	fmt.Printf("\033[2J")
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
	var m sync.RWMutex
	go func() {
		for {
			m.RLock()
			fmt.Printf("\033[5;%vH%v", index, content)
			m.RUnlock()
			// var content string
			// fmt.Scanln(&content)
			// fmt.Printf("\033[33m\033[1A\033[999D\033[K")
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC:
				fmt.Printf("You press ctrl c")
				return
			case termbox.KeyEnter:
				m.Lock()
				index = 0
				content = ""
				fmt.Printf("\033[33m\033[5;1H\033[K")
				m.Unlock()

			case termbox.KeyBackspace2:
				m.Lock()
				index--
				content = ""
				fmt.Printf("\033[34m\033[5;%vH\033[K", index)
				m.Unlock()

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
	fmt.Printf("\033[%v;%vHi'm text", _row, _column)
	time.Sleep(time.Second)
	fmt.Printf("\033[%v;%vH\033[999D\033[K", _row, _column)
}
