package main

import (
	"debug/elf"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/arch/x86/x86asm"
)

func main() {

	if len(os.Args) != 2 {
		panic("invalid params")
	}

	// pid
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	// 通过pid找到可执行程序路径
	exePath, err := GetExecutable(pid)
	if err != nil {
		panic(err)
	}
	fmt.Println(exePath)

	// 读取指令信息并反汇编
	elfFile, err := elf.Open(exePath)
	if err != nil {
		panic(err)
	}
	section := elfFile.Section(".text")
	buf, err := section.Data()
	if err != nil {
		panic(err)
	}

	offset := 0
	for {
		inst, err := x86asm.Decode(buf[offset:], 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%8x %s\n", offset, inst.String())
		offset += inst.Len
	}
}

// GetExecutable 根据pid获取可执行程序路径
func GetExecutable(pid int) (string, error) {
	exeLink := fmt.Sprintf("/proc/%d/exe", pid)
	exePath, err := os.Readlink(exeLink)
	if err != nil {
		return "", fmt.Errorf("find executable by pid err: %w", err)
	}
	return exePath, nil
}
