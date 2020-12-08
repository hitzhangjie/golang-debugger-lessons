package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
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
	//fmt.Println("read dwarf ok")

	err = parseDebugInfo(dw)
	if err != nil {
		fmt.Println(err)
	}

	//for _, cu := range compileUnits {
	//	fmt.Printf("cu: %+v\n", cu)
	//}

	//printFuncMainMain(dw)
}

type Variable struct {
	Name string
}

type Function struct {
	Name      string
	DeclFile  string
	Variables []*Variable
}

type CompileUnit struct {
	Source []string
	Funcs  []*Function
}

var compileUnits = []*CompileUnit{}

func parseDebugInfo(dw *dwarf.Data) error {
	rd := dw.Reader()

	var curCompileUnit *CompileUnit
	var curFunction *Function

	for idx := 0; ; idx++ {
		entry, err := rd.Next()
		if err != nil {
			return fmt.Errorf("iterate entry error: %v", err)
		}
		if entry == nil {
			return nil
		}

		if entry.Tag == dwarf.TagCompileUnit {
			lrd, err := dw.LineReader(entry)
			if err != nil {
				return err
			}

			cu := &CompileUnit{}
			curCompileUnit = cu

			for _, v := range lrd.Files() {
				if v == nil {
					continue
				}
				cu.Source = append(cu.Source, v.Name)
			}
			compileUnits = append(compileUnits, cu)
		}

		if entry.Tag == dwarf.TagSubprogram {
			fn := &Function{
				Name:     entry.Val(dwarf.AttrName).(string),
				DeclFile: curCompileUnit.Source[entry.Val(dwarf.AttrDeclFile).(int64)-1],
			}
			curFunction = fn
			curCompileUnit.Funcs = append(curCompileUnit.Funcs, fn)

			if fn.Name == "main.main" {
				printEntry(entry)
				fmt.Printf("main.main is defined in %s\n", fn.DeclFile)
			}
		}

		if entry.Tag == dwarf.TagVariable {
			variable := &Variable{
				Name: entry.Val(dwarf.AttrName).(string),
			}
			curFunction.Variables = append(curFunction.Variables, variable)
			if curFunction.Name == "main.main" {
				printEntry(entry)
			}
		}
	}
	return nil
}

func printEntry(entry *dwarf.Entry) {
	fmt.Println("children:", entry.Children)
	fmt.Println("offset:", entry.Offset)
	fmt.Println("tag:", entry.Tag.String())
	for _, f := range entry.Field {
		fmt.Println("attr:", f.Attr, f.Val, f.Class)
	}
}
