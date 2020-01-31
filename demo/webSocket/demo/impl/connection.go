package impl

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	mutex     sync.Mutex
	isClosed  bool
}

func NewConnection(webWsConn *websocket.Conn) *Connection {
	conn := &Connection{
		wsConn:    webWsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	go conn.writeLoop()

	go conn.readLoop()

	return conn
}

// API
func (conn *Connection) ReadMessage() (data []byte, err error) {
	fmt.Println("ReadMessage start")
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	fmt.Println("ReadMessage end")
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	fmt.Println("WriteMessage start")
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	fmt.Println("WriteMessage end")
	return
}

func (conn *Connection) Close() {

	//线程安全
	conn.wsConn.Close()
	//保证这里只执行一次
	conn.mutex.Lock()
	fmt.Println("调用close",conn.isClosed)
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()

}

//内部实现
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
		m int
	)
	for {
		if m, data, err = conn.wsConn.ReadMessage(); err != nil {
			fmt.Println("readLoop",m, data, err )
			goto ERR
		}
		fmt.Println("readLoop",m, data, err )
		select {
		//满时会阻塞
		case conn.inChan <- data:
			fmt.Println("读消息")
			//closeChan关闭时
		case a:=<-conn.closeChan:
			fmt.Println("读消息关闭",a)
			goto ERR
		}
	}
ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
		case a:=<-conn.closeChan:
			fmt.Println("写消息关闭",a)
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}

	}
ERR:
	conn.Close()
}
