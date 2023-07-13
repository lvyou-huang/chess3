package main

import (
	"chess/api"
	"chess/service"
	"fmt"
	"net/http"
)

func main() {
	hub := service.NewHub()
	chess := service.NewChess()
	go service.Run(*hub)
	//Router := gin.Default()
	//Router.POST("/login", api.Login)
	//Router.POST("/register", api.Register)
	//Router.Run(":8686")
	http.HandleFunc("/", api.ServeHome)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		api.ServeWS(hub, w, r, chess)
	})

	if err := http.ListenAndServe(":5656", nil); err != nil {
		fmt.Printf("start http service error: %s\n", err)
	}
}
