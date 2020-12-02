package main

import (
	"debug/elf"
	"debug/gosym"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/arch/x86/x86asm"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go <prog>")
		os.Exit(1)
	}
	prog := os.Args[1]

	file, err := elf.Open(prog)
	if err != nil {
		panic(err)
	}
	//spew.Dump(file)

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 3, ' ', 0)
	fmt.Fprintf(tw, "No.\tType\tFlags\tVAddr\tMemSize\n")
	for idx, p := range file.Progs {
		fmt.Fprintf(tw, "%d\t%v\t%v\t%#x\t%d\n", idx, p.Type, p.Flags, p.Vaddr, p.Memsz)
	}
	tw.Flush()
	println()

	text := file.Progs[2]
	buf := make([]byte, text.Filesz, text.Filesz)
	n, err := text.ReadAt(buf, int64(64+56))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", buf[:64])
	fmt.Printf("i have read some data: %d bytes\n", n)

	i := 0
	count := 0
	for {
		if count > 10 {
			break
		}
		inst, err := x86asm.Decode(buf[i:], 64)
		if err != nil {
			panic(err)
		}
		asm := x86asm.GoSyntax(inst, uint64(inst.PCRel), nil)
		fmt.Printf("%#x %-16x %s\n", i, buf[i:i+inst.Len], asm)
		i += inst.Len
		count++
	}
	fmt.Println()
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

func printStackTrace(pclntab *gosym.Table, pc uint64, depth int) {
	for i := 0; i < depth; i++ {
		file, ln, fn := pclntab.PCToLine(pc)
		fmt.Printf("func: %s, pos:%s:%d\n", fn.Name, file, ln)
		pc = fn.Entry - 1
	}
}
