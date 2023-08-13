package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const (
	usage = "Usage: go run main.go exec <path/to/prog>"

	cmdExec   = "exec"
	cmdAttach = "attach"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usage)
		os.Exit(1)
	}
	cmd := os.Args[1]

	switch cmd {
	case cmdExec:
		prog := os.Args[2]

		// check whether prog exists
		fin, err := os.Lstat(prog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lstat error: %v\n\n", err)
			os.Exit(1)
		}

		if fin.IsDir() {
			fmt.Fprintf(os.Stderr, "%s is directory\n\n", prog)
			os.Exit(1)
		}

		// check whether prog executable
		if fin.Mode()&0111 == 0 {
			fmt.Fprintf(os.Stderr, "%s not executable\n\n", prog)
			os.Exit(1)
		}

		// run prog
		progCmd := exec.Command(prog)
		buf, err := progCmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s exec error: %v, \n\n%s\n\n", err, string(buf))
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(buf))

	case cmdAttach:
		// attach to process
		pid, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s invalid pid\n\n", os.Args[2])
			os.Exit(1)
		}

		if !checkPid(int(pid)) {
			fmt.Fprintf(os.Stderr, "process %d not existed\n\n", pid)
			os.Exit(1)
		}

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
		n, err = syscall.PtracePokeData(int(pid), uintptr(regs.PC()), []byte{0xCC})
		if err != nil || n != 1 {
			fmt.Fprintf(os.Stderr, "process %d write text error: %v, n: %d\n\n", pid, err, n)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "process %d write text: %#x\n\n", pid, 0xCC)

		// after pathcing, you'll see the tracee maybe notified by SIGTRAP,
		// for example, run `top`, and `go run main.go attach <top-pid>`,
		// `top` will be notified by signal SIGTRAP and killed.
		// That's ok so far, we'll move on later.

	default:
		fmt.Fprintf(os.Stderr, "%s unknown cmd\n\n", cmd)
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
