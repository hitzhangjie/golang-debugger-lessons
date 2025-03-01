package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var usage = `Usage:
	go run main.go <pid>

	args:
	- pid: specify the pid of process to attach
`

func main() {
	runtime.LockOSThread()

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

	// step1: supposing running dlv attach here
	fmt.Fprintf(os.Stdout, "===step1===: supposing running `dlv attach pid` here\n")

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

	regs := syscall.PtraceRegs{}
	if err := syscall.PtraceGetRegs(int(pid), &regs); err != nil {
		fmt.Fprintf(os.Stderr, "get regs fail: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "tracee stopped at %0x\n", regs.PC())

	// step2: supposing running `dlv> b <addr>`  and `dlv> continue` here
	time.Sleep(time.Second * 2)
	fmt.Fprintf(os.Stdout, "===step2===: supposing running `dlv> b <addr>`  and `dlv> continue` here\n")

	// read the address
	var input string
	fmt.Fprintf(os.Stdout, "enter return address of loop()\n")
	_, err = fmt.Fscanf(os.Stdin, "%s", &input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read address fail\n")
		os.Exit(1)
	}
	addr, err := strconv.ParseUint(input, 0, 64)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "you entered %0x\n", addr)

	// add breakpoint and run there
	var orig [1]byte
	if n, err := syscall.PtracePeekText(int(pid), uintptr(addr), orig[:]); err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "peek text fail, n: %d, err: %v\n", n, err)
		os.Exit(1)
	}
	if n, err := syscall.PtracePokeText(int(pid), uintptr(addr), []byte{0xCC}); err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "poke text fail, n: %d, err: %v\n", n, err)
		os.Exit(1)
	}
	if err := syscall.PtraceCont(int(pid), 0); err != nil {
		fmt.Fprintf(os.Stderr, "ptrace cont fail, err: %v\n", err)
		os.Exit(1)
	}

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

	// step3: supposing change register RAX value from true to false
	time.Sleep(time.Second * 2)
	fmt.Fprintf(os.Stdout, "===step3===: supposing change register RAX value from true to false\n")
	if err := syscall.PtraceGetRegs(int(pid), &regs); err != nil {
		fmt.Fprintf(os.Stderr, "ptrace get regs fail, err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "before RAX=%x\n", regs.Rax)

	regs.Rax &= 0xffffffff00000000
	if err := syscall.PtraceSetRegs(int(pid), &regs); err != nil {
		fmt.Fprintf(os.Stderr, "ptrace set regs fail, err: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "after RAX=%x\n", regs.Rax)

	// step4: let tracee continue and check it behavior (loop3.go should exit the for-loop)
	if n, err := syscall.PtracePokeText(int(pid), uintptr(addr), orig[:]); err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "restore instruction data fail: %v\n", err)
		os.Exit(1)
	}
	if err := syscall.PtraceCont(int(pid), 0); err != nil {
		fmt.Fprintf(os.Stderr, "ptrace cont fail, err: %v\n", err)
		os.Exit(1)
	}
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
