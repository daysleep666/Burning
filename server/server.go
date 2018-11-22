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
	_, err := c.conn.Write([]byte(_content))
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
	c.conns = append(c.conns, _conn)
	return nil
}

func (c *connections) delete(_conn *connection) (*connection, error) {
	for i, v := range c.conns {
		if v.ip == _conn.ip && v.nickName == _conn.nickName {
			if len(c.conns) == 1 {
				c.conns = []*connection{}
			} else {
				c.conns = append(c.conns[0:i], c.conns[i+1:]...)
			}
			v.conn.Close()
			return v, nil
		}
	}
	return nil, fmt.Errorf("not exist")
}

func (c *connections) init(_conn *connection) error {
	var (
		err      error = fmt.Errorf("nickname")
		nickname string
	)

	// nickname
	for err != nil {
		_conn.send("please input your nickname")
		nickname, err = _conn.reading()
		if err != nil {
			c.delete(_conn)
			return err
		}
		_conn.nickName = nickname
	}
	err = c.add(_conn)
	if err != nil {
		_conn.send("failed")
		return err
	}
	_conn.send("welcome. HelpCommand:[\\help]")
	return nil
}

//------------------------------------------------------------------------------------------------------------------------------

func newConnection(_conn net.Conn) *connection {
	return &connection{ip: _conn.LocalAddr().String(), conn: _conn}
}

func IsQuery(_content string) (string, bool) {
	switch _content {
	case "\\help":
		return fmt.Sprintf("查询总人数：【\\tp】"), true

	case "\\tp":
		return fmt.Sprintf("TotalPopulation:%v", totalCount), true

	default:
		return "", false
	}
}

var totalCount int64

func main() {
	port := os.Args[1]
	fmt.Println("Burning after reading v0.2    bind in ", port)
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
			totalCount++
			go func() {
				for {
					content, err := connection.reading()
					if err != nil {
						connection.send(fmt.Sprintf("错误:%v", err))
						conns.delete(connection)
						totalCount--
						fmt.Printf("%v left\n", connection.nickName)
						break
					}

					if str, isQuery := IsQuery(content); isQuery {
						connection.send(str)
						continue
					}

					content = connection.nickName + ": " + content
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
