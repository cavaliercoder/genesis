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

	// z80
	s.m.Map(m68kmem.NewNop(), 0xA00000, 0xA0FFFF)

	// version register
	s.m.Map(m68kmem.NewNop(), 0xA10000, 0xA10001)

	// peripherals
	s.m.Map(m68kmem.NewRAM(16), 0xA10002, 0xA1001F)

	// TMSS register
	s.m.Map(m68kmem.NewNop(), 0xA14000, 0xA14003)

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
