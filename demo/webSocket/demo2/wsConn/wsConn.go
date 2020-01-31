package wsConn

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type wsMessage struct {
	messageType int
	data        []byte
}
type wsConnection struct {
	wsSocket *websocket.Conn
	inChan   chan *wsMessage
	outChan  chan *wsMessage

	//关闭通道相关
	mutex     sync.Mutex
	closeChan chan byte
	isClosed  bool
}

func NewWsConn(w http.ResponseWriter, r *http.Request) (*wsConnection, error) {
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
	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	return wsConn, nil
}

func (wsConn *wsConnection) Init(){
	// 心跳检测
	go wsConn.ping()
	// 读
	go wsConn.readLoop()
	// 写
	go wsConn.writeLoop()
}

// 阻塞读取客户端发来的消息
// 如果没有消息发送过来就一直阻塞
func (wsConn *wsConnection) readLoop() {
	for {
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			wsConn.Close()
			return
		}
		select {
		case wsConn.inChan <- &wsMessage{messageType: msgType, data: data}:
		case <-wsConn.closeChan:
			return
		}
	}

}

//从出通道中取信息发给客户端
func (wsConn *wsConnection) writeLoop() {
	for {
		select {
		case msg := <-wsConn.outChan:
			// 发送给客户端
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				wsConn.Close()
				return
			}
		case <-wsConn.closeChan:
			return
		}
	}
}

// 心跳检测
func (wsConn *wsConnection) ping() {
	for {
		time.Sleep(time.Second * 2)
		if err := wsConn.WriteMsg(websocket.TextMessage, []byte("ping")); err != nil {
			wsConn.Close()
		}
	}
}

// 往通道中写信息
func (wsConn *wsConnection) WriteMsg(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType: messageType, data: data}:
		return nil
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
}

// 从通道中读取信息
func (wsConn *wsConnection) ReadMsg() (*wsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
		return nil, errors.New("websocket closed")
	}
}

// 关闭socket连接，关闭通道
func (wsConn *wsConnection) Close() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		close(wsConn.closeChan)
		wsConn.isClosed = true
	}

}
