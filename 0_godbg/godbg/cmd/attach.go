/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"

	"godbg/cmd/debug"

	"github.com/spf13/cobra"
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:   "attach <traceePID>",
	Short: "调试运行中进程",
	Long:  `调试运行中进程`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//fmt.Printf("attach %s\n", strings.Join(args, ""))

		// issue: https://github.com/golang/go/issues/7699
		//
		// 为什么syscall.PtraceDetach, detach error: no such process?
		// 因为ptrace请求应该来自相同的tracer线程，
		//
		// ps: 如果恰好不是，可能需要对tracee的状态显示进行更复杂的处理，需要考虑信号？目前看系统调用传递的参数是这样
		runtime.LockOSThread()

		if len(args) != 1 {
			return errors.New("参数错误")
		}

		pid, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s invalid traceePID\n\n", os.Args[2])
			os.Exit(1)
		}
		traceePID = int(pid)

		// check traceePID
		if !checkPid(int(pid)) {
			fmt.Fprintf(os.Stderr, "process %d not existed\n\n", pid)
			os.Exit(1)
		}

		// attach
		err = syscall.PtraceAttach(int(pid))
		if err != nil {
			fmt.Printf("process %d attached error: %v\n", pid, err)
			os.Exit(1)
		}
		fmt.Printf("process %d attached succ\n", pid)

		// wait
		var (
			status syscall.WaitStatus
			rusage syscall.Rusage
		)
		_, err = syscall.Wait4(int(pid), &status, syscall.WSTOPPED, &rusage)
		if err != nil {
			fmt.Printf("process %d waited error: %v\n", pid, err)
			os.Exit(1)
		}
		fmt.Printf("process %d stopped: %v\n", pid, status.Stopped())
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		debug.NewDebugShell().Run()

		// MUST: call runtime.LockOSThead() first
		return syscall.PtraceDetach(int(traceePID))
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)

	attachCmd.Flags().Uint32P("traceePID", "p", 0, "process's traceePID to attach")
}

// checkPid check whether traceePID is valid process's id
//
// On Unix systems, os.FindProcess always succeeds and returns a Process for
// the given traceePID, regardless of whether the process exists.
func checkPid(pid int) bool {
	out, err := exec.Command("kill", "-s", "0", strconv.Itoa(pid)).CombinedOutput()
	if err != nil {
		panic(err)
	}

	// output error message, means traceePID is invalid
	if string(out) != "" {
		return false
	}

	return true
}
