package main

import (
	"fmt"
	"log"
	"net/http"

	r "gopkg.in/gorethink/gorethink.v3"
)

func main() {
	conn, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "symposium",
	})
	if err != nil {
		log.Panic(err.Error())
	}
	defer conn.Close()

	router := NewRouter(conn)

	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subChannel)
	router.Handle("channel unsubscribe", unsubChannel)

	router.Handle("user edit", editUser)
	router.Handle("user subscribe", subUser)
	router.Handle("user unsubscribe", unsubUser)

	router.Handle("message add", addChannelMessage)
	router.Handle("message subscribe", subChannelMessage)
	router.Handle("message unsubscribe", unsubChannelMessage)

	http.Handle("/", router)

	fmt.Println("Server listening at 8888")
	http.ListenAndServe(":8888", nil)
}
