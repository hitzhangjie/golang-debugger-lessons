package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
		}
	}()
	go func() {
		http.ListenAndServe(":8888", nil)
	}()
	for {
		panic("fuck")
		fmt.Println("vim-go")
		time.Sleep(time.Second)
	}
}
