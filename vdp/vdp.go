package vdp

import (
	"fmt"
	"io"
)

// ref: http://md.squee.co/VDP
// ref: https://emu-docs.org/Genesis/Graphics/genvdp.txt
// ref: https://megacatstudios.com/blogs/press/sega-genesis-mega-drive-vdp-graphics-guide-v1-2a-03-14-17

// VDP emulates the Video Display Processor of the Sega Genesis. Also known as
// the YM7101, derived from the Texas Instruments TMS9918A.
type VDP struct {
	R    [24]byte
	SR   uint16
	VRAM *VRAM

	A  uint16 // address register
	CD uint16 // code register
	WP bool   // write pending flag
}

// New initializes and returns a VDP.
func New() *VDP {
	return &VDP{
		SR:   0x3400,
		VRAM: NewVRAM(),
	}
}

func (c *VDP) Read(addr int, p []byte) (n int, err error) {
	switch addr {
	case 0x04, 0x06: // read status register
		if len(p) == 2 {
			c.WP = false
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
	// TODO: support single byte writes
	var nn int
	for i := addr; len(p) > 0; i += 2 {
		switch addr {
		case 0x00, 0x02:
			nn, err = c.writeData(p[:2])

		case 0x04, 0x06:
			nn, err = c.writeControl(p[:2])

		default:
			err = fmt.Errorf("vdp access violation: 0x%X", uint32(addr))
		}
		if err != nil {
			return
		}
		n += nn
		addr += nn
		p = p[2:]
	}
	return
}

// writeData writes two bytes to the data port.
func (c *VDP) writeData(p []byte) (n int, err error) {
	if len(p) > 2 {
		panic("data port is only two bytes wide")
	}
	if len(p) == 1 {
		// single-byte writes are interpreted as two bytes with the same value.
		p = []byte{p[0], p[0]}
	}
	if c.CD&0x01 == 0 {
		err = fmt.Errorf("vdp data read operation expected")
		return
	}

	c.WP = false
	switch c.CD & 0x0E {
	case 0: // VRAM
		return c.VRAM.Write(int(c.A), p)
	}
	return
}

// writeControl writes two bytes to the Control Port.
func (c *VDP) writeControl(p []byte) (n int, err error) {
	if len(p) > 2 {
		panic("control port is only two bytes wide")
	}
	if len(p) == 1 {
		// single-byte writes are interpreted as two bytes with the same value.
		p = []byte{p[0], p[0]}
	}
	n = 2

	if c.WP {
		// set second half of memory access mode
		c.CD = (c.CD & 0x03) | uint16(p[1])>>2
		c.A = (c.A & 0x3F) | uint16(p[1]&0x03)<<14
		c.WP = false
		return
	}

	// if bits 14 - 15 are 1 0, write to register
	if p[0]&0xC0 == 0x80 {
		i := int(p[0] & 0x3F)
		if i < len(c.R) {
			c.R[i] = p[1]
		}
		return
	}

	// set first half of memory access mode
	c.CD = (c.CD & 0xFC) | uint16(p[0])>>6
	c.A = (c.A & 0xC000) | uint16(p[0]&0x3F)<<8 | uint16(p[1])
	c.WP = true
	return
}

// Dump prints the current state of the VDP to the given io.Writer.
func Dump(w io.Writer, c *VDP) {
	fmt.Fprintf(w, "SR: %04X\n", c.SR)
	for i := 0; i < len(c.R); i += 2 {
		fmt.Fprintf(w, "R%02d: %02X R%02d: %02X\n", i, c.R[i], i+1, c.R[i+1])
	}
}
