package wsConn

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type wsMessage struct {
	messageType int
	data        []byte
}
type WsConnection struct {
	wsSocket *websocket.Conn
	inChan   chan *wsMessage
	outChan  chan *wsMessage

	//关闭通道相关
	mutex     sync.Mutex
	closeChan chan byte
	isClosed  bool
}

func NewWsConn(w http.ResponseWriter, r *http.Request) (*WsConnection, error) {
	wsUp := websocket.Upgrader{
		HandshakeTimeout: time.Second * 5,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsSocket, err := wsUp.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	//初始化
	wsConn := &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	return wsConn, nil
}

func (s *WsConnection) Init() {
	// 心跳检测
	//go wsConn.ping()
	// 读
	go s.readLoop()
	// 写
	go s.writeLoop()
}

// 阻塞读取客户端发来的消息
// 如果没有消息发送过来就一直阻塞
func (s *WsConnection) readLoop() {
	for {
		msgType, data, err := s.wsSocket.ReadMessage()
		fmt.Println(msgType, string(data))
		if err != nil {
			s.Close()
			return
		}
		select {
		case s.inChan <- &wsMessage{messageType: msgType, data: data}:
		case <-s.closeChan:
			return
		}
	}

}

//从出通道中取信息发给客户端
func (s *WsConnection) writeLoop() {
	for {
		select {
		case msg := <-s.outChan:
			// 发送给客户端
			if err := s.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				s.Close()
				return
			}
		case <-s.closeChan:
			return
		}
	}
}

// 心跳检测
func (s *WsConnection) ping() {
	for {
		time.Sleep(time.Second * 2)
		if err := s.WriteMsg(websocket.TextMessage, []byte("ping")); err != nil {
			s.Close()
		}
	}
}
func (s *WsConnection) ApiWriteMsg(messageType int, data []byte) error {
	msg := wsMessage{
		messageType: messageType,
		data:        data,
	}
	fmt.Println(msg)
	s.outChan <- &msg
	return nil
}

// 往通道中写信息
func (s *WsConnection) WriteMsg(messageType int, data []byte) error {
	select {
	case s.outChan <- &wsMessage{messageType: messageType, data: data}:
		return nil
	case <-s.closeChan:
		return errors.New("websocket closed")
	}
}

// 从通道中读取信息
func (s *WsConnection) ReadMsg() (*wsMessage, error) {
	select {
	case msg := <-s.inChan:
		//写入输出通道
		fmt.Println("写信息", msg)
		err := s.WriteMsg(msg.messageType, msg.data)
		fmt.Println(err)
		return msg, nil
	case <-s.closeChan:
		return nil, errors.New("websocket closed")
	}
}

// 关闭socket连接，关闭通道
func (s *WsConnection) Close() {
	s.wsSocket.Close()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.isClosed {
		close(s.closeChan)
		s.isClosed = true
	}

}
