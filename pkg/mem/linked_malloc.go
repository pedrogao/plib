package mem

import (
	"fmt"
	"syscall"
	"unsafe"
)

type Area struct {
	addr      uintptr // 首地址
	mem       []byte  // 内存
	size      int     // 内存大小
	allocated int     // 已分配内存
	blocks    int     // 内存块个数
}

// TODO 暂时不关注16字节对齐，把注意力放在核心实现上
type header struct {
	size      int
	allocated bool
	next      *header
}

var (
	placeholder = [16]byte{0}
	headerSize  = int(unsafe.Sizeof(header{}))
)

func Init(size int) (*Area, error) {
	mem, err := syscall.Mmap(-1, 0, size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		return nil, fmt.Errorf("mmap err: %s", err)
	}

	addr := unsafe.Pointer(&mem[0])
	hdr := (*header)(addr)
	hdr.size = size - headerSize
	m := &Area{
		addr:      uintptr(addr),
		mem:       mem,
		size:      size,
		allocated: 0,
		blocks:    0,
	}

	return m, nil
}

func (m *Area) Alloc(size int) (unsafe.Pointer, error) {
	if m.allocated+size+headerSize > m.size {
		return nil, fmt.Errorf("can't alloc any more")
	}

	h := (*header)(unsafe.Pointer(m.addr))
	// 如果内存已分配，且大小不够，则下一个
	for h.allocated || h.size < size {
		h = h.next
	}
	// 没有找到可用的 block
	if h == nil {
		return nil, fmt.Errorf("can't alloc any more")
	}

	prevSz := h.size
	cur := unsafe.Pointer(h)
	p := unsafe.Pointer(uintptr(cur) + uintptr(headerSize))
	// 下一个
	hdr := (*header)(unsafe.Pointer(uintptr(cur) + uintptr(headerSize) + uintptr(size)))
	hdr.size = prevSz - size - headerSize
	h.size = size
	h.allocated = true
	h.next = hdr
	m.allocated += size + headerSize
	m.blocks++
	return p, nil
}

func (m *Area) Free(p unsafe.Pointer) error {
	h := (*header)(unsafe.Pointer(uintptr(p) - uintptr(headerSize)))
	if !h.allocated {
		return fmt.Errorf("can't free %x twice", p)
	}
	h.allocated = false
	m.allocated -= h.size + headerSize
	m.blocks--
	return nil
}

func (m *Area) Reset() {
	m.blocks = 0
	m.allocated = 0
	for i := 0; i < m.size; i++ {
		m.mem[i] = 0
	}
}

func (m *Area) Release() error {
	return syscall.Munmap(m.mem)
}

func (m *Area) Stat() string {
	return fmt.Sprintf("{ total: %d, allocated: %d, remain: %d, blocks: %d }",
		m.size, m.allocated, m.size-m.allocated, m.blocks)
}
