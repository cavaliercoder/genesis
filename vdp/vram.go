package vdp

import (
	"fmt"
	"io"
)

type VRAM struct {
	b [0x20000]byte
}

func NewVRAM() *VRAM {
	return &VRAM{}
}

func (m *VRAM) Read(addr int, p []byte) (n int, err error) {
	if addr < 0 || addr >= len(m.b) {
		return 0, fmt.Errorf("access violation: 0x%X", uint32(addr))
	}
	n = copy(p, m.b[addr:])
	return
}

func (m *VRAM) Write(addr int, p []byte) (n int, err error) {
	if addr < 0 || addr >= len(m.b) {
		return 0, fmt.Errorf("access violation: 0x%X", uint32(addr))
	}
	n = copy(m.b[addr:], p)
	if n < len(p) {
		return n, io.ErrShortWrite
	}
	return
}
