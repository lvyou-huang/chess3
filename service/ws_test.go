package service

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	go func() {
		hub := NewHub()
		chess := NewChess()
		go Run(*hub)
		http.HandleFunc("/", ServeHome)
		http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
			ServeWS(hub, w, r, chess)
		})

		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("start http service error: %s\n", err)
		}
	}()
	time.Sleep(600 * time.Second)
}
func TestRead(t *testing.T) {
	go func() {
		hub := NewHub()
		chess := NewChess()
		go Run(*hub)
		http.HandleFunc("/", ServeHome)
		http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
			ServeWS(hub, w, r, chess)
		})

		if err := http.ListenAndServe(":5656", nil); err != nil {
			fmt.Printf("start http service error: %s\n", err)
		}
	}()
	time.Sleep(600 * time.Second)
}
