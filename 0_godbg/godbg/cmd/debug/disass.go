package debug

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/arch/x86/x86asm"
)

var disassCmd = &cobra.Command{
	Use:   "disass <locspec>",
	Short: "反汇编机器指令",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupSource,
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		max, _ := cmd.Flags().GetUint("n")
		syntax, _ := cmd.Flags().GetString("syntax")

		// 读取PC值
		regs := syscall.PtraceRegs{}
		err := syscall.PtraceGetRegs(TraceePID, &regs)
		if err != nil {
			return err
		}

		buf := make([]byte, 1)
		n, err := syscall.PtracePeekText(TraceePID, uintptr(regs.PC()), buf)
		if err != nil || n != 1 {
			return fmt.Errorf("peek text error: %v, bytes: %d", err, n)
		}
		//fmt.Printf("read %d bytes, value of %x\n", n, buf[0])

		// read a breakpoint
		if buf[0] == 0xCC {
			regs.SetPC(regs.PC() - 1)
		}

		// 查找，如果之前设置过断点，将恢复
		dat := make([]byte, 1024)
		n, err = syscall.PtracePeekText(TraceePID, uintptr(regs.PC()), dat)
		if err != nil {
			return fmt.Errorf("peek text error: %v, bytes: %d", err, n)
		}
		//fmt.Printf("size of text: %d\n", n)

		// 反汇编这里的指令数据
		offset := 0
		count := 0
		for uint(count) < max {
			inst, err := x86asm.Decode(dat[offset:], 64)
			if err != nil {
				return fmt.Errorf("x86asm decode error: %v", err)
			}
			//fmt.Printf("%#x %s\n", regs.PC()+uint64(offset), inst.String())

			asm, err := instSyntax(inst, syntax)
			if err != nil {
				return fmt.Errorf("x86asm syntax error: %v", err)
			}
			fmt.Printf("%#x %s\n", regs.PC()+uint64(offset), asm)
			offset += inst.Len
			count++
		}
		return nil
	},
}

func init() {
	debugRootCmd.AddCommand(disassCmd)

	disassCmd.Flags().UintP("n", "n", 10, "反汇编指令数量")
	disassCmd.Flags().StringP("syntax", "s", "gnu", "反汇编指令语法，支持：go, gnu, intel")
}

// GetExecutable 根据pid获取可执行程序路径
func GetExecutable(pid int) (string, error) {
	exeLink := fmt.Sprintf("/proc/%d/exe", pid)
	exePath, err := os.Readlink(exeLink)
	if err != nil {
		return "", err
	}
	return exePath, nil
}

func instSyntax(inst x86asm.Inst, syntax string) (string, error) {
	asm := ""
	switch syntax {
	case "go":
		asm = x86asm.GoSyntax(inst, uint64(inst.PCRel), nil)
	case "gnu":
		asm = x86asm.GNUSyntax(inst, uint64(inst.PCRel), nil)
	case "intel":
		asm = x86asm.IntelSyntax(inst, uint64(inst.PCRel), nil)
	default:
		return "", fmt.Errorf("invalid asm syntax error")
	}
	return asm, nil
}
