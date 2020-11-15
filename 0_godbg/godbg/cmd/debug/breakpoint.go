package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "breakpoint",
	Short: "在源码中添加断点",
	Long: `在源码中添加断点，源码位置可以通过locspec格式指定。

当前支持的locspec格式，包括两种:
- [文件名:]行号
- [文件名:]函数名`,
	Aliases: []string{"b", "break"},
	Annotations: map[string]string{
		commandGroup: "breakpoint",
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("breakpoint command called")
	},
}

var clearCmd = &cobra.Command{
	Use:     "clear <n>",
	Short:   "清除指定编号的断点",
	Long:    `清除指定编号的断点`,
	Aliases: []string{"b", "break"},
	Annotations: map[string]string{
		commandGroup: "breakpoint",
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("breakpoint command called")
	},
}

var listFunctions = &cobra.Command{
	Use:   "food",
	Short: "Get some food",
	Annotations: map[string]string{
		commandGroup: "breakpoint",
	},
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		for _, v := range args {
			if verbose {
				fmt.Println("Here you go, take this:", v)
			} else {
				fmt.Println(v)
			}
		}
	},
}

func init() {
	breakCmd.AddCommand(listFunctions)
	breakCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose log")
}
