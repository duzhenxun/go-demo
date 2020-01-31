package sendMsg

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type SendMsg struct {
	wsConnect   *websocket.Conn
	inChan      chan []byte
	outChan     chan []byte
	closeChan   chan byte
	mutex       sync.Mutex
	IsClosed    bool
}

func NewSendMsg(w http.ResponseWriter,r *http.Request) (sendMsg *SendMsg, err error) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	coon,err:=upGrader.Upgrade(w,r,nil)
	sendMsg = &SendMsg{
		wsConnect: coon,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}
	go sendMsg.readLoop()
	go sendMsg.writeLoop()
	return
}

func (s *SendMsg) readLoop() {
	for {
		msgType, data, err := s.wsConnect.ReadMessage()
		fmt.Println("=====================")
		fmt.Println(msgType, data, err)
		fmt.Println("=====================")
		if err != nil {
			fmt.Println(s.wsConnect.RemoteAddr(), "无法读取", err)
			s.Close()
		}
		fmt.Println(msgType)
		//将data数据放到通道中
		select {
		case s.inChan <- data:
			fmt.Println("data有数据 读出...")
		case <-s.closeChan:
			fmt.Println("关才通道有数据写入...")
			s.Close()
		}
	}
}

func (s *SendMsg) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-s.outChan:
			fmt.Println("data有数据写入...")
		case <-s.closeChan:
			fmt.Println("关闭通道有数据写入...")
			s.Close()
		}
		if err = s.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			s.Close()
		}
	}
}
func (s *SendMsg) Close() {
	s.wsConnect.Close()
	s.mutex.Lock()
	if !s.IsClosed {
		close(s.closeChan)
		s.IsClosed = true
	}
	s.mutex.Unlock()
}

func (s *SendMsg) Write(data []byte) error {
	select {
	case s.outChan <- data:
	case <-s.closeChan:
		return errors.New("closed")
	}
	return nil
}
func (s *SendMsg) Read() ([]byte, error) {
	var data []byte
	select {
	case data = <-s.inChan:
	case <-s.closeChan:
		return nil, errors.New("closed")
	}
	return data, nil
}
