package main

import (
	"bytes"
	"compress/zlib"
	"debug/dwarf"
	"debug/elf"
	"encoding/binary"
	"errors"
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

	err = parseDwarf(dw)
	if err != nil {
		panic(err)
	}

	pc, err := find("/root/debugger101/testdata/loop2.go", 16)
	if err != nil {
		panic(err)
	}
	fmt.Printf("found pc: %#x\n", pc)

	// .debug_frame & .zdebug_frame
	dat, err :=readDebugFrame(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read .debug/zdebug_frame ok, size: %d\n", len(dat))
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
	Lines  []*dwarf.LineEntry
}

var compileUnits = []*CompileUnit{}

func parseDwarf(dw *dwarf.Data) error {
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

			for {
				var e dwarf.LineEntry
				err := lrd.Next(&e)
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
				curCompileUnit.Lines = append(curCompileUnit.Lines, &e)
			}
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

func find(file string, lineno int) (pc uint64, err error) {
	for _, cu := range compileUnits {
		for _, e := range cu.Lines {
			if e.File.Name != file {
				continue
			}
			if e.Line != lineno {
				continue
			}
			if !e.IsStmt {
				continue
			}
			return e.Address, nil
		}
	}
	return 0, errors.New("not found")
}

func readDebugFrame(file *elf.File) ([]byte, error) {

	s := file.Section(".debug_frame")
	if s != nil {
		return s.Data()
	}

	s = file.Section(".zdebug_frame")
	buf, err := s.Data()
	if err != nil {
		return nil, err
	}

	return decompressMaybe(buf)
}


func decompressMaybe(b []byte) ([]byte, error) {
	if len(b) < 12 || string(b[:4]) != "ZLIB" {
		// not compressed
		return b, nil
	}

	dlen := binary.BigEndian.Uint64(b[4:12])
	dbuf := make([]byte, dlen)
	r, err := zlib.NewReader(bytes.NewBuffer(b[12:]))
	if err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(r, dbuf); err != nil {
		return nil, err
	}
	if err := r.Close(); err != nil {
		return nil, err
	}
	return dbuf, nil
}