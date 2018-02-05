package genesis

import (
	"os"
	"testing"

	"github.com/cavaliercoder/genesis/vdp"

	"github.com/cavaliercoder/genesis/rom"
	"github.com/cavaliercoder/go-m68k/dump"
)

func TestGame(t *testing.T) {
	s := New()
	game, err := rom.LoadROM("roms/barts_nightmare.bin")
	if err != nil {
		t.Fatal(err)
	}

	err = s.Load(game)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		vdp.Dump(os.Stderr, s.v)
		dump.Memory(os.Stderr, s.v.VRAM)
		t.Fatal(err)
	}
}
