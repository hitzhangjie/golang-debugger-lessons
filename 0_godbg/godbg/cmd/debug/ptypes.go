package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var ptypesCmd = &cobra.Command{
	Use:   "ptypes <variable>",
	Short: "打印变量类型信息",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(ptypesCmd)
}
