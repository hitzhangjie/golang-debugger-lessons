package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "结束调试会话",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupOthers,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(exitCmd)
}
