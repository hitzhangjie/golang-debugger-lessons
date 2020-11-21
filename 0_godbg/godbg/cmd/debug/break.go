package debug

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"godbg/target"

	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break <locspec>",
	Short: "在源码中添加断点",
	Long: `在源码中添加断点，源码位置可以通过locspec格式指定。

当前支持的locspec格式，包括两种:
- 指令地址
- [文件名:]行号
- [文件名:]函数名`,
	Aliases: []string{"b", "breakpoint"},
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupBreakpoints,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("break %s\n", strings.Join(args, " "))

		if len(args) != 1 {
			return errors.New("参数错误")
		}

		locStr := args[0]
		v, err := strconv.ParseUint(locStr, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid locspec: %v", err)
		}
		addr := uintptr(v)

		orig := [1]byte{}
		n, err := syscall.PtracePeekText(TraceePID, addr, orig[:])
		if err != nil || n != 1 {
			return fmt.Errorf("peek text, %d bytes, error: %v", n, err)
		}
		breakpoint, err := target.NewBreakpoint(addr, orig[0], "")
		if err != nil {
			return fmt.Errorf("add breakpoint error: %v", err)
		}
		breakpoints[addr] = &breakpoint

		n, err = syscall.PtracePokeText(TraceePID, addr, []byte{0xCC})
		if err != nil || n != 1 {
			return fmt.Errorf("poke text, %d bytes, error: %v", n, err)
		}
		fmt.Printf("添加断点成功\n")
		return nil
	},
}

func init() {
	debugRootCmd.AddCommand(breakCmd)
}
