package main

import (
	"fmt"
	"os"

	"github.com/cavaliercoder/genesis"
	"github.com/cavaliercoder/genesis/rom"
)

func main() {
	if len(os.Args) != 2 {
		usage(1)
	}

	s := genesis.New()
	game, err := rom.LoadROM(os.Args[1])
	dieOn(err)
	dieOn(s.Load(game))
	dieOn(s.Run())
}

func usage(code int) {
	w := os.Stdout
	if code != 0 {
		w = os.Stderr
	}
	fmt.Fprintf(w, "usage: %s FILE\n", os.Args[0])
	os.Exit(code)
}

func dieOn(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
