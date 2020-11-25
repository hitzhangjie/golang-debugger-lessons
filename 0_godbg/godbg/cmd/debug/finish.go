package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "退出当前函数",
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupCtrlFlow,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("finish")
	},
}

func init() {
	debugRootCmd.AddCommand(finishCmd)
}
