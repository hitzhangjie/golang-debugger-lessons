package target

import (
	"go.uber.org/atomic"
)

var (
	seqNo atomic.Int64
)

// Breakpoint 断点
type Breakpoint struct {
	ID       int64
	Addr     uintptr
	Orig     byte
	Location string
}

type Breakpoints []Breakpoint

func (b Breakpoints) Len() int {
	return len(b)
}

func (b Breakpoints) Less(i, j int) bool {
	if b[i].ID <= b[j].ID {
		return true
	}
	return false
}

func (b Breakpoints) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func NewBreakpoint(addr uintptr, orig byte, location string) (Breakpoint, error) {
	b := Breakpoint{
		ID:       seqNo.Add(1),
		Addr:     0,
		Orig:     0,
		Location: "",
	}
	return b, nil
}

func AddBreakpoint(addr uintptr) error {
	return nil
}

func ClearBreakpoint(ddr uintptr) error {
	return nil
}
