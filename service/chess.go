package service

import "chess/model"

// 兵=1 车=2 马=3 象=4 后=5 王=6
func NewChess() *model.Chess_Table {
	return &model.Chess_Table{
		Chess: [8][8]int{
			{-2, -3, -4, -6, -5, -4, -3, -2},
			{-1, -1, -1, -1, -1, -1, -1, -1},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{1, 1, 1, 1, 1, 1, 1, 1},
			{2, 3, 4, 5, 6, 4, 3, 2},
		},
		Player: 1,
	}
}

func Blockbool(chess *model.Chess_Table, a int, b int, x int, y int) bool {
	if a != x && b != y {
		if x > a && y > b {
			for i := 1; a+i < x; i++ {
				if chess.Chess[a+i][b+i] != 0 {
					return false
				}
			}
		} else if x > a && y < b {
			for i := 1; a+i < x; i++ {
				if chess.Chess[a+i][b-i] != 0 {
					return false
				}
			}
		} else if x < a && y > b {
			for i := 1; x+i < a; i++ {
				if chess.Chess[x+i][b+i] != 0 {
					return false
				}
			}
		} else if x < a && y < b {
			for i := 1; x+i < a; i++ {
				if chess.Chess[x+i][y+i] != 0 {
					return false
				}
			}
		}
	} else if a == x {
		if y > b {
			for i := 1; b+i < y; i++ {
				if chess.Chess[a][b+i] != 0 {
					return false
				}
			}
		} else if y < b {
			for i := 1; y+i < b; i++ {
				if chess.Chess[a][y+i] != 0 {
					return false
				}
			}
		}
	} else if b == y {
		if x > a {
			for i := 1; a+i < x; i++ {
				if chess.Chess[a+i][b] != 0 {
					return false
				}
			}
		} else if x < a {
			for i := 1; x+i < a; i++ {
				if chess.Chess[x+i][y] != 0 {
					return false
				}
			}
		}
	}
	return true

}
