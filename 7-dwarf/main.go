package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go <prog>")
		os.Exit(1)
	}
	prog := os.Args[1]

	// open elf
	file, err := elf.Open(prog)
	if err != nil {
		panic(err)
	}
	//spew.Dump(file)

	// dwarf调试信息遍历
	dw, err := file.DWARF()
	if err != nil {
		panic(err)
	}

	rd := dw.Reader()

	// next compilation unit
	for {
		entry, err := rd.Next()
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		if entry == nil {
			break
		}

		if entry.Tag == dwarf.TagCompileUnit {
			fmt.Println("CompilationUnit:", entry.Field)
		}

		if entry.Tag == dwarf.TagSubprogram {
			// 读取.debug_line关联的行号表信息
			lrd, err := dw.LineReader(entry)
			if err != nil {
				fmt.Println(err)
				break
			}
			if lrd == nil {
				continue
			}
			for {
				lentry := dwarf.LineEntry{}
				err = lrd.Next(&lentry)
				if err == io.EOF {
					break
				}
				fmt.Printf("lineTable: %s:%d\n", lentry.File.Name, lentry.Line)
			}
		}

		fmt.Println(entry)
	}
}
