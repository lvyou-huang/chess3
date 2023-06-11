package service

func NewHub() *Hub {
	return &Hub{Broadcast: make(chan string),
		Clients:    make(map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client)}
}
func (hub Hub) Run() {
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
