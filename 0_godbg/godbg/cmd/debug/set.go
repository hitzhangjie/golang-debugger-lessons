package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <var|reg>=<value>",
	Short: "设置变量或寄存器值",
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(setCmd)
}
