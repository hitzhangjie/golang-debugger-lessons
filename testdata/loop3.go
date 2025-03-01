package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.LockOSThread()

	for loop() {
		fmt.Println("pid:", os.Getpid())
		time.Sleep(time.Second)
	}
}

//go:noinline
func loop() bool {
	return true
}
