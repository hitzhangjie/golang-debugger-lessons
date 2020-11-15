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
		prog := os.Args[2]

		// run prog
		progCmd := exec.Command(prog)
		buf, err := progCmd.CombinedOutput()

		fmt.Fprintf(os.Stdout, "tracee pid: %d\n", progCmd.Process.Pid)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s exec error: %v, \n\n%s\n\n", err, string(buf))
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(buf))

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

// 如果在ptrace之前不先调用runtime.LockOSThread()，直接使用syscall.PtraceDetach(pid)会报错，提示no such process，
// 但是这个函数却可以成功，对比了下，走的都是相同的系统调用syscall.Syscall6，区别是有一个参数不同，但是我不知道这个参数的具体作用是什么
//
// 在没有执行lockOsThread操作之前，测试系统调用参数影响，可见对信号的处理是需要考虑的
// 在执行lockOsThread之后，信号参数就没影响了
//err = ptraceDetach(int(pid), 1) // sig:1,fail, detach失败: no such process
//err = ptraceDetach(int(pid), 0) // sig:0,ok, detach成功
//
// 在没有os.LockOSThread()之前，
// ptraceDetach(pid, 0)，ok，可以正常detach，但是ptraceDetach(pid, 1)不行，这里的第2个参数对应syscall.Syscall6里面的一个参数，
// 暂时不清楚对应参数的含义，没有明确说明！看上去像是控制信号的，tracee stop状态和signal stop状态有区别，需要显示处理。
//
// 参考man手册里面的说明，这里先了解下这个问题好了，后面再仔细研究下。
/*
Detaching of the tracee is performed by:

ptrace(PTRACE_DETACH, pid, 0, sig);

PTRACE_DETACH is a restarting operation; therefore it requires the
tracee to be in ptrace-stop.  If the tracee is in signal-delivery-
stop, a signal can be injected.  Otherwise, the sig parameter may be
silently ignored.

If the tracee is running when the tracer wants to detach it, the
usual solution is to send SIGSTOP (using tgkill(2), to make sure it
goes to the correct thread), wait for the tracee to stop in signal-
delivery-stop for SIGSTOP and then detach it (suppressing SIGSTOP in‐
jection).  A design bug is that this can race with concurrent
SIGSTOPs.  Another complication is that the tracee may enter other
ptrace-stops and needs to be restarted and waited for again, until
SIGSTOP is seen.  Yet another complication is to be sure that the
tracee is not already ptrace-stopped, because no signal delivery hap‐
pens while it is—not even SIGSTOP.

If the tracer dies, all tracees are automatically detached and
restarted, unless they were in group-stop.  Handling of restart from
group-stop is currently buggy, but the "as planned" behavior is to
leave tracee stopped and waiting for SIGCONT.  If the tracee is
restarted from signal-delivery-stop, the pending signal is injected.
*/
func ptraceDetach(tid, sig int) error {
	_, _, err := syscall.Syscall6(syscall.SYS_PTRACE, syscall.PTRACE_DETACH, uintptr(tid), 1, uintptr(sig), 0, 0)
	if err != syscall.Errno(0) {
		return err
	}
	return nil
}
