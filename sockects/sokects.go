package sockects

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	"net/http"
	// "strconv"
	"time"

	// "../api"
	"../common"
	"github.com/gorilla/websocket"
)

type WsMsg struct {
	MType string `json:"Type" `
}
type NoticeMsg struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Safe    int    `json:"safe"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

var logs = common.Loger
var notceContet string

//全局静态连接
var conns []*websocket.Conn
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理ws请求
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	var conn *websocket.Conn
	var err error
	logs.Info("WebSocketHandler enter")
	//Log("进入web socket"+time.Now().Format("2006-01-02 15:04:05"))
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Info("err:%v", err)
		return
	}
	conns = append(conns, conn)

	//gin通过协程调用该handler函数，一旦退出函数，ws会被主动销毁
	for {
		mt, message, err := conn.ReadMessage()
		//common.Log("接收消息：" + string(message))
		logs.Info("recv:%s", message)
		if err != nil {
			break
		}
		for _, item := range conns {

			err = item.WriteMessage(mt, message)
			if err != nil {
				continue
			}
		}

	}
}

func Send(msg string) {
	for _, item := range conns {
		err := item.WriteMessage(1, []byte(msg))
		if err != nil {
			continue
		}
	}
}
