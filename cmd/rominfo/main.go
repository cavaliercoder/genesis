package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/cavaliercoder/genesis/rom"
)

const romFormat = `Title       : {{ .Title }}
Copyright   : {{ .Copyright }}
Console     : {{ .Console }}

`

var romTemplate = func() *template.Template {
	tmpl, err := template.New("rom").
		Funcs(template.FuncMap{
			"join": func(a []string) string {
				return strings.Join(a, ", ")
			},
		}).
		Parse(romFormat)
	if err != nil {
		panic(err)
	}
	return tmpl
}()

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("--- %s ---\n", os.Args[i])
		r, err := rom.LoadROM(os.Args[i])
		if err != nil {
			panic(err)
		}
		PrintROM(r)
	}
}

func PrintROM(rom *rom.ROM) {
	romTemplate.Execute(os.Stdout, rom)
}
