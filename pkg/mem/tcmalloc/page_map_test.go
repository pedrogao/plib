package tcmalloc

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_pageMap_operations(t *testing.T) {
	assert := assert.New(t)

	pm := newPageMap()
	addr, err := mmapAnonymous(pageSize)
	*(*int64)(unsafe.Pointer(addr)) = 1024
	assert.Nil(err)
	err = pm.insertAt(1, addr)
	assert.Nil(err)
	p := pm.get(1)
	n := (*int64)(unsafe.Pointer(p))
	assert.Equal(*n, int64(1024))

	type person struct {
		age int
	}
	ps := &person{age: 26}
	err = pm.insertAt(2, uintptr(unsafe.Pointer(ps)))
	assert.Nil(err)
	p = pm.get(2)
	pp := (*person)(unsafe.Pointer(p))
	assert.Equal(pp.age, 26)
}

func Test_getIndex(t *testing.T) {
	type args struct {
		pageNum int
		level   int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1 1",
			args: args{
				pageNum: 1,
				level:   1,
			},
			want: 0,
		},
		{
			name: "1 2",
			args: args{
				pageNum: 1,
				level:   2,
			},
			want: 0,
		},
		{
			name: "1 3",
			args: args{
				pageNum: 1,
				level:   3,
			},
			want: 1,
		},
		{
			name: "262144 1",
			args: args{
				pageNum: 262144,
				level:   1,
			},
			want: 0,
		},
		{
			name: "262144 2",
			args: args{
				pageNum: 262144,
				level:   2,
			},
			want: 1,
		},
		{
			name: "262144 3",
			args: args{
				pageNum: 262144,
				level:   3,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getIndex(tt.args.pageNum, tt.args.level), "getIndex(%v, %v)", tt.args.pageNum, tt.args.level)
		})
	}
}
