package genesis

import (
	"os"

	"github.com/cavaliercoder/go-m68k/dump"

	"github.com/cavaliercoder/genesis/rom"
	"github.com/cavaliercoder/genesis/vdp"
	"github.com/cavaliercoder/go-m68k"
	"github.com/cavaliercoder/go-m68k/m68kmem"
)

// ref: https://emu-docs.org/Genesis/sega2f.htm

type System struct {
	p *m68k.Processor
	m *m68kmem.Mapper
	v *vdp.VDP
	w *tracer
}

// New returns a new Sega Genesis system emulator.
func New() *System {
	// configure system
	s := &System{
		p: &m68k.Processor{},
		m: m68kmem.NewMapper(),
		v: vdp.New(),
		w: newTracer(os.Stdout),
	}

	// attach memory mapper to the processor
	s.p.M = m68kmem.NewDecoder(s.m)

	// attach system tracer to processor
	s.p.TraceWriter = s.w

	// map devices into memory - specification is here:
	// https://en.wikibooks.org/wiki/Genesis_Programming#Memory_map

	// z80 address
	s.m.Map(m68kmem.NewTracer("z80", s.w, m68kmem.NewNop()), 0xA00000, 0xA0FFFF)

	// version register
	s.m.Map(m68kmem.NewTracer("version", s.w, m68kmem.NewROM([]byte{0x00, 0x02})), 0xA10000, 0xA10001)

	// controllers
	s.m.Map(m68kmem.NewTracer("ctr1", s.w, m68kmem.NewROM([]byte{0xFF, 0xFF, 0xFF, 0xFF})), 0xA10008, 0xA10009)
	s.m.Map(m68kmem.NewTracer("ctr2", s.w, m68kmem.NewROM([]byte{0xFF, 0xFF, 0xFF, 0xFF})), 0xA1000A, 0xA1000B)
	s.m.Map(m68kmem.NewTracer("ctr3", s.w, m68kmem.NewROM([]byte{0xFF, 0xFF, 0xFF, 0xFF})), 0xA1000C, 0xA1000D)

	// peripherals
	// s.m.Map(m68kmem.NewRAM(16), 0xA10002, 0xA1001F)

	// z80 bus request
	s.m.Map(m68kmem.NewTracer("z80bus", s.w, m68kmem.NewNop()), 0xA11100, 0xA11101)

	// z80 reset
	s.m.Map(m68kmem.NewTracer("z80reset", s.w, m68kmem.NewNop()), 0xA11200, 0xA11201)

	// tmss register
	s.m.Map(m68kmem.NewTracer("tmss", s.w, m68kmem.NewNop()), 0xA14000, 0xA14003)

	// vdp
	s.m.Map(m68kmem.NewTracer("vdp", s.w, s.v), 0xC00000, 0xC00009)

	// psg
	s.m.Map(m68kmem.NewTracer("psg", s.w, m68kmem.NewNop()), 0xC00011, 0xC00011)

	// system ram
	sram := m68kmem.Mirror(m68kmem.NewRAM(0x10000), 0x10000, 0x1FFFFF)
	s.m.Map(m68kmem.NewTracer("ram", s.w, sram), 0xE00000, 0xFFFFFF)
	// s.m.Map(sram, 0xE00000, 0xFFFFFF)

	return s
}

// Reset resets the state of the Genesis system, as if the Reset button on the
// console had been pushed.
func (c *System) Reset() error {
	// TODO: implement system reset
	c.w.Printf("system reset\n")
	return nil
}

// Load resets the system state before mapping the given ROM into memory and
// initializing the system ready to execute the ROM's main program.
func (c *System) Load(b *rom.ROM) (err error) {
	// reset system
	if err = c.Reset(); err != nil {
		return
	}

	// map ROM into memory
	if err = c.m.Map(b, 0, 0x3FFFFF); err != nil {
		return
	}
	c.w.Printf("rom loaded at 0x3FFFFF: %s\n", b)

	// read initial stack pointer from first word of ROM
	c.p.A[7], err = c.p.M.Long(0)
	if err != nil {
		return
	}
	c.w.Printf("stack pointer set to 0x%X\n", c.p.A[7])

	// read initial program pointer from second word of ROM
	c.p.PC, err = c.p.M.Long(4)
	if err != nil {
		return
	}
	c.w.Printf("program counter set to 0x%X\n", c.p.PC)

	return
}

// Run executes the loaded ROM until an error is encountered.
func (c *System) Run() (err error) {
	b := make([]byte, 4)
	for err == nil {
		err = c.p.Step()

		c.w.Printf("CC: %s\n", dump.FormatConditionCode(c.p.SR))
		dump.Processor(c.w, c.p)
		if c.p.A[7] != 0 {
			if _, err := c.p.M.Read(int(c.p.A[7]), b); err == nil {
				c.w.Printf("SP: %08X\n", b)
			}
		}
		c.w.Printf("\n")
	}
	return
}
