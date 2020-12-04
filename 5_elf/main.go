package main

import (
	"debug/elf"
	"debug/gosym"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	//-----------------------------------------------------------------------
	// elf file

	fmt.Println("processing elf binary:", prog)
	println()

	// open elf
	file, err := elf.Open(prog)
	if err != nil {
		panic(err)
	}
	//spew.Dump(file)

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
	if err != nil {
		panic(err)
	}
	_ = dat
	fmt.Printf("% x\n", dat[:32])
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
}
