package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stepCmd = &cobra.Command{
	Use:   "step",
	Short: "执行一条指令",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupCtrlFlow,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("step")
	},
}

func init() {
	debugRootCmd.AddCommand(stepCmd)
}
