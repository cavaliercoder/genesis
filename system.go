package genesis

import (
	"os"

	"github.com/cavaliercoder/genesis/rom"
	"github.com/cavaliercoder/genesis/vdp"
	"github.com/cavaliercoder/go-m68k"
	"github.com/cavaliercoder/go-m68k/m68kmem"
)

type System struct {
	p *m68k.Processor
	m *m68kmem.Mapper
	v *vdp.VDP
}

func New() *System {
	// configure system
	s := &System{
		p: &m68k.Processor{
			TraceWriter: os.Stdout,
		},
		m: m68kmem.NewMapper(),
		v: vdp.New(),
	}

	// attach memory mapper to the processor
	s.p.M = m68kmem.NewDecoder(s.m)

	// map devices into memory - specification is here:
	// https://en.wikibooks.org/wiki/Genesis_Programming#Memory_map

	// system ram
	s.m.Map(m68kmem.NewRAM(0x10000), 0xFF0000, 0xFFFFFF)

	// z80 address
	s.m.Map(m68kmem.NewTracer("z80", os.Stderr, m68kmem.NewNop()), 0xA00000, 0xA0FFFF)

	// version register
	s.m.Map(m68kmem.NewTracer("version", os.Stderr, m68kmem.NewROM([]byte{0x00, 0x02})), 0xA10000, 0xA10001)

	s.m.Map(m68kmem.NewTracer("c1", os.Stderr, m68kmem.NewROM([]byte{0xFF, 0xFF, 0xFF, 0xFF})), 0xA10008, 0xA10009)

	// peripherals
	// s.m.Map(m68kmem.NewRAM(16), 0xA10002, 0xA1001F)

	// z80 bus request
	s.m.Map(m68kmem.NewTracer("z80bus", os.Stderr, m68kmem.NewNop()), 0xA11100, 0xA11101)

	// z80 reset
	s.m.Map(m68kmem.NewTracer("z80reset", os.Stderr, m68kmem.NewNop()), 0xA11200, 0xA11201)

	// tmss register
	s.m.Map(m68kmem.NewTracer("tmss", os.Stderr, m68kmem.NewNop()), 0xA14000, 0xA14003)

	// vdp
	s.m.Map(m68kmem.NewTracer("vdp", os.Stderr, s.v), 0xC00000, 0xC00009)

	return s
}

func (c *System) Load(rom *rom.ROM) error {
	return c.m.Map(rom, 0, 0x3FFFFF)
}

func (c *System) Run() (err error) {
	c.p.PC = 0x200
	for {
		err = c.p.Step()
		if err != nil {
			break
		}
		// dump.Processor(os.Stderr, c.p)
	}
	return
}
