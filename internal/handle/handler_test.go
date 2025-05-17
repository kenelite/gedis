package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"reflect"
	"testing"
)

func Test_get(t *testing.T) {
	type args struct {
		args []response.Value
	}
	tests := []struct {
		name string
		args args
		want response.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hget(t *testing.T) {
	type args struct {
		args []response.Value
	}
	tests := []struct {
		name string
		args args
		want response.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hget(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hgetall(t *testing.T) {
	type args struct {
		args []response.Value
	}
	tests := []struct {
		name string
		args args
		want response.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hgetall(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hgetall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hset(t *testing.T) {
	type args struct {
		args []response.Value
	}
	tests := []struct {
		name string
		args args
		want response.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hset(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ping(t *testing.T) {
	type args struct {
		args []response.Value
	}
	tests := []struct {
		name string
		args args
		want response.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ping(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set(t *testing.T) {
	type args struct {
		args []response.Value
	}
	tests := []struct {
		name string
		args args
		want response.Value
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("set() = %v, want %v", got, tt.want)
			}
		})
	}
}
