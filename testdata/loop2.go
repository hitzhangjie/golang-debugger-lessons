package main

import "fmt"
import "os"
import "time"

func init() {
	go func() {
		for {
			fmt.Println("main.func1 pid:", os.Getpid())
			time.Sleep(time.Second)
		}
	}()
}
func main() {
	pid := os.Getpid()
	for {
		fmt.Println("main.main pid:", pid)
		time.Sleep(time.Second * 3)
	}
}
