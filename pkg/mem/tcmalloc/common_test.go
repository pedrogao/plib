package tcmalloc

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_mmapAnonymous(t *testing.T) {
	assert := assert.New(t)

	p, err := mmapAnonymous(100)
	assert.Nil(err)
	a := unsafe.Pointer(p)
	*(*int64)(a) = 10
	n := *(*int64)(a)
	assert.Equal(n, int64(10))
}

func Test_classFromSize(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "23",
			args: args{size: 23},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, size2Class(tt.args.size), "size2Class(%v)", tt.args.size)
		})
	}
}

func Test_classGetSize(t *testing.T) {
	type args struct {
		class int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "2",
			args: args{class: 2},
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, class2Size(tt.args.class), "class2Size(%v)", tt.args.class)
		})
	}
}
