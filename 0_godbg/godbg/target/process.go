package target

import (
	"fmt"
	"os"
)

// TargetProcess 被调试进程
type TargetProcess struct {
	Process     *os.Process
	Breakpoints Breakpoints
}

func (t *TargetProcess) ListBreakpoints() {
	for _, b := range t.Breakpoints {
		fmt.Printf("breakpoint[%d] addr:%#x, loc:%s\n", b.ID, b.Addr, b.Location)
	}
}

func (t *TargetProcess) AddBreakpoint(addr uintptr, loc string) {

}

func ProcessStart(executable string) (*os.Process, error) {
	return nil, nil
}

func ProcessAttach(pid uint32) error {
	return nil
}
