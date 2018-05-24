package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	r "gopkg.in/gorethink/gorethink.v3"
)

type Router struct {
	rules   map[string]Handler
	session *r.Session
}
type Handler func(*Client, interface{})

func testOrigin(r *http.Request) bool { return true }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     testOrigin,
}

func (r *Router) Handle(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	client := NewClient(socket, e.FindHandler, e.session)
	defer client.Close()
	go client.Write()
	client.Read()
}

func NewRouter(sess *r.Session) *Router {
	return &Router{
		rules:   make(map[string]Handler),
		session: sess,
	}
}

func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}
