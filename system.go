package genesis

import (
	"os"

	"github.com/cavaliercoder/genesis/rom"
	"github.com/cavaliercoder/go-m68k"
	"github.com/cavaliercoder/go-m68k/m68kmem"
)

type System struct {
	p *m68k.Processor
	m *m68kmem.Mapper
}

func New() *System {
	// map devices into memory - specification is here:
	// https://en.wikibooks.org/wiki/Genesis_Programming#Memory_map
	mm := m68kmem.NewMapper()

	// system ram
	mm.Map(m68kmem.NewRAM(0x10000), 0xFF0000, 0xFFFFFF)

	// peripherals
	mm.Map(m68kmem.NewRAM(16), 0xA10000, 0xA1001F)

	// vdp
	mm.Map(m68kmem.NewNop(), 0xC00000, 0xC00010)

	// configure processor
	p := &m68k.Processor{
		M:           m68kmem.NewDecoder(mm),
		TraceWriter: os.Stdout,
	}
	return &System{p: p, m: mm}
}

func (c *System) Load(rom *rom.ROM) error {
	return c.m.Map(rom, 0, 0x3FFFFF)
}

func (c *System) Run() error {
	c.p.PC = 0x200
	return c.p.Run()
}
