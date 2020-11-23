package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var pmemCmd = &cobra.Command{
	Use:   "pmem ",
	Short: "打印内存数据",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(pmemCmd)
}
