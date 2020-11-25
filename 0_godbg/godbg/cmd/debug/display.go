package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var displayCmd = &cobra.Command{
	Use:   "display <var|reg>",
	Short: "始终显示变量或寄存器值",
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(displayCmd)
}
