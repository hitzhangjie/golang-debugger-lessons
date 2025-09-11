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

	// step2: supposing running list and disass <locspec> go get the address of interested code
	time.Sleep(time.Second * 2)

	var input string
	fmt.Fprintf(os.Stdout, "enter a address you want to modify data\n")
	_, err = fmt.Fscanf(os.Stdin, "%s", &input)
	if err != nil {
		panic("read address fail")
	}
	addr, err := strconv.ParseUint(input, 0, 64)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "you entered %0x\n", addr)

	fmt.Fprintf(os.Stdout, "enter a value you want to change to\n")
	_, err = fmt.Fscanf(os.Stdin, "%s", &input)
	if err != nil {
		panic("read value fail")
	}
	val, err := strconv.ParseUint(input, 0, 64)
	if err != nil {
		panic("read value fail")
	}
	fmt.Fprintf(os.Stdout, "you entered %x\n", val)
	fmt.Fprintf(os.Stdout, "we'll set *(%x) = %x\n", addr, val)

	// step2: supposing runnig step here
	time.Sleep(time.Second * 2)
	fmt.Fprintf(os.Stdout, "===step2===: supposing running `dlv> set *addr = 0xaf` here\n")

	var data [1]byte
	n, err := syscall.PtracePeekData(int(pid), uintptr(addr), data[:])
	if err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "read data fail: %v\n", err)
		os.Exit(1)
	}

	n, err = syscall.PtracePokeData(int(pid), uintptr(addr), []byte{byte(val)})
	if err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "write data fail: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "change data from %x to %d succ\n", data[0], val)
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
