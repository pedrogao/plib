package rt

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test___isspace(t *testing.T) {
	type args struct {
		ch byte
	}
	tests := []struct {
		name    string
		args    args
		wantRet byte
	}{
		{
			name:    "false",
			args:    args{ch: '0'},
			wantRet: 0,
		},
		{
			name:    "true",
			args:    args{ch: '\n'},
			wantRet: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := __isspace(tt.args.ch); gotRet != tt.wantRet {
				t.Errorf("__isspace() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func Test___u32toa_small(t *testing.T) {
	var buf [32]byte
	type args struct {
		out *byte
		val uint32
	}
	tests := []struct {
		name    string
		args    args
		wantRet int
	}{
		{
			name: "10",
			args: args{
				out: &buf[0],
				val: 10,
			},
			wantRet: 2,
		},
		{
			name: "9999",
			args: args{
				out: &buf[0],
				val: 9999,
			},
			wantRet: 4,
		},
		{
			name: "1234",
			args: args{
				out: &buf[0],
				val: 1234,
			},
			wantRet: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := __u32toa_small(tt.args.out, tt.args.val)
			assert.Equalf(t, tt.wantRet, got, "__u32toa_small(%v, %v)", tt.args.out, tt.args.val)
			assert.Equalf(t, tt.name, string(buf[:tt.wantRet]), "ret string must equal name")
		})
	}
}

//BenchmarkGoConv
//BenchmarkGoConv-12    	60740782	        19.52 ns/op
func BenchmarkGoConv(b *testing.B) {
	val := int(rand.Int31() % 10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		strconv.Itoa(val)
	}
}

//goos: darwin
//goarch: amd64
//cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
//BenchmarkFastConv
//BenchmarkFastConv-12    	122945924	         9.455 ns/op
func BenchmarkFastConv(b *testing.B) {
	var buf [32]byte
	val := uint32(rand.Int31() % 10000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		__u32toa_small(&buf[0], val)
	}
}
