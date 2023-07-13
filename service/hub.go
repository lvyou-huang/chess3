package service

import "chess/model"

func NewHub() *model.Hub {
	return &model.Hub{Broadcast: make(chan string),
		Clients:    make(map[*model.Client]struct{}),
		Register:   make(chan *model.Client),
		Unregister: make(chan *model.Client)}
}

func Run(hub model.Hub) {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = struct{}{}
		case client := <-hub.Unregister:
			delete(hub.Clients, client)
			close(client.Send)
		case msg := <-hub.Broadcast:
			for client := range hub.Clients {
				select {
				case client.Send <- []byte(msg): //如果管道不能立即写入数据，就认为该client出故障了
				default:
					close(client.Send)
					delete(hub.Clients, client)
				}
			}
		}
	}
}
