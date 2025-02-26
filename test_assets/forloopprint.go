package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	for {
		fmt.Println("pid:", os.Getpid())
		time.Sleep(time.Second)
	}
}
