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
	"fmt"

	"godbg/cmd/debug"

	"github.com/spf13/cobra"
)

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "调试运行中进程",
	Long:  `调试运行中进程`,
	Run: func(cmd *cobra.Command, args []string) {
		pid, _ := cmd.Flags().GetUint32("pid")
		fmt.Printf("attach to process %d\n", pid)
		debug.NewDebugShell().Run()
	},
}

func init() {
	rootCmd.AddCommand(attachCmd)

	attachCmd.Flags().Uint32P("pid", "p", 0, "process's pid to attach")
}
