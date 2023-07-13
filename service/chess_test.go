package service

import (
	"chess/model"
	"testing"
)

func TestBlockbool(t *testing.T) {
	type args struct {
		chess *model.Chess_Table
		a     int
		b     int
		x     int
		y     int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				chess: NewChess(),
				a:     0,
				b:     0,
				x:     2,
				y:     2,
			},
			want: false,
		},
		{
			name: "2",
			args: args{
				chess: NewChess(),
				a:     0,
				b:     2,
				x:     2,
				y:     0,
			},
			want: false,
		},
		{
			name: "3",
			args: args{
				chess: NewChess(),
				a:     2,
				b:     0,
				x:     0,
				y:     2,
			},
			want: false,
		},
		{
			name: "4",
			args: args{
				chess: NewChess(),
				a:     2,
				b:     2,
				x:     0,
				y:     0,
			},
			want: false,
		},
		{
			name: "5",
			args: args{
				chess: NewChess(),
				a:     0,
				b:     0,
				x:     0,
				y:     2,
			},
			want: false,
		},
		{
			name: "6",
			args: args{
				chess: NewChess(),
				a:     0,
				b:     2,
				x:     0,
				y:     0,
			},
			want: false,
		},
		{
			name: "7",
			args: args{
				chess: NewChess(),
				a:     2,
				b:     0,
				x:     0,
				y:     0,
			},
			want: false,
		},
		{
			name: "8",
			args: args{
				chess: NewChess(),
				a:     0,
				b:     0,
				x:     2,
				y:     0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Blockbool(tt.args.chess, tt.args.a, tt.args.b, tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Blockbool() = %v, want %v", got, tt.want)
			}
		})
	}
}
