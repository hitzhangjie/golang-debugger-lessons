package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	loop := true
	for loop {
		fmt.Println("pid:", os.Getpid())
		time.Sleep(time.Second)
	}
}
