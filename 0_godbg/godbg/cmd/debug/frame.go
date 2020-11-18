package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var frameCmd = &cobra.Command{
	Use:   "frame",
	Short: "选择调用栈中栈帧",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(frameCmd)
}
