package debug

import (
	"fmt"
	"syscall"

	"godbg/target"

	"github.com/spf13/cobra"
)

var clearallCmd = &cobra.Command{
	Use:   "clearall <n>",
	Short: "清除所有的断点",
	Long:  `清除所有的断点`,
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupBreakpoints,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("clearall")

		for _, brk := range breakpoints {
			n, err := syscall.PtracePokeData(TraceePID, brk.Addr, []byte{brk.Orig})
			if err != nil || n != 1 {
				return fmt.Errorf("清空断点失败: %v", err)
			}
		}

		breakpoints = map[uintptr]*target.Breakpoint{}
		fmt.Println("清空断点成功")
		return nil
	},
}

func init() {
	debugRootCmd.AddCommand(clearallCmd)
}
