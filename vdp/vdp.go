package vdp

import (
	"fmt"
	"io"
)

// ref: http://md.squee.co/VDP

// VDP emulates the Video Display Processor of the Sega Genesis. Also known as
// the YM7101, derived from the Texas Instruments TMS9918A.
type VDP struct {
	R    [24]byte
	SR   uint16
	VRAM *VRAM
}

// New initializes and returns a VDP.
func New() *VDP {
	return &VDP{
		SR:   0x34FF,
		VRAM: NewVRAM(),
	}
}

func (c *VDP) Read(addr int, p []byte) (n int, err error) {
	switch addr {
	case 0x04, 0x06: // read status register
		if len(p) == 2 {
			p[0] = byte(c.SR >> 8)
			p[1] = byte(c.SR)
			n = 2
			return
		}
	}
	err = fmt.Errorf("access violation: 0x%X", uint32(addr))
	return
}

func (c *VDP) Write(addr int, p []byte) (n int, err error) {
	switch len(p) {
	case 2: // write to register
		i := int(p[0] & 0x7F)
		if i < len(c.R) {
			c.R[i] = p[1]
			n = 2
			return
		}

		// case 4: // read/write to ram
		// 	cd := p[0]>>6 | p[3]>>2
		// 	addr := uint32(p[1]) | uint32(p[0]&0x3F)<<8 | uint32(p[3])<<6
		// 	switch cd {
		// 	case 0x01: // vram read
		// 		return c.VRAM.Read(int(addr), p)
		// 	}
	}
	err = fmt.Errorf("access violation: 0x%X", uint32(addr))
	return
}

func Dump(w io.Writer, c *VDP) {
	fmt.Fprintf(w, "SR: %04X\n", c.SR)
	for i := 0; i < len(c.R); i += 2 {
		fmt.Fprintf(w, "R%02d: %02X R%02d: %02X\n", i, c.R[i], i+1, c.R[i+1])
	}
}
