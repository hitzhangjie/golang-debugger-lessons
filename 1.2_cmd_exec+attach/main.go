package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	usage = "Usage: go run main.go exec <path/to/prog>"

	cmdExec   = "exec"
	cmdAttach = "attach"
)

func main() {
	// issue: https://github.com/golang/go/issues/7699
	//
	// 为什么syscall.PtraceDetach, detach error: no such process?
	// 因为ptrace请求应该来自相同的tracer线程，
	//
	// ps: 如果恰好不是，可能需要对tracee的状态显示进行更复杂的处理，需要考虑信号？目前看系统调用传递的参数是这样
	runtime.LockOSThread()

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usage)
		os.Exit(1)
	}
	cmd := os.Args[1]

	switch cmd {
	case cmdExec:
		args := os.Args[2:]
		fmt.Printf("exec %s\n", strings.Join(args, ""))

		if len(args) != 1 {
			fmt.Println("参数错误")
			os.Exit(1)
		}

		// start process but don't wait it finished
		progCmd := exec.Command(args[0])
		progCmd.Stdin = os.Stdin
		progCmd.Stdout = os.Stdout
		progCmd.Stderr = os.Stderr
		progCmd.SysProcAttr = &syscall.SysProcAttr{
			Ptrace: true,
		}

		if err := progCmd.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// wait target process stopped
		var (
			status syscall.WaitStatus
			rusage syscall.Rusage
		)
		pid := progCmd.Process.Pid
		if _, err := syscall.Wait4(pid, &status, syscall.WALL, &rusage); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("process %d stopped:%v\n", pid, status.Stopped())

	case cmdAttach:
		pid, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s invalid pid\n\n", os.Args[2])
			os.Exit(1)
		}

		// check pid
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

		// wait
		var (
			status syscall.WaitStatus
			rusage syscall.Rusage
		)
		_, err = syscall.Wait4(int(pid), &status, syscall.WSTOPPED, &rusage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "process %d wait error: %v\n\n", pid, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "process %d wait succ, status:%v, rusage:%v\n\n", pid, status, rusage)

		// detach
		fmt.Printf("we're doing some debugging...\n")
		time.Sleep(time.Second * 10)

		// MUST: call runtime.LockOSThead() first
		err = syscall.PtraceDetach(int(pid))

		// 在没有执行lockOsThread操作之前，测试系统调用参数影响，可见对信号的处理是需要考虑的
		// 在执行lockOsThread之后，信号参数就没影响了
		//err = ptraceDetach(int(pid), 1) // sig:1,fail
		//err = ptraceDetach(int(pid), 0) // sig:0,ok
		if err != nil {
			fmt.Fprintf(os.Stderr, "process %d detach error: %v\n\n", pid, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "process %d detach succ\n\n", pid)

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
