package controllers

import (
	//"ad/models"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

/*


调用demo
	beego.Router("/", &controllers.MyWebSocketController{}, "get:ChatRoom")
	beego.Router("/ws", &controllers.MyWebSocketController{})




//服务器接收消息
type Message struct {
	Imei string `json:"imei"` //哪台机器
	Act  string `json:"act"`  //什么指令
	Key  string `json:"key"`  //指令键名
	Val  string `json:"val"`  //指令内容
}

//{"imei":"1112312312", "act":"setnick", "key":"1112312312", "val":"helloKitty"}               //设置昵称
//{"imei":"1112312312", "act":"getalipayuserid", "key":"1112312312", "val":"1"}                //获取ID<-
//{"imei":"1112312312", "act":"getalipayuserid", "key":"208888888", "val":"华硕品质，坚若磐石"}  //获取ID-> 服务器返回的
//{"imei":"1112312312", "act":"setalipayuser", "key":"20888888888", "val":"201900000099"}      //设置ID
//{"imei":"1112312312", "act":"setusercomefrom", "key":"URL", "val":"ip"}                      //设置用户反馈 KEY 是用户访问的URL, VAL:是用户返回的ip
//用户反馈功能还没有想好，因为不知道用户是这么来到我的网页的


*/

type MyWebSocketController struct {
	beego.Controller
}

//var upgrader = websocket.Upgrader{}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	var ret *Hub
	ret = &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go ret.run()
	return ret
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

/*
func (c *MyWebSocketController) Get() {
	imei := c.GetString("imei")
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	clients[ws] = imei

	//不断的广播发送到页面上
	//msg := models.Message{Message: "这是向页面发送的数据 " + time.Now().Format("2006-01-02 15:04:05")}
	//broadcast <- msg
	//for {
	//目前存在问题 定时效果不好 需要在业务代码替换时改为beego toolbox中的定时器
	//time.Sleep(time.Second * 3)
	//msg := models.Message{Message: "这是向页面发送的数据 " + time.Now().Format("2006-01-02 15:04:05")}
	//broadcast <- msg
	//}
}*/
var hub *Hub = newHub()

func (c *MyWebSocketController) Get() {
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	c.Ctx.WriteString("")
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			ret := string(message)
			fmt.Println(ret) //string(message),
			//message = []byte("服务器说:好的,我收到了")
			if 1 == 1 {
				for client := range h.clients {
					select {
					case client.send <- []byte(ret):
						//client.send <- []byte("<div align=right>end</div>\r\n")
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//message = []byte("helloWorld")
		c.hub.broadcast <- message
	}
}
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			//w.Write([]byte("message"))
			w.Write([]byte(Query(message)))
			//fmt.Println(message)
			//在这里调用
			//w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *MyWebSocketController) ChatRoom() {
	c.Ctx.WriteString(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<title>WebSocket测试</title>
		<script type="text/javascript">
		window.onload = function () {
		    var conn;
		    var msg = document.getElementById("msg");
		    var log = document.getElementById("log");

		    function appendLog(item) {
		        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
		        log.appendChild(item);
		        if (doScroll) {
		            log.scrollTop = log.scrollHeight - log.clientHeight;
		        }
		    }

		    document.getElementById("form").onsubmit = function () {
		        if (!conn) {
		            return false;
		        }
		        if (!msg.value) {
		            return false;
		        }
		        conn.send(msg.value);
		        msg.value = "";
		        return false;
		    };

		    if (window["WebSocket"]) {
		        conn = new WebSocket("ws://" + document.location.host + "/ws");
		        conn.onclose = function (evt) {
		            var item = document.createElement("div");
		            item.innerHTML = "<b>Connection closed.</b>";
		            appendLog(item);
		        };
		        conn.onmessage = function (evt) {
		            var messages = evt.data.split('\n');
		            for (var i = 0; i < messages.length; i++) {
		                var item = document.createElement("div");
		                item.innerText = messages[i];
		                appendLog(item);
		            }
		        };
		    } else {
		        var item = document.createElement("div");
		        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
		        appendLog(item);
		    }
		};
		</script>
		<style type="text/css">
		html {
		    overflow: hidden;
		}

		body {
		    overflow: hidden;
		    padding: 0;
		    margin: 0;
		    width: 100%;
		    height: 100%;
		    background: gray;
		}

		#log {
		    background: white;
		    margin: 0;
		    padding: 0.5em 0.5em 0.5em 0.5em;
		    position: absolute;
		    top: 0.5em;
		    left: 0.5em;
		    right: 0.5em;
		    bottom: 3em;
		    overflow: auto;
		}

		#form {
		    padding: 0 0.5em 0 0.5em;
		    margin: 0;
		    position: absolute;
		    bottom: 1em;
		    left: 0px;
		    width: 100%;
		    overflow: hidden;
		}

		</style>
		</head>
		<body>
		<div id="log"></div>
		<form id="form">
		    <input type="submit" value="Send" />
		    <input type="text" id="msg" size="64"/>
		</form>
		</body>
		</html>`)
}
