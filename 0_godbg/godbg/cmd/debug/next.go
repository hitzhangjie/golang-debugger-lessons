package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "执行一条语句",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupCtrlFlow,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("next")
	},
}

func init() {
	debugRootCmd.AddCommand(nextCmd)
}
