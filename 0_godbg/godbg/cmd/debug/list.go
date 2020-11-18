package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list <linespec>",
	Short:   "查看源码信息",
	Aliases: []string{"l"},
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupSource,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list codes in file")
	},
}

func init() {
	debugRootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("list", "l", "", "显示位置附近源码")
}
