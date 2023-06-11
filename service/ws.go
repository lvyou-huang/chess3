package service

import (
	"bytes"
	"chess/model"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Hub struct {
	Broadcast  chan string          //broadcast管道里有数据时把它写入每一个Client的send管道中
	Clients    map[*Client]struct{} //Hub持有每个client的指针
	Register   chan *Client         //注册管道
	Unregister chan *Client         //注销管道
}

type Client struct {
	Hub   *Hub
	Conn  *websocket.Conn
	Send  chan []byte
	Name  []byte
	Chess *model.Chess_Table
}

const (
	writeWait  = 10 * time.Second  //
	pongWait   = 60 * time.Second  // 每60秒向websocket发送一次pong
	pingPeriod = 9 * pongWait / 10 //连接不断时每隔54秒向client发送一次ping
	maxMsgSize = 512               //消息的长度不能超过512
)

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, chess *model.Chess_Table) {
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

	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256), Chess: chess}
	hub.Register <- client
	go client.read()
	go client.write()
}
func (client *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop() //ticker不用就stop，防止协程泄漏
		fmt.Printf("close connection to %s\n", client.Conn.RemoteAddr().String())
		client.Conn.Close() //给前端写数据失败，就可以关系连接了
	}()
	for {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				fmt.Println("管道已经关闭")
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait)) //10秒内必须把信息写给前端（写到websocket连接里去），否则就关闭连接
			if writer, err := client.Conn.NextWriter(websocket.TextMessage); err != nil {
				log.Println(err)
				return
			} else {
				writer.Write(msg)
				writer.Write([]byte{'\n'})
				// 有消息一次全写出去
				n := len(client.Send)
				for i := 0; i < n; i++ {
					writer.Write(<-client.Send)
					writer.Write([]byte{'\n'})
				}
				if err := writer.Close(); err != nil { //必须调close，否则下次调用client.conn.NextWriter时本条消息才会发送给浏览器
					log.Println(err)
					return //结束一切
				}
			}
			//超时
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// 从websocket读取数据
func (client *Client) read() {
	defer func() {
		client.Hub.Unregister <- client //向hub发送注销
		fmt.Printf("%s offline\n", client.Name)
		fmt.Printf("close connection to %s\n", client.Conn.RemoteAddr().String())
		client.Conn.Close() //关闭ws连接
	}()

	// conn细节设置
	client.Conn.SetReadLimit(maxMsgSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait)) //设置最长可读时间
	client.Conn.SetPongHandler(func(appData string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait)) //每次接收到ping后都将最长可读时间延后60秒
		return nil
	})

	for {
		_, p, err := client.Conn.ReadMessage() //返回消息类型，消息，error
		if err != nil {
			//如果以意料之外的关闭状态关闭，就打印日志
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				fmt.Printf("close websocket conn error: %v\n", err)
			}
			break //只要ReadMessage失败，就关闭websocket管道、注销client，退出
		} else {
			// trimspace:消去首尾空格， replace：将换行符换位空格，-1：全部转换
			message := bytes.TrimSpace(bytes.Replace(p, []byte{'\n'}, []byte{' '}, -1))
			order := string(message)

			//处理棋子的位置
			if order[0] == 'c' && order[1] == 'h' {
				a, _ := strconv.Atoi(string(order[3]))
				b, _ := strconv.Atoi(string(order[4]))
				x, _ := strconv.Atoi(string(order[6]))
				y, _ := strconv.Atoi(string(order[7]))
				if (a > 7 || a < 0) || (b > 7 || b < 0) || (x > 7 || x < 0) || (y > 7 || y < 0) {
					client.Hub.Broadcast <- "超界了"
				} else {
					before := [2]int{a, b}
					next := [2]int{x, y}
					if client.Chess.Chess[x][y]*client.Chess.Chess[a][b] > 0 {
						client.Hub.Broadcast <- "不能吃自己的子"
					} else {
						if client.Chess.Chess[before[0]][before[1]] == -3 || client.Chess.Chess[before[0]][before[1]] == 3 {
							if (math.Abs(float64(before[0]-next[0])) == 1 && math.Abs(float64(before[1]-next[1])) == 2) || (math.Abs(float64(before[1]-next[1])) == 1 && math.Abs(float64(before[0]-next[0])) == 2) {
								client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
							} else {
								client.Hub.Broadcast <- "非法移动"
							}

						}
						if client.Chess.Chess[a][b] == 6 || client.Chess.Chess[a][b] == -6 {
							if math.Abs(float64(a-x)) <= 1 && math.Abs(float64(b-y)) <= 1 {
								client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
							}
						}
						if client.Chess.Chess[a][b] == 2 || client.Chess.Chess[a][b] == -2 {
							if Blockbool(client.Chess, a, b, x, y) {
								if a-x == 0 || b-y == 0 {
									client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
								} else {
									client.Hub.Broadcast <- "非法移动"
								}
							} else {
								client.Hub.Broadcast <- "堵塞不行"
							}
						}
						if client.Chess.Chess[a][b] == 4 || client.Chess.Chess[a][b] == -4 {
							if Blockbool(client.Chess, a, b, x, y) {
								if math.Abs(float64(a-x)) == math.Abs(float64(b-y)) {
									client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
								} else {
									client.Hub.Broadcast <- "非法移动"
								}
							} else {
								client.Hub.Broadcast <- "堵塞不行"
							}
						}
						if client.Chess.Chess[a][b] == 5 || client.Chess.Chess[a][b] == -5 {
							if Blockbool(client.Chess, a, b, x, y) {
								if math.Abs(float64(a-x)) == math.Abs(float64(b-y)) || a-x == 0 || b-y == 0 {
									client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
								} else {
									client.Hub.Broadcast <- "非法移动"
								}
							} else {
								client.Hub.Broadcast <- "堵塞不行"
							}
						}
						/*if (client.Chess.Chess[a][b] == 6 || client.Chess.Chess[a][b] == -6) && (client.Chess.Chess[x][y] != 6 || client.Chess.Chess[x][y] != -6) {
							if (a-x == 0 && b-y == 1) || (a-x == 1 && b-y == 0) {
								client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
							} else {
								client.Hub.Broadcast <- "非法移动"
							}
						}*/
						if client.Chess.Chess[a][b] == 1 || client.Chess.Chess[a][b] == -1 {
							if client.Chess.Chess[a][b] > 0 {
								if client.Chess.Chess[x][y] == 0 && y-b == 0 && x-a == -1 {
									client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
									if x == 0 || x == 7 {
										client.Hub.Broadcast <- "请选择升变的棋子序号"
									}
								} else if client.Chess.Chess[x][y] != 0 && math.Abs(float64(y-b)) == 1 && x-a == -1 {
									client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
								} else {
									client.Hub.Broadcast <- "非法移动"
								}
							} else if client.Chess.Chess[a][b] < 0 {
								if client.Chess.Chess[x][y] == 0 && y-b == 0 && x-a == 1 {
									client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
									if x == 0 || x == 7 {
										client.Hub.Broadcast <- "请选择升变的棋子序号"

									}
								} else if client.Chess.Chess[x][y] != 0 && math.Abs(float64(y-b)) == 1 && x-a == 1 {
									client.Chess.Chess[a][b], client.Chess.Chess[x][y] = 0, client.Chess.Chess[a][b]
								} else {
									client.Hub.Broadcast <- "非法移动"
								}
							}
						}
						white := 0
						black := 0
						for i := 0; i < 8; i++ {
							for j := 0; j < 8; j++ {
								if client.Chess.Chess[i][j] == 6 {
									white = 1
								} else if client.Chess.Chess[i][j] == -6 {
									black = 1
								}
							}
						}
						if black == 0 {
							client.Hub.Broadcast <- "白方胜利"
						} else if white == 0 {
							client.Hub.Broadcast <- "黑方胜利"
						}

					}
				}
			} else if order[0] == 's' && order[1] == 'b' {
				a, _ := strconv.Atoi(string(order[3]))
				b, _ := strconv.Atoi(string(order[4]))
				c, _ := strconv.Atoi(string(order[6]))
				if c == 2 || c == 3 || c == 4 || c == 5 {
					client.Chess.Chess[a][b] = c
				} else {
					client.Hub.Broadcast <- "序列号错误"
				}
			} else if order[0] == 'w' && order[1] == 'c' {
				a, _ := strconv.Atoi(string(order[3]))
				b, _ := strconv.Atoi(string(order[4]))
				x, _ := strconv.Atoi(string(order[6]))
				y, _ := strconv.Atoi(string(order[7]))
				if Blockbool(client.Chess, a, b, x, y) == true {
					if client.Chess.Chess[a][b] == 6 && client.Chess.Chess[x][y] == 2 {
						if x > a {
							client.Chess.Chess[a][b], client.Chess.Chess[a+2][b] = 0, 6
							client.Chess.Chess[x][y], client.Chess.Chess[a+1][b] = 0, 2
						}
					}
				} else {
					client.Hub.Broadcast <- "堵塞不行"
				}

			}
			if len(client.Name) == 0 {
				client.Name = message
			} else {
				client.Hub.Broadcast <- string(bytes.Join([][]byte{client.Name, message}, []byte(": ")))
				mes := ""
				for i := 0; i < 8; i++ {
					for j := 0; j < 8; j++ {
						mes += strconv.Itoa(client.Chess.Chess[i][j]) + "\t"
					}
					mes += "\n"
				}
				client.Hub.Broadcast <- mes
			}

		}
	}
}

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
