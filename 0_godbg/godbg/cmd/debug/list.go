package debug

import (
	"fmt"

	"github.com/spf13/cobra"
	cobraprompt "github.com/stromland/cobra-prompt"
)

var listCmd = &cobra.Command{
	Use:     "list <linespec>",
	Short:   "查看源码信息",
	Aliases: []string{"l"},
	Annotations: map[string]string{
		cmdGroupKey:                     cmdGroupSource,
		cobraprompt.CALLBACK_ANNOTATION: suggestionListSourceFiles,
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list codes in file")
	},
}

func init() {
	debugRootCmd.AddCommand(listCmd)
}
