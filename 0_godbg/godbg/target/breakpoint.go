package target

// Breakpoint 断点
type Breakpoint struct {
	Address uintptr
	Orig    byte
}

func AddBreakpoint(addr uintptr) error {
	return nil
}

func ClearBreakpoint(addr uintptr) error {
	return nil
}
