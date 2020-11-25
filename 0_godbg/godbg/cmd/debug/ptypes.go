package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var ptypeCmd = &cobra.Command{
	Use:     "ptype <variable|type>",
	Short:   "打印变量类型信息",
	Aliases: []string{"pt"},
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(ptypeCmd)
}
