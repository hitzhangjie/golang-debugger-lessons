package main

import (
	"debug/elf"
	"debug/gosym"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/arch/x86/x86asm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go <prog>")
		os.Exit(1)
	}
	prog := os.Args[1]

	//--------------------------------------------------------------------
	// open elf file

	fmt.Println("processing elf file: ", prog)

	file, err := elf.Open(prog)
	if err != nil {
		panic(err)
	}
	//spew.Dump(file)
	println()

	//---------------------------------------------------------------------
	// 构建行号表

	fmt.Println("constructing lineno table")

	gosymtab, _ := file.Section(".gosymtab").Data()
	gopclntab, _ := file.Section(".gopclntab").Data()

	pclntab := gosym.NewLineTable(gopclntab, file.Section(".text").Addr)
	table, _ := gosym.NewTable(gosymtab, pclntab)
	println()

	// 利用行号表
	fmt.Println("using lineno table: file:lineno to vaddr")
	fmt.Println(strings.Repeat("-", 120))

	type arg struct {
		file   string
		lineno int
		pc     uintptr
	}

	args := []arg{
		{file: "/root/debugger101/testdata/loop2.go", lineno: 3},
		{file: "/root/debugger101/testdata/loop2.go", lineno: 9},
		{file: "/root/debugger101/testdata/loop2.go", lineno: 11},
		{file: "/root/debugger101/testdata/loop2.go", lineno: 17},
	}

	for _, arg := range args {
		pc, fn, err := table.LineToPC(arg.file, arg.lineno)
		if err != nil {
			fmt.Printf("file:lineno to vaddr, file:lineno = %s:%d, error: %v\n", arg.file, arg.lineno, err)
			continue
		}
		fmt.Printf("file:lineno to vaddr, file:lineno = %s:%d, pc = %#x, fn = %s\n", arg.file, arg.lineno, pc, fn.Name)
	}
	println()

	fmt.Println("using lineno table: vaddr to file:lineno")
	fmt.Println(strings.Repeat("-", 120))

	f, l, fn := table.PCToLine(0x4b86cf)
	fmt.Printf("vaddr to file:lineno, pc = %#x, fn = %s, file:lineno = %s:%d\n", 0x4b86cf, fn.Name, f, l)
	println()

	fmt.Println("using pclntab: print stacktrace")
	fmt.Println(strings.Repeat("-", 120))

	printStackTrace(table, 0x4b86cf, 20)
}

func main2() {
	wd, _ := os.Getwd()
	file, err := elf.Open(filepath.Join(wd, "testdata/loop2"))
	if err != nil {
		panic(err)
	}
	spew.Printf("elf open ok, file:\n")
	spew.Dump(file)

	os.Exit(1)
	// 查看section列表
	for idx, sec := range file.Sections {
		buf, _ := sec.Data()
		fmt.Printf("[%d] %-16s %-16s size:%d\n", idx, sec.Name, sec.Type.String(), len(buf))
	}
	fmt.Println()

	// 查看指令部分
	text, _ := file.Section(".text").Data()
	i := 0
	count := 0
	for {
		if count > 10 {
			break
		}
		inst, err := x86asm.Decode(text[i:], 64)
		if err != nil {
			break
		}
		fmt.Printf("%#x %-16x %s\n", i, text[i:i+inst.Len], inst.String())
		i += inst.Len
		count++
	}
	fmt.Println()

	// 查看符号表部分
	symtab, _ := file.Section(".gosymtab").Data()
	fmt.Println("symtab size:", len(symtab))

	pclntab, _ := file.Section(".gopclntab").Data()
	linetab := gosym.NewLineTable(pclntab, file.Entry)

	table, err := gosym.NewTable(symtab, linetab)
	if err != nil {
		panic(err)
	}
	for k, _ := range table.Files {
		fmt.Println("file:", k)
	}

	for _, f := range table.Funcs {
		fmt.Printf("func: %s\n", f.Name)
	}
	//fmt.Println("objs size:", len(table.Objs))

	//for _, sym := range table.Syms {
	//	fmt.Println("sym: %s\n", sym.Name)
	//}
	fmt.Println()

	fmt.Println("main.main entry:", table.LookupFunc("main.main").Entry)
	pc, fn, err := table.LineToPC("/root/main.go", 15)
	fmt.Println(pc, fn.Name, err)

	f, l, fn := table.PCToLine(pc)
	fmt.Println(f, l, fn.Name)

	// print call stack
	printStackTrace(table, fn.Entry, 3)
}

var lastFn string

func printStackTrace(pclntab *gosym.Table, pc uint64, depth int) {
	for i := 0; i < depth; i++ {
		for {
			file, ln, fn := pclntab.PCToLine(pc)
			_ = file
			_ = ln
			//fmt.Printf("func: %s, pos:%s:%d\n", fn.Name, file, ln)
			if fn.Name != lastFn {
				lastFn = fn.Name
				fmt.Printf("%s pc:%#x\n", fn.Name, pc)
				pc = fn.Entry - 1
				break
			}
			pc--
		}
	}
}
