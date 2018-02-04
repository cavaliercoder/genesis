package rom

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func TestRom(t *testing.T) {
	dir := "../roms/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".bin") {
			// t.Logf("Loading %s...", file.Name())
			rom, err := LoadROM(dir + file.Name())
			if err != nil {
				t.Fatalf("%v", err)
			}

			if err := rom.ValidateChecksum(); err != nil {
				t.Errorf("%s [%v:%v]: %v", file.Name(), rom.Start(), rom.End(), err)
			}
		}
	}
}
