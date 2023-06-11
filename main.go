package main

import (
	"chess/api"
	"chess/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	hub := service.NewHub()
	chess := service.NewChess()
	go hub.Run()
	Router := gin.Default()
	Router.POST("/login", api.Login)
	Router.POST("/register", api.Register)
	Router.Run(":8686")
	http.HandleFunc("/", service.ServeHome)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		service.ServeWS(hub, w, r, chess)
	})

	if err := http.ListenAndServe(":5656", nil); err != nil {
		fmt.Printf("start http service error: %s\n", err)
	}

}
