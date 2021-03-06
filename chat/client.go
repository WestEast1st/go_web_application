package main

/*
clientのモデル化
*/
import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//client is one chat user
type client struct {
	//client is websocket
	socket *websocket.Conn
	//send is message channel
	send chan *message
	//chan is client chat room
	room *room
	// userDataはユーザーに関する情報を保持します
	userData map[string]interface{}
}

//method read
func (c *client) read() {
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		fmt.Println(err)
		if err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
			log.Println("受信しました")
		} else {
			break
		}
	}
	c.socket.Close()
}

//method write
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
		log.Println("送信しました")
	}
	c.socket.Close()
}
