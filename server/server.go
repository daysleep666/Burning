package main

import (
	"fmt"
	"net"
	"os"

	"github.com/daysleep666/Burning/tool"
)

//------------------------------------------------------------------------------------------------------------------------------

type connection struct {
	nickName string
	ip       string
	conn     net.Conn
}

func (c *connection) send(_content string) error {
	content := c.nickName + ":" + _content
	_, err := c.conn.Write([]byte(content))
	return err
}

func (c *connection) reading() (string, error) {
	userContent := make([]byte, 10000)
	n, err := c.conn.Read(userContent)
	return string(userContent[:n]), err
}

func (c *connections) sendAll(_contentChan chan string) {
	content := <-_contentChan
	for _, v := range c.conns {
		v.send(content)
	}
}

//------------------------------------------------------------------------------------------------------------------------------

type connections struct {
	conns []*connection
}

func (c *connections) add(_conn *connection) error {
	for _, v := range c.conns {
		if v.ip == _conn.ip {
			return fmt.Errorf("%v has joined", v.ip)
		}
	}
	c.conns = append(c.conns, _conn)
	return nil
}

func (c *connections) delete(_ip string) (*connection, error) {
	for i, v := range c.conns {
		if v.ip == _ip {
			if len(c.conns) == 1 {
				c.conns = []*connection{}
			} else {
				c.conns = append(c.conns[0:i], c.conns[i+1:]...)
			}
			return v, nil
		}
	}
	return nil, fmt.Errorf("not exist %v", _ip)
}

func (c *connections) init(_conn *connection) error {
	var (
		err      error = fmt.Errorf("nickname")
		nickname string
	)

	// nickname
	for err != nil {
		// type： writer or reader
		_conn.send("please input your nickname\n")
		nickname, err = _conn.reading()
		if err != nil {
			c.delete(_conn.ip)
			return err
		}
		_conn.nickName = nickname
	}
	err = c.add(_conn)
	if err != nil {
		_conn.send("failed")
		return err
	}
	_conn.send("welcome")
	return nil
}

//------------------------------------------------------------------------------------------------------------------------------

func newConnection(_conn net.Conn) *connection {
	return &connection{ip: _conn.LocalAddr().String(), conn: _conn}
}

func main() {
	port := os.Args[1]
	fmt.Println("Burning after reading v1.0")
	fmt.Println("		bind in ", port)
	ln, err := net.Listen("tcp", "0.0.0.0"+":"+port)
	tool.CheckErr(err)
	contentChan := make(chan string)
	var conns connections

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			connection := newConnection(conn)
			err = conns.init(connection)
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("错误:%v", err)))
				continue
			}
			fmt.Printf("%v joined\n", connection.nickName)

			go func() {
				for {
					content, err := connection.reading()
					if err != nil {
						conn.Write([]byte(fmt.Sprintf("错误:%v", err)))
						break
					}
					contentChan <- content
				}
			}()
		}
	}()

	for {
		conns.sendAll(contentChan)
	}

	var w chan int
	<-w
}