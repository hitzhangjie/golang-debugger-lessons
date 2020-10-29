package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"io"
	"os"
)

func main() {
	// TODO 指定pid，获取path，如linux下通过/proc/<pid>/exe来获取path

	// elf file open
	target := "/root/debugger101/4_dwarf/main"
	fin, err := os.Open(target)
	if err != nil {
		panic(err)
	}

	file, err := elf.NewFile(fin)
	if err != nil {
		panic(err)
	}

	// read .debug_info section
	s := file.Section(".debug_info")
	if s == nil {
		s = file.Section(".zdebug_info")
	}

	dat, err := s.Data()
	if err != nil {
		panic(err)
	}
	_ = dat

	dwarfData, err := file.DWARF()
	if err != nil {
		panic(err)
	}

	reader := dwarfData.Reader()

	for {
		// each compilation unit
		entry, err := reader.Next()
		if err != nil {
			panic(err)
		}
		if entry == nil {
			break
		}

		// compilation line table
		switch entry.Tag {
		case dwarf.TagCompileUnit:
			lineReader, err := dwarfData.LineReader(entry)
			if err != nil {
				panic(err)
			}
			lineEntry := dwarf.LineEntry{}
			for {
				err = lineReader.Next(&lineEntry)
				if err == io.EOF {
					break
				}
				cu := entry.AttrField(dwarf.AttrName).Val.(string)
				if len(cu) == 0 {
					continue
				}
				if cu != "main" {
					continue
				}
				fmt.Fprintf(os.Stdout, "CompileUnit: %s, address: %v, position: %s:%d:%d\n",
					cu, lineEntry.Address, lineEntry.File.Name, lineEntry.Line, lineEntry.Column)
			}
		default:
		}

	}

}
