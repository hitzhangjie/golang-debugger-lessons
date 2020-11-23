package debug

import (
	"errors"
	"fmt"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var pmemCmd = &cobra.Command{
	Use:   "pmem ",
	Short: "打印内存数据",
	Annotations: map[string]string{
		cmdGroupKey: cmdGroupInfo,
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		count, _ := cmd.Flags().GetUint("count")
		format, _ := cmd.Flags().GetString("fmt")
		size, _ := cmd.Flags().GetUint("size")
		addr, _ := cmd.Flags().GetString("addr")

		// check params
		err := checkPmemArgs(count, format, size, addr)
		if err != nil {
			return err
		}

		// calculate size of memory to read
		readAt, _ := strconv.ParseUint(addr, 0, 64)
		bytes := count * size

		buf := make([]byte, bytes, bytes)
		n, err := syscall.PtracePeekData(TraceePID, uintptr(readAt), buf)
		if err != nil || n != int(bytes) {
			return fmt.Errorf("read %d bytes, error: %v", n, err)
		}

		fmt.Printf("read %d bytes ok:", n)
		for _, b := range buf[:n] {
			fmt.Printf("%x", b)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	debugRootCmd.AddCommand(pmemCmd)
	// 类似gdb的命令x/FMT，其中FMT=重复数字+格式化修饰符+size
	pmemCmd.Flags().Uint("count", 16, "查看数值数量")
	pmemCmd.Flags().String("fmt", "hex", "数值打印格式: b(binary), o(octal), x(hex), d(decimal), ud(unsigned decimal)")
	pmemCmd.Flags().Uint("size", 4, "数值占用字节")
	pmemCmd.Flags().String("addr", "", "读取的内存地址")
}

func checkPmemArgs(count uint, format string, size uint, addr string) error {
	if count == 0 {
		return errors.New("invalid count")
	}
	if size == 0 {
		return errors.New("invalid size")
	}
	formats := map[string]struct{}{
		"b":  {},
		"o":  {},
		"x":  {},
		"d":  {},
		"ud": {},
	}
	if _, ok := formats[format]; !ok {
		return errors.New("invalid format")
	}
	// TODO make it compatible
	_, err := strconv.ParseUint(addr, 0, 64)
	return err
}
