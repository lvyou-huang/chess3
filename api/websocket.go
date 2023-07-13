package api

import (
	"chess/model"
	"chess/service"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
func ServeWS(hub *model.Hub, w http.ResponseWriter, r *http.Request, chess *model.Chess_Table) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 2 * time.Second, //握手超时时间
		ReadBufferSize:   1024,            //读缓冲大小
		WriteBufferSize:  1024,            //写缓冲大小
		CheckOrigin:      func(r *http.Request) bool { return true },
		Error:            func(w http.ResponseWriter, r *http.Request, status int, reason error) {},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("connect to client %s\n", conn.RemoteAddr().String())

	client := &model.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256), Chess: chess}
	hub.Register <- client
	go service.Read(client)
	go service.Write(client)
}
