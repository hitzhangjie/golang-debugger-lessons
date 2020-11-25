package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var frameCmd = &cobra.Command{
	Use:   "frame <frame no.>",
	Short: "选择调用栈中栈帧",
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(frameCmd)
}
