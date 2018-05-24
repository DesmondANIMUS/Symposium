package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	r "gopkg.in/gorethink/gorethink.v3"
)

type Channel struct {
	UID  string `json:"uid" gorethink:"uid,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

type User struct {
	UID  string `json:"uid" gorethink:"uid,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

type ChannelMessage struct {
	UID       string `json:"uid" gorethink:"uid,omitempty"`
	ChannelID string `json:"channelId" gorethink:"channelId"`
	Body      string `json:"body" gorethink:"body"`
	Author    string `json:"author" gorethink:"author"`
	CreatedAt string `json:"createdAt" gorethink:"createdAt"`
}

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

const (
	channelTable = "channel"
	userTable    = "user"
	messageTable = "message"
)

const (
	ChannelStop = iota
	UserStop
	MessageStop
)

func addChannel(client *Client, data interface{}) {
	var channel Channel
	fmt.Println("inside channel ADD event handler")

	if err := mapstructure.Decode(data, &channel); err != nil {
		log.Println(err)
		client.send <- Message{"error", "500"}
	}

	go func() {
		if err := r.Table(channelTable).Insert(channel).Exec(client.session); err != nil {
			log.Println(err)
			client.send <- Message{"error", "500"}
		}
	}()
}

func subChannel(client *Client, data interface{}) {
	stop := client.NewStopChannel(ChannelStop)
	result := make(chan r.ChangeResponse)

	cursor, err := r.Table(channelTable).Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)
	if err != nil {
		log.Println(err)
		client.send <- Message{"error", "500"}
		return
	}

	go func() {
		var change r.ChangeResponse
		for cursor.Next(&change) {
			result <- change
		}
	}()

	go func() {
		for {
			select {
			case <-stop:
				cursor.Close()
				return
			case change := <-result:
				if change.NewValue != nil && change.OldValue == nil {
					client.send <- Message{"channel add", change.NewValue}
					fmt.Println(change.NewValue)
					fmt.Println("new channel added")
				}
			}
		}
	}()
}

func unsubChannel(client *Client, data interface{}) {
	client.StopForKey(ChannelStop)
}

func editUser(client *Client, data interface{}) {
	var user User
	err := mapstructure.Decode(data, &user)
	if err != nil {
		client.send <- Message{"error", err.Error()}
	}
	client.userName = user.Name
	go func() {
		_, err := r.Table(userTable).Get(client.id).Update(user).RunWrite(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
		}
	}()
}
func subUser(client *Client, data interface{}) {
	go func() {
		stop := client.NewStopChannel(UserStop)
		cursor, err := r.Table(userTable).Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
			return
		}

		changeFeedHelper(cursor, userTable, client.send, stop)
	}()
}
func unsubUser(client *Client, data interface{}) {
	client.StopForKey(UserStop)
}

func addChannelMessage(client *Client, data interface{}) {
	var channelMessage ChannelMessage
	err := mapstructure.Decode(data, &channelMessage)
	if err != nil {
		client.send <- Message{"error", err.Error()}
	}

	go func() {
		channelMessage.CreatedAt = time.Now().Format(time.RFC850)
		channelMessage.Author = client.userName
		err := r.Table(messageTable).Insert(channelMessage).Exec(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
		}
	}()
}
func subChannelMessage(client *Client, data interface{}) {
	go func() {
		eventData := data.(map[string]interface{})
		val, ok := eventData["channelId"]
		if !ok {
			return
		}
		channelID, ok := val.(string)
		if !ok {
			return
		}
		stop := client.NewStopChannel(MessageStop)
		cursor, err := r.Table(messageTable).OrderBy(r.OrderByOpts{Index: r.Desc("createdAt")}).
			Filter(r.Row.Field("channelId").Eq(channelID)).Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)

		if err != nil {
			client.send <- Message{"error", err.Error()}
			return
		}

		changeFeedHelper(cursor, messageTable, client.send, stop)
	}()
}
func unsubChannelMessage(client *Client, data interface{}) {
	client.StopForKey(MessageStop)
}

func changeFeedHelper(cursor *r.Cursor, changeEventName string, send chan<- Message, stop <-chan bool) {
	change := make(chan r.ChangeResponse)
	cursor.Listen(change)
	for {
		eventName := ""
		var data interface{}
		select {
		case <-stop:
			cursor.Close()
			return
		case val := <-change:
			if val.NewValue != nil && val.OldValue == nil {
				eventName = changeEventName + " add"
				data = val.NewValue
			} else if val.NewValue == nil && val.OldValue != nil {
				eventName = changeEventName + " remove"
				data = val.OldValue
			} else if val.NewValue != nil && val.OldValue != nil {
				eventName = changeEventName + " edit"
				data = val.NewValue
			}
			send <- Message{eventName, data}
		}
	}
}
