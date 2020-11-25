package debug

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"unsafe"

	"github.com/spf13/cobra"
)

var pmemCmd = &cobra.Command{
	Use:   "pmem [flags] <address>",
	Short: "打印内存数据",
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupInfo,
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
		buf = buf[:n]

		fmt.Printf("read %d bytes ok:\n", n)
		s := prettyPrintMem(uintptr(readAt), buf, isLittleEndian(), format[0], int(size))
		fmt.Println(s)

		return nil
	},
}

func init() {
	debugRootCmd.AddCommand(pmemCmd)

	// 类似gdb的命令x/FMT，其中FMT=重复数字+格式化修饰符+size
	pmemCmd.Flags().Uint("count", 16, "查看数值数量")
	pmemCmd.Flags().String("fmt", "hex", "数值打印格式: b(binary), o(octal), x(hex), d(decimal)") // TODO signed/unsigned
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

// prettyPrintMem examine the memory and format data
//
// `format` specifies the data format (or data type), `size` specifies size of each data,
// like 4byte integer, 1byte character, etc. `count` specifies the number of values.
func prettyPrintMem(address uintptr, memArea []byte, littleEndian bool, format byte, size int) string {

	var (
		cols      int
		colFormat string
		colBytes  = size

		addrLen int
		addrFmt string
	)

	// Diffrent versions of golang output differently about '#'.
	// See https://ci.appveyor.com/project/derekparker/delve-facy3/builds/30179356.
	switch format {
	case 'b':
		cols = 4 // Avoid emitting rows that are too long when using binary format
		colFormat = fmt.Sprintf("%%0%db", colBytes*8)
	case 'o':
		cols = 8
		colFormat = fmt.Sprintf("0%%0%do", colBytes*3) // Always keep one leading zero for octal.
	case 'd':
		cols = 8
		colFormat = fmt.Sprintf("%%0%dd", colBytes*3)
	case 'x':
		cols = 8
		colFormat = fmt.Sprintf("0x%%0%dx", colBytes*2) // Always keep one leading '0x' for hex.
	default:
		return fmt.Sprintf("not supprted format %q\n", string(format))
	}
	colFormat += "\t"

	l := len(memArea)
	rows := l / (cols * colBytes)
	if l%(cols*colBytes) != 0 {
		rows++
	}

	// Avoid the lens of two adjacent address are different, so always use the last addr's len to format.
	if l != 0 {
		addrLen = len(fmt.Sprintf("%x", uint64(address)+uint64(l)))
	}
	addrFmt = "0x%0" + strconv.Itoa(addrLen) + "x:\t"

	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 3, ' ', 0)

	for i := 0; i < rows; i++ {
		fmt.Fprintf(w, addrFmt, address)

		for j := 0; j < cols; j++ {
			offset := i*(cols*colBytes) + j*colBytes
			if offset+colBytes <= len(memArea) {
				n := byteArrayToUInt64(memArea[offset:offset+colBytes], littleEndian)
				fmt.Fprintf(w, colFormat, n)
			}
		}
		fmt.Fprintln(w, "")
		address += uintptr(cols)
	}
	w.Flush()
	return b.String()
}

// 将byte slice转成uint64数值，注意字节序影响
func byteArrayToUInt64(buf []byte, isLittleEndian bool) uint64 {
	var n uint64
	if isLittleEndian {
		for i := len(buf) - 1; i >= 0; i-- {
			n = n<<8 + uint64(buf[i])
		}
	} else {
		for i := 0; i < len(buf); i++ {
			n = n<<8 + uint64(buf[i])
		}
	}
	return n
}

// 判断是否是小端字节序
func isLittleEndian() bool {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		return true
	case [2]byte{0xAB, 0xCD}:
		return false
	default:
		panic("Could not determine native endianness.")
	}
}
