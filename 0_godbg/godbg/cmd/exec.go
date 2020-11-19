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

	"godbg/cmd/debug"

	"github.com/spf13/cobra"
)

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

		progCmd := exec.Command(args[0])
		buf, err := progCmd.CombinedOutput()

		fmt.Fprintf(os.Stdout, "tracee pid: %d\n", progCmd.Process.Pid)

		if err != nil {
			return fmt.Errorf("%s exec error: %v, \n\n%s\n\n", err, string(buf))
		}
		fmt.Printf("%s\n", string(buf))
		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		debug.NewDebugShell().Run()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
