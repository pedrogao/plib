package rt

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func Test_memmove(t *testing.T) {
	assert := assert.New(t)

	src := []byte{1, 2, 3, 4, 5, 6}

	assert.Equal(len(src), 6)
	assert.Equal(cap(src), 6)

	dest := make([]byte, 10, 10)
	assert.Equal(len(dest), 10)
	assert.Equal(cap(dest), 10)

	spew.Dump(src)
	spew.Dump(dest)

	srcp := (*GoSlice)(unsafe.Pointer(&src))
	destp := (*GoSlice)(unsafe.Pointer(&dest))

	memmove(destp.Ptr, srcp.Ptr, unsafe.Sizeof(byte(0))*6)

	t.Logf("moved\n")

	spew.Dump(src)
	spew.Dump(dest)

	str := "pedro"
	// 注意：这里的len不能为0，否则数据没有分配，就无法复制
	data := make([]byte, 10, 10)
	spew.Dump(str)
	spew.Dump(data)
	t.Logf("moved\n")
	memmove((*GoSlice)(unsafe.Pointer(&data)).Ptr, (*GoString)(unsafe.Pointer(&str)).Ptr,
		unsafe.Sizeof(byte(0))*5)
	spew.Dump(str)
	spew.Dump(data)
}

func Test_copy(t *testing.T) {
	src := []byte{1, 2, 3, 4, 5, 6}
	dest := make([]byte, 10, 10)

	spew.Dump(src)
	spew.Dump(dest)

	copy(dest, src)

	spew.Dump(src)
	spew.Dump(dest)
}

func Test_growslice(t *testing.T) {
	assert := assert.New(t)

	var typeByte = UnpackType(reflect.TypeOf(byte(0)))

	spew.Dump(typeByte)

	dest := make([]byte, 0, 10)

	assert.Equal(len(dest), 0)
	assert.Equal(cap(dest), 10)

	ds := (*GoSlice)(unsafe.Pointer(&dest))
	*ds = growslice(typeByte, *ds, 100)

	assert.Equal(len(dest), 0)
	assert.Equal(cap(dest), 112)
}
