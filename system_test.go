package genesis

import (
	"os"
	"testing"

	"github.com/cavaliercoder/go-m68k/dump"

	"github.com/cavaliercoder/genesis/rom"
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
		dump.Processor(os.Stderr, s.p)
		t.Fatal(err)
	}
}
