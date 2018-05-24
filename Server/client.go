package main

import (
	"log"

	"github.com/gorilla/websocket"
	r "gopkg.in/gorethink/gorethink.v3"
)

type FindTheHandle func(string) (Handler, bool)

type Client struct {
	send         chan Message
	socket       *websocket.Conn
	findHandler  FindTheHandle
	session      *r.Session
	stopChannels map[int]chan bool
	userName     string
	id           string
}

func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	close(c.send)
	r.Table("user").Get(c.id).Delete().Exec(c.session)
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}

func (c *Client) NewStopChannel(stopKey int) chan bool {
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			log.Println(err)
			break
		}
	}

	client.socket.Close()
}

func (client *Client) Read() {
	var msg Message
	for {
		if err := client.socket.ReadJSON(&msg); err != nil {
			log.Println(msg)
			break
		}

		if handler, found := client.findHandler(msg.Name); found {
			handler(client, msg.Data)
		}
	}

	client.socket.Close()
}

func NewClient(sock *websocket.Conn, find FindTheHandle, sess *r.Session) *Client {
	var user User
	user.Name = "Incognito"
	res, err := r.Table(userTable).Insert(user).RunWrite(sess)
	if err != nil {
		log.Println(err.Error())
	}

	var uid string
	if len(res.GeneratedKeys) > 0 {
		uid = res.GeneratedKeys[0]
	}

	return &Client{
		send:         make(chan Message),
		socket:       sock,
		findHandler:  find,
		session:      sess,
		stopChannels: make(map[int]chan bool),
		userName:     user.Name,
		id:           uid,
	}
}
