package storage

import (
	"github.com/kenelite/gedis/internal/response"
	"testing"
)

func TestNewAof(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Aof
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				path: "./database.aof.test",
			},
			want:    &Aof{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAof(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestAof_Write(t *testing.T) {
	type args struct {
		path  string
		value response.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				path: "./database.aof.test",
				value: response.Value{
					Typ: "string",
					Str: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aof, _ := NewAof(tt.args.path)
			if err := aof.Write(tt.args.value); err != nil {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
