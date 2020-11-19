package target

import "os"

// TargetProcess 被调试进程
type TargetProcess struct {
	P *os.Process
}

func ProcessStart(executable string) (*os.Process, error) {
	return nil, nil
}

func ProcessAttach(pid uint32) error {
	return nil
}
