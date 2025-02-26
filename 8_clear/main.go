package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

var usage = `Usage:
	go run main.go <pid>

	args:
	- pid: specify the pid of process to attach
`

var seq int32

type Breakpoint struct {
	Seq int32

	Addr uint64
	Data byte
	// 其他信息，比如源码位置，我们将在符号级调试器开发章节补充
}

var breakpoints = map[uint64]Breakpoint{}

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

	// step2: supposing runnig break here
	time.Sleep(time.Second * 2)
	fmt.Fprintf(os.Stdout, "===step2===: supposing running `dlv> break <locspec>` here\n")

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

	// record the new breakpoint
	b := Breakpoint{
		Seq:  atomic.AddInt32(&seq, 1),
		Addr: regs.PC(),
		Data: orig[0],
	}
	breakpoints[regs.PC()] = b

	// step3: supposing running `dlv> breakpoints` here
	time.Sleep(time.Second * 2)
	fmt.Fprintf(os.Stdout, "===step3===: supposing running `dlv> breakpoints` here\n")
	for _, b := range breakpoints {
		fmt.Fprintf(os.Stdout, "breakpoint[%d] at %0x\n", b.Seq, b.Addr)
	}

	// step4: supposing running `dlv> continue` here
	time.Sleep(time.Second * 2)
	fmt.Fprintf(os.Stdout, "===step4===: supposing running `dlv> continue` here\n")
	if err := syscall.PtraceCont(int(pid), 0); err != nil {
		fmt.Fprintf(os.Stderr, "run to breakpoint fail: %v\n", err)
		os.Exit(1)
	}
	_, _ = syscall.Wait4(int(pid), nil, 0, nil)
	regs2 := syscall.PtraceRegs{}
	_ = syscall.PtraceGetRegs(int(pid), &regs2)
	fmt.Fprintf(os.Stdout, "stopped at %x\n", regs2.PC())

	// step4: supposing running `dlv> clear -n 1` to remove unused breakpoint
	time.Sleep(time.Second * 2)
	fmt.Fprintf(os.Stdout, "===step5===: supposing running `dlv> clear -n 1` here\n")
	// - 先根据断点编号查找是否存在此断点
	var bpToDelete Breakpoint
	var found bool
	for _, b := range breakpoints {
		if b.Seq == 1 {
			bpToDelete = b
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stdout, "breakpoint[%d] not found\n", 1)
		return
	}
	// - 尝试还原断点位置数据
	n, err = syscall.PtracePokeText(int(pid), uintptr(bpToDelete.Addr), []byte{bpToDelete.Data})
	if err != nil || n != 1 {
		fmt.Fprintf(os.Stderr, "clear breakpoint fail: %v\n", err)
		os.Exit(1)
	}
	// - 从断点列表中删掉
	delete(breakpoints, bpToDelete.Addr)
	// - 尝试检查当前PC位置是否刚好越过断点位置，如果是需要执行这条被patched的指令
	if err := syscall.PtraceGetRegs(int(pid), &regs); err != nil {
		fmt.Fprintf(os.Stderr, "get regs fail: %v\n", err)
		os.Exit(1)
	}
	if regs.PC()-1 == bpToDelete.Addr {
		if err := syscall.PtraceSetRegs(int(pid), &regs); err != nil {
			fmt.Fprintf(os.Stderr, "set regs (PC=PC-1) fail: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stdout, "breakpoint[1] cleared")
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
