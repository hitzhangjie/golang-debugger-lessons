package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break <locspec>",
	Short: "在源码中添加断点",
	Long: `在源码中添加断点，源码位置可以通过locspec格式指定。

当前支持的locspec格式，包括两种:
- [文件名:]行号
- [文件名:]函数名`,
	Aliases: []string{"b", "breakpoint"},
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupBreakpoints,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("breakpoint added at <locspec>")
	},
}

func init() {
	debugRootCmd.AddCommand(breakCmd)
}
