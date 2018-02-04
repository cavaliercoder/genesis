# genesis

Genesis is an experimental, work in progress emulator for the Sega Genesis
written entirely in Go.

The emulator ecosystem for the Genesis is very well established and mature.
Essentially there is no need for another contender. The purpose of this project
mainly for personal growth and to spark interest in the Go community.

Kudos to existing projects like:

- [fogleman's nes](https://github.com/fogleman/nes)
- [nwidget's nintengo](https://github.com/nwidger/nintengo)

## Design goals

- Every component should be modular and reusable
- Componant APIs must be well documented and follow established Go idioms
- Performance is sacrificed (within reason) for code simplicity
- Emulator internals should be highly observable

## Notes in code generation

My initial approach to this project was to use `go generate` to generate
discrete functions for all < 65536 possible operation codes. This would
significantly reduce the number of machine intructions and branches required to
execute any given Motorola 68000 instruction.

The approach worked, (see early commits to m68k) and performed well at runtime,
but resulted in > 30 second compile times and > 80MB binaries. This was not
conducive to a buiding a thriving developer community.
