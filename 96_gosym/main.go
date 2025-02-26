package main

import (
	"debug/elf"
	"debug/gosym"
	"fmt"
	"os"
	"strings"
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

	//---------------------------------------------------------------------
	// 调用栈信息

	fmt.Println("using lineno table: stacktrace")
	fmt.Println(strings.Repeat("-", 120))
}

