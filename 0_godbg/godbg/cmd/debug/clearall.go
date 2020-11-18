package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearallCmd = &cobra.Command{
	Use:   "clearall <n>",
	Short: "清除所有的断点",
	Long:  `清除所有的断点`,
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupBreakpoints,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clearall breakpoints")
	},
}

func init() {
	debugRootCmd.AddCommand(clearallCmd)
}
