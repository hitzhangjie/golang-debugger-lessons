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

	fmt.Fprintf(os.Stdout, "===step1===: check target process existed or not\n")
	// pid
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if !checkPid(int(pid)) {
		fmt.Fprintf(os.Stderr, "process %d not existed\n\n", pid)
		os.Exit(1)
	}

	// enumerate all threads
	fmt.Fprintf(os.Stdout, "===step2===: enumerate created threads by reading /proc\n")

	// read dir entries of /proc/<pid>/task/
	threads, err := readThreadIDs(pid)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "threads: %v\n", threads)

	// prompt user which thread to attach
	var last int64

	// attach thread <n>, or switch thread to another one thread <m>
	for {
		fmt.Fprintf(os.Stdout, "===step3===: supposing running `dlv> thread <n>` here\n")
		var target int64
		n, err := fmt.Fscanf(os.Stdin, "%d\n", &target)
		if n == 0 || err != nil || target <= 0 {
			panic("invalid input, thread id should > 0")
		}

		if last > 0 {
			if err := syscall.PtraceDetach(int(last)); err != nil {
				fmt.Fprintf(os.Stderr, "switch from thread %d to thread %d error: %v\n", last, target, err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "switch from thread %d thread %d\n", last, target)
		}

		// attach
		err = syscall.PtraceAttach(int(target))
		if err != nil {
			fmt.Fprintf(os.Stderr, "thread %d attach error: %v\n\n", target, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "process %d attach succ\n\n", target)

		// check target process stopped or not
		var status syscall.WaitStatus
		var rusage syscall.Rusage
		_, err = syscall.Wait4(int(target), &status, 0, &rusage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "process %d wait error: %v\n\n", target, err)
			os.Exit(1)
		}
		if !status.Stopped() {
			fmt.Fprintf(os.Stderr, "process %d not stopped\n\n", target)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "process %d stopped\n\n", target)

		regs := syscall.PtraceRegs{}
		if err := syscall.PtraceGetRegs(int(target), &regs); err != nil {
			fmt.Fprintf(os.Stderr, "get regs fail: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "tracee stopped at %0x\n", regs.PC())

		last = target
		time.Sleep(time.Second)
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

// reads all thread IDs associated with a given process ID.
func readThreadIDs(pid int) ([]int, error) {
	dir := fmt.Sprintf("/proc/%d/task", pid)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var threads []int
	for _, file := range files {
		tid, err := strconv.Atoi(file.Name())
		if err != nil { // Ensure that it's a valid positive integer
			continue
		}
		threads = append(threads, tid)
	}
	return threads, nil
}
