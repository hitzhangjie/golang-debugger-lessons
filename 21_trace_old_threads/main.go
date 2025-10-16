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

	// attach to all threads for tracing
	fmt.Fprintf(os.Stdout, "===step3===: attach to all threads for tracing\n")

	attachedThreads := make(map[int]bool)

	// attach to all threads
	for _, tid := range threads {
		err := syscall.PtraceAttach(tid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "thread %d attach error: %v\n", tid, err)
			continue
		}
		attachedThreads[tid] = true
		fmt.Fprintf(os.Stdout, "thread %d attached successfully\n", tid)
	}

	fmt.Fprintf(os.Stdout, "attached to %d threads total\n\n", len(attachedThreads))

	// wait for all attached threads to stop
	fmt.Fprintf(os.Stdout, "===step4===: wait for all threads to stop\n")
	for tid := range attachedThreads {
		var status syscall.WaitStatus
		var rusage syscall.Rusage
		_, err := syscall.Wait4(tid, &status, 0, &rusage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "thread %d wait error: %v\n", tid, err)
			continue
		}
		if !status.Stopped() {
			fmt.Fprintf(os.Stderr, "thread %d not stopped\n", tid)
			continue
		}
		fmt.Fprintf(os.Stdout, "thread %d stopped\n", tid)
	}

	// show current state of all traced threads
	fmt.Fprintf(os.Stdout, "\n===step5===: show current state of all traced threads\n")
	for tid := range attachedThreads {
		regs := syscall.PtraceRegs{}
		if err := syscall.PtraceGetRegs(tid, &regs); err != nil {
			fmt.Fprintf(os.Stderr, "thread %d get regs fail: %v\n", tid, err)
			continue
		}
		fmt.Fprintf(os.Stdout, "thread %d stopped at %0x\n", tid, regs.PC())
	}

	fmt.Fprintf(os.Stdout, "\nAll threads are now being traced. Use Ctrl+C to exit.\n")
	time.Sleep(time.Minute)
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
