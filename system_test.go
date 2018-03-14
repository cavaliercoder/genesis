package genesis

import (
	"encoding/binary"
	"testing"

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
		t.Fatal(err)
	}
}

func TestRamMirroring(t *testing.T) {
	s := New()

	// write one word at a time across entire ram range (0xE00000 - 0xFFFFFF)
	b := []byte{0, 0, 0, 0}
	for i := uint32(0xE00000); i <= 0xFFFFFF; i += 4 {
		binary.BigEndian.PutUint32(b, i)
		if _, err := s.m.Write(int(i), b); err != nil {
			t.Errorf("error writing to 0x%08X: %v", i, err)
			return
		}

		// read back from within 0xFF0000 - 0xFFFFFF
		b = []byte{0, 0, 0, 0}
		addr := 0xFF0000 + ((i - 0xE00000) % 0x10000)
		if _, err := s.m.Read(int(addr), b); err != nil {
			t.Errorf("error reading at 0x%08X: %v", addr, err)
			return
		}

		// ensure value matches
		v := binary.BigEndian.Uint32(b)
		if v != i {
			t.Errorf("error reading mirrored memory: expected %v, got %v", i, v)
			return
		}
	}
}
