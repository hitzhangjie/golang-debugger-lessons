package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	usage = "Usage: go run main.go exec <path/to/prog>"

	cmdExec = "exec"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s\n\n", usage)
		os.Exit(1)
	}
	cmd := os.Args[1]

	switch cmd {
	case cmdExec:
		prog := os.Args[2]

		fin, err := os.Lstat(prog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lstat error: %v\n\n", err)
			os.Exit(1)
		}

		if fin.IsDir() {
			fmt.Fprintf(os.Stderr, "%s is directory\n\n", prog)
			os.Exit(1)
		}

		if fin.Mode()&0111 == 0 {
			fmt.Fprintf(os.Stderr, "%s not executable\n\n", prog)
			os.Exit(1)
		}

		progCmd := exec.Command(prog)
		buf, err := progCmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s exec error: %v, \n\n%s\n\n", err, string(buf))
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(buf))
	default:
		fmt.Fprintf(os.Stderr, "%s unknown cmd\n\n", cmd)
		os.Exit(1)
	}

}
