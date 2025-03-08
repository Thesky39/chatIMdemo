package service

import (
	"chat/cache"
	"chat/conf"
	"chat/pkg/e"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"strconv"
	"time"

	//"golang.org/x/crypto/openpgp/packet"
	"net/http"
)

const month = 60 * 60 * 24 * 30 //一个月

type SenMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}

func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil) //升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建一个用户实例
	client := &Client{
		ID:     CreateID(uid, toUid), // 1->2
		SendID: CreateID(toUid, uid), // 2->1
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户组测到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SenMsg)
		//c.Socket.ReadJSON(sendMsg)
		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			fmt.Println("数据格式不对", err)
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		if sendMsg.Type == 1 {
			r1, _ := cache.RedisClient.Get(c.ID).Result()
			r2, _ := cache.RedisClient.Get(c.SendID).Result()
			if r1 > "3" && r2 == "" { // 1给2发消息，发三条，若2没回或未读，
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "达到极限",
				}
				msg, _ := json.Marshal(replyMsg) //序列化
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(c.ID)
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30*3).Result()
				//防止过快“分手”，建立连接三个月

			}
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 {
			//获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content)
			if err != nil {
				timeT = 999999
			}
			results, _ := FindMany(conf.MongoDBName, c.SendID, c.ID, int64(timeT), 10) // 获取10条历史消息
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "达到限制",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case msssage, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: fmt.Sprintf("%s", string(msssage)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
