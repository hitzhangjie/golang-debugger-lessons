package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print <var|reg>",
	Short: "打印变量或寄存器值",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(printCmd)
}
