package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var disassCmd = &cobra.Command{
	Use:   "disass <locspec>",
	Short: "反汇编机器指令",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupSource,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("disass <locspec>")
	},
}

func init() {
	debugRootCmd.AddCommand(disassCmd)
}
