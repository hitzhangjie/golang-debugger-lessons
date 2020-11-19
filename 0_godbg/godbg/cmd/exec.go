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
	"strings"
	"syscall"

	"godbg/cmd/debug"

	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var pid int

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <prog>",
	Short: "调试可执行程序",
	Long:  `调试可执行程序`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("exec %s\n", strings.Join(args, ""))

		if len(args) != 1 {
			return errors.New("参数错误")
		}

		// start process but don't wait it finished
		progCmd := exec.Command(args[0])
		progCmd.Stdin = os.Stdin
		progCmd.Stdout = os.Stdout
		progCmd.Stderr = os.Stderr
		progCmd.SysProcAttr = &syscall.SysProcAttr{
			Ptrace: true,
			Setpgid: true,
			Foreground: true,
		}

		err := progCmd.Start()
		if err != nil {
			return err
		}

		// attach target process
		pid = progCmd.Process.Pid
		err = syscall.PtraceAttach(pid)
		if err != nil {
			return err
		}
		fmt.Printf("process %d attach succ\n", pid)

		// wait target process stopped
		var (
			status syscall.WaitStatus
			rusage syscall.Rusage
		)
		_, err = syscall.Wait4(int(pid), &status, syscall.WSTOPPED, &rusage)
		if err != nil {
			return err
		}
		fmt.Printf("process %d stoppped\n", pid)

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		debug.NewDebugShell().Run()
		return syscall.PtraceDetach(pid)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}

func attachProcessToTTY(process *exec.Cmd, tty string) (*os.File, error) {
	f, err := os.OpenFile(tty, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	if !isatty.IsTerminal(f.Fd()) {
		f.Close()
		return nil, fmt.Errorf("%s is not a terminal", f.Name())
	}
	process.Stdin = f
	process.Stdout = f
	process.Stderr = f
	process.SysProcAttr.Setpgid = false
	process.SysProcAttr.Setsid = true
	process.SysProcAttr.Setctty = true

	return f, nil
}
