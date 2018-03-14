package rom

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/cavaliercoder/go-m68k/m68kmem"
)

// ref: https://www.zophar.net/fileuploads/2/10614uauyw/Genesis_ROM_Format.txt

var (
	ErrInvalidChecksum = errors.New("invalid checksum")
)

type ROM struct {
	Data []byte
}

func ReadROM(r io.Reader) (*ROM, error) {
	b := &bytes.Buffer{}
	_, err := io.Copy(b, r)
	if err != nil {
		return nil, err
	}

	return &ROM{
		Data: b.Bytes(),
	}, nil
}

func LoadROM(path string) (*ROM, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadROM(f)
}

func (r *ROM) String() string {
	return r.Title()
}

func (r *ROM) Read(addr int, p []byte) (n int, err error) {
	if addr < 0 || addr >= len(r.Data) {
		return 0, io.EOF
	}
	n = copy(p, r.Data[addr:])
	return
}

func (r *ROM) Write(addr int, p []byte) (n int, err error) {
	// TODO: deny write on Read-only memory
	if addr < 0 || addr >= len(r.Data) {
		return 0, m68kmem.AccessViolationError(uint32(addr))
	}
	n = copy(r.Data[addr:], p)
	if n < len(p) {
		return n, io.ErrShortWrite
	}
	return
}

func hdrString(b []byte) string {
	out := &bytes.Buffer{}
	for i := 0; i < len(b); i++ {
		if b[i] == ' ' {
			if i+1 < len(b) && b[i+1] != ' ' {
				out.WriteByte(b[i])
			}
		} else {
			out.WriteByte(b[i])
		}
	}
	return out.String()
}

func (r *ROM) Console() string {
	return hdrString(r.Data[0x100:0x110])
}

func (r *ROM) Copyright() string {
	return hdrString(r.Data[0x110:0x120])
}

func (r *ROM) Title() string {
	return r.DomesticName()
}

func (r *ROM) DomesticName() string {
	return hdrString(r.Data[0x120:0x150])
}

func (r *ROM) OverseasName() string {
	return hdrString(r.Data[0x150:0x180])
}

func (r *ROM) ProductType() string {
	return hdrString(r.Data[0x180:0x182])
}

func (r *ROM) ProductCode() string {
	return hdrString(r.Data[0x182:0x18E])
}

func (r *ROM) Checksum() uint16 {
	return uint16(r.Data[0x18E])<<8 + uint16(r.Data[0x18F])
}

func (r *ROM) Start() uint32 {
	return binary.BigEndian.Uint32(r.Data[0x1A0:0x1A4])
}

func (r *ROM) End() uint32 {
	return binary.BigEndian.Uint32(r.Data[0x1A4:0x1A8])
}

func (r *ROM) RAMStart() uint32 {
	return binary.BigEndian.Uint32(r.Data[0x1A8:0x1AC])
}

func (r *ROM) RAMEnd() uint32 {
	return binary.BigEndian.Uint32(r.Data[0x1AC:0x1B0])
}

func (r *ROM) ValidateChecksum() error {
	var sum uint16
	for i := 0x200; i < len(r.Data)-1; i += 2 {
		sum += uint16(r.Data[i])<<8 + uint16(r.Data[i+1])
	}
	if sum != r.Checksum() {
		return ErrInvalidChecksum
	}
	return nil
}

func (r *ROM) Countries() []string {
	res := make([]string, 0)
	for i := 0x1F0; i < 0x1F3; i++ {
		switch r.Data[i] {
		case 'E':
			res = append(res, "europe")

		case 'J':
			res = append(res, "japan")

		case 'A':
			res = append(res, "asia")

		case 'B', 0x04:
			res = append(res, "brazil")

		case 'F':
			res = append(res, "france")

		case 0x08:
			res = append(res, "hong kong")

		case ' ':
			// skip space
		}
	}
	return res
}
