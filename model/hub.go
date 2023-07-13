package model

type Hub struct {
	Broadcast  chan string          //broadcast管道里有数据时把它写入每一个Client的send管道中
	Clients    map[*Client]struct{} //Hub持有每个client的指针
	Register   chan *Client         //注册管道
	Unregister chan *Client         //注销管道
}
