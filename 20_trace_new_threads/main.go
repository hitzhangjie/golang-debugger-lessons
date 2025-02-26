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
	var rusage syscall.Rusage
	_, err = syscall.Wait4(int(pid), &status, 0, &rusage)
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

	// step2: setup to trace all new threads creation events
	time.Sleep(time.Second * 2)

	opts := syscall.PTRACE_O_TRACEFORK | syscall.PTRACE_O_TRACEVFORK | syscall.PTRACE_O_TRACECLONE
	if err := syscall.PtraceSetOptions(int(pid), opts); err != nil {
		fmt.Fprintf(os.Stderr, "set options fail: %v\n", err)
		os.Exit(1)
	}

	for {
		// 放行主线程，因为每次主线程都会因为命中clone就停下来
		if err := syscall.PtraceCont(int(pid), 0); err != nil {
			fmt.Fprintf(os.Stderr, "cont fail: %v\n", err)
			os.Exit(1)
		}

		// 检查主线程状态，检查如果status是clone事件，则继续获取clone出的线程的lwp pid
		var status syscall.WaitStatus
		rusage := syscall.Rusage{}
		_, err := syscall.Wait4(pid, &status, syscall.WSTOPPED|syscall.WCLONE, &rusage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wait4 fail: %v\n", err)
			break
		}
		// 检查下状态信息是否是clone事件 (see `man 2 ptrace` 关于选项PTRACE_O_TRACECLONE的说明部分)
		isclone := status>>8 == (syscall.WaitStatus(syscall.SIGTRAP) | syscall.WaitStatus(syscall.PTRACE_EVENT_CLONE<<8))
		fmt.Fprintf(os.Stdout, "tracee stopped, tracee pid:%d, status: %s, trapcause is clone: %v\n",
			pid,
			status.StopSignal().String(),
			isclone)

		// 获取子线程对应的lwp的pid
		msg, err := syscall.PtraceGetEventMsg(int(pid))
		if err != nil {
			fmt.Fprintf(os.Stderr, "get event msg fail: %v\n", err)
			break
		}
		fmt.Fprintf(os.Stdout, "eventmsg: new thread lwp pid: %d\n", msg)

		// 放行子线程继续执行
		_ = syscall.PtraceDetach(int(msg))

		time.Sleep(time.Second * 2)
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
