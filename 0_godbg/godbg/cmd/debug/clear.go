package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear <n>",
	Short: "清除指定编号的断点",
	Long:  `清除指定编号的断点`,
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupBreakpoints,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clear ${n}th breakpoint")
	},
}

func init() {
	debugRootCmd.AddCommand(clearCmd)

	clearCmd.Flags().Uint32P("no", "n", 1, "断点编号")
}
