package common

import (
	"reflect"
	"testing"
)

func TestString2Bytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "simple",
			args: args{
				s: "pedro",
			},
			want: []byte("pedro"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String2Bytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String2Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes2String(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple",
			args: args{
				b: []byte("mike"),
			},
			want: "mike",
		},
		{
			name: "mul",
			args: args{
				b: []byte("0c2Y43LKAMtoyF4xgLxkTki0zCIzKPZOXKrMAiQoPB8OCLlQ4XmyeZucoI9yd7o0GXrYD3jip8N7cN5wb1xyTqn6oVPPnpQYI8XlVOMu9gOUSkCYCFlPUewB"),
			},
			want: "0c2Y43LKAMtoyF4xgLxkTki0zCIzKPZOXKrMAiQoPB8OCLlQ4XmyeZucoI9yd7o0GXrYD3jip8N7cN5wb1xyTqn6oVPPnpQYI8XlVOMu9gOUSkCYCFlPUewB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes2String(tt.args.b); got != tt.want {
				t.Errorf("Bytes2String() = %v, want %v", got, tt.want)
			}
		})
	}
}
