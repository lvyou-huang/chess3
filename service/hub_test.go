package service

import (
	"chess/model"
	"reflect"
	"testing"
)

func TestNewHub(t *testing.T) {
	tests := []struct {
		name string
		want *model.Hub
	}{
		// TODO: Add test cases.
		{name: "1",
			want: &model.Hub{Broadcast: make(chan string),
				Clients:    make(map[*model.Client]struct{}),
				Register:   make(chan *model.Client),
				Unregister: make(chan *model.Client)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHub(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		hub model.Hub
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{hub: *NewHub()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.hub)
		})
	}
}
