package main

import (
	"debug/elf"
	"debug/gosym"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/arch/x86/x86asm"
)

func main() {
	wd, _ := os.Getwd()
	file, err := elf.Open(filepath.Join(wd, "testdata/loop2"))
	if err != nil {
		panic(err)
	}
	spew.Printf("elf open ok, file:\n")
	//spew.Dump(file)

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

}
