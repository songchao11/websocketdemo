package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  64 * 1024,
	WriteBufferSize: 64 * 1024,
	//跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RegisterData struct {
	ConnKey string `json:"connKey"`
}

type WebSocketMessage struct {
	Msg     string  `json:"msg"`
}

//连接map
var connMap = sync.Map{}

//创建hatvideo websocket连接
func WsHandler(c *gin.Context) {
	param := c.Query("param")

	//解析参数获取websocket注册的connKey
	registerData := &RegisterData{}
	err := json.Unmarshal([]byte(param), &registerData)
	if err != nil {
		log.Errorf("websocket连接请求参数解析错误:%s", err)
		return
	}

	//初始化连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("初始化websocket失败:{}", err)
		return
	}
	log.Infof("初始化websocket成功,connKey:%s", registerData.ConnKey)

	// 业务处理协程
	done := make(chan int)
	go func() {
		defer conn.Close()
		defer close(done)

		//将websocket连接放入ConnMap
		connMap.Store(registerData.ConnKey, conn)

		readMessage(conn, registerData.ConnKey)
	}()

	// ping协程
	go func() {
		tick := time.NewTicker(5 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
			case <-done:
				return
			}
		}
	}()
}

//读websocket message
func readMessage(conn *websocket.Conn, connKey string) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Errorf("读取message失败:%s, connKey:%s", err, connKey)
			}
			//删除connMap存放的websocket连接
			connMap.Delete(connKey)
			break
		}

		wsm := &WebSocketMessage{}
		err = json.Unmarshal(message, &wsm)
		if err != nil {
			log.Errorf("json Unmarshal failed:%s ", err)
			continue
		}
		//向每个连接发送消息
		connMap.Range(func(k, v interface{}) bool {
			msg := connKey + ":" + wsm.Msg
			WriteBytesMsg(v.(*websocket.Conn), []byte(msg))
			return true
		})
	}
}

//发送消息
func WriteBytesMsg(conn *websocket.Conn, message []byte) {
	log.Infof("con :%v,content:%s", conn.RemoteAddr(),string(message))
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Errorf("write message errs:%s", err)
	}
}