package model

// 兵=1 车=2 马=3 象=4 后=5 王=6
// 0空地
type Chess_Table struct {
	Chess  [8][8]int //棋盘 表示符号黑白子 白是+ 黑是负
	Player int       //+为白，-为黑
}
