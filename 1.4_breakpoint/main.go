package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

var usage = `Usage:
	go run main.go <pid>

	args:
	- pid: specify the pid of process to attach
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	// pid
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if !checkPid(int(pid)) {
		fmt.Fprintf(os.Stderr, "process %d not existed\n\n", pid)
		os.Exit(1)
	}

	// attach
	err = syscall.PtraceAttach(int(pid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "process %d attach error: %v\n\n", pid, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "process %d attach succ\n\n", pid)

	// check target process stopped or not
	var status syscall.WaitStatus
	var options int
	var rusage syscall.Rusage

	_, err = syscall.Wait4(int(pid), &status, options, &rusage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process %d wait error: %v\n\n", pid, err)
		os.Exit(1)
	}
	if !status.Stopped() {
		fmt.Fprintf(os.Stderr, "process %d not stopped\n\n", pid)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "process %d stopped\n\n", pid)

	// try to patch the instruction at PC with 0xCC

	// read PC
	var regs syscall.PtraceRegs
	err = syscall.PtraceGetRegs(int(pid), &regs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process %d read regs error: %v\n\n", pid, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "process %d read regs: %#x\n\n", pid, regs.PC())

	// backup original instruction data
	orig := [1]byte{}
	n, err := syscall.PtracePeekText(int(pid), uintptr(regs.PC()), orig[:])
	if err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "process %d read text error: %v, n: %d\n\n", pid, err, n)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "process %d read text: %#x\n\n", pid, orig[0])

	// patching instruction with 0xCC
	n, err = syscall.PtracePokeText(int(pid), uintptr(regs.PC()), []byte{0xCC})
	if err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "process %d write text error: %v, n: %d\n\n", pid, err, n)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "process %d write text: %#x\n\n", pid, 0xCC)
	fmt.Fprintf(os.Stdout, "add breakpoint success\n")

	// after pathcing, you'll see the tracee maybe notified by SIGTRAP,
	// for example, run `top`, and `go run main.go attach <top-pid>`,
	// `top` will be notified by signal SIGTRAP and killed.
	// That's ok so far, we'll move on later.
}

// checkPid check whether pid is valid process's id
//
// On Unix systems, os.FindProcess always succeeds and returns a Process for
// the given pid, regardless of whether the process exists.
func checkPid(pid int) bool {
	out, err := exec.Command("kill", "-s", "0", strconv.Itoa(pid)).CombinedOutput()
	if err != nil {
		panic(err)
	}

	// output error message, means pid is invalid
	if string(out) != "" {
		return false
	}

	return true
}
