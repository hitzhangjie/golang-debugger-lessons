package elf_test

import (
	"debug/elf"
	"debug/gosym"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/tabwriter"

	"github.com/stretchr/testify/require"
	"golang.org/x/arch/x86/x86asm"
)

var testELFBinary = "../testdata/loop2"

func TestParseELFFile(t *testing.T) {
	// open elf
	file, err := elf.Open(testELFBinary)
	require.Nil(t, err)
	// prettyprint(file)

	//-----------------------------------------------------------------------
	// prog headers
	fmt.Println("dumping the elf.Progs")
	fmt.Println(strings.Repeat("-", 120))

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 3, ' ', 0)
	fmt.Fprintf(tw, "No.\tType\tFlags\tVAddr\tMemSize\n")
	for idx, p := range file.Progs {
		fmt.Fprintf(tw, "%d\t%v\t%v\t%#x\t%d\n", idx, p.Type, p.Flags, p.Vaddr, p.Memsz)
	}
	tw.Flush()
	println()

	//-----------------------------------------------------------------------
	// section headers

	fmt.Println("dumping the elf.Sections")
	fmt.Println(strings.Repeat("-", 120))

	tw = tabwriter.NewWriter(os.Stdout, 0, 4, 3, ' ', 0)
	heading := "No.\tName\tType\tFlags\tAddr\tOffset\tSize\tLink\tInfo\tAddrAlign\tEntSize\tFileSize\n"
	fmt.Fprintf(tw, heading)
	for idx, s := range file.Sections {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%#x\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			idx, s.Name, s.Type.String(), s.Flags.String(), s.Addr, s.Offset,
			s.Size, s.Link, s.Info, s.Addralign, s.Entsize, s.FileSize)
	}
	tw.Flush()
	println()

	//-----------------------------------------------------------------------
	// using sections

	fmt.Println("using .text section")
	fmt.Println(strings.Repeat("-", 120))

	// .text section
	dat, err := file.Section(".text").Data()
	if err != nil {
		panic(err)
	}
	fmt.Printf("% x\n", dat[:32])

	offset := 0
	for i := 0; i < 10; i++ {
		inst, err := x86asm.Decode(dat[offset:], 64)
		if err != nil {
			break
		}
		fmt.Println(x86asm.GNUSyntax(inst, 0, nil))
		offset += inst.Len
	}
	println()

	fmt.Println("using .data section")
	fmt.Println(strings.Repeat("-", 120))

	dataSection := file.Section(".data")
	dat, err = dataSection.Data()
	require.Nil(t, err)
	fmt.Printf("% x\n", dat[:32])

	// 查看gopclntab部分
	symtab, _ := file.Section(".gosymtab").Data()
	fmt.Println("symtab size:", len(symtab))

	pclntab, _ := file.Section(".gopclntab").Data()
	linetab := gosym.NewLineTable(pclntab, file.Entry)

	table, err := gosym.NewTable(symtab, linetab)
	require.Nil(t, err)

	fmt.Println("dump files:")
	i := 0
	for k, _ := range table.Files {
		fmt.Println("file:", k)
		i++
		if i > 10 {
			break
		}
	}
	fmt.Println("number of files:", len(table.Files))
	fmt.Println()

	fmt.Println("dump funcs:")
	for i, f := range table.Funcs {
		fmt.Printf("func: %s\n", f.Name)
		if i > 10 {
			break
		}
	}
	fmt.Println("number of funcs:", len(table.Funcs))
	fmt.Println()

	fmt.Println("dump syms:")
	for i, sym := range table.Syms {
		fmt.Printf("sym: %s\n", sym.Name)
		if i > 10 {
			break
		}
	}
	fmt.Println("number of syms:", len(table.Syms))
	fmt.Println()

	fmt.Println("main.main entry:", table.LookupFunc("main.main").Entry)
	fp, lineno := func() (string, int) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return filepath.Join(filepath.Dir(wd), "testdata/loop2.go"), 15
	}()
	pc, fn, err := table.LineToPC(fp, lineno)
	require.Nil(t, err)
	fmt.Println(pc, fn.Name, err)

	f, l, fn := table.PCToLine(pc)
	fmt.Println(f, l, fn.Name)

	// stacktrace，go运行时加载了.gosymtab, .gopclntab之后，直接指定pc就能拿到stacktrace吗？
	// 不能，还是需要目标程序必须运行起来才可以，这个stack肯定是回溯栈帧记录才可以打印出这个堆栈的。
	//
	// 只有可能回溯这个调用栈记录，才能一次次回溯找到caller的pc地址，才能根据这个pclntab找到源码
	// 层面的调用栈信息 ... OK, 到此结束。
	//
	// runtime.Stack()
}

func prettyprint(v any) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
