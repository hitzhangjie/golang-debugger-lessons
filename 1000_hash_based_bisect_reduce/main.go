/*
	how to run this test?


	```bash
	$ bisect FEAT1=PATTERN FEAT2=PATTERN FEAT3=PATTERN ./main

	bisect FEAT1=PATTERN FEAT2=PATTERN FEAT3=PATTERN ./main
	bisect: checking target with all changes disabled
	bisect: run: FEAT1=n FEAT2=n FEAT3=n ./main... ok (0 matches)
	bisect: run: FEAT1=n FEAT2=n FEAT3=n ./main... ok (0 matches)
	bisect: checking target with all changes enabled
	bisect: run: FEAT1=y FEAT2=y FEAT3=y ./main... FAIL (1 matches)
	bisect: run: FEAT1=y FEAT2=y FEAT3=y ./main... FAIL (1 matches)
	bisect: target succeeds with no changes, fails with all changes
	bisect: searching for minimal set of enabled changes causing failure
	bisect: confirming failing change set
	bisect: run: FEAT1=v+x9 FEAT2=v+x9 FEAT3=v+x9 ./main... ok (1 matches)
	bisect: run: FEAT1=v+x9 FEAT2=v+x9 FEAT3=v+x9 ./main... ok (1 matches)
	bisect: confirmation run succeeded unexpectedly
	bisect: FOUND failing change set
	--- change set #1 (enabling changes causes failure)
	/home/zhangjie/hitzhangjie/codemaster/bisectv2/main.go:96 feat2() called
	---
	bisect: checking for more failures
	bisect: run: FEAT1=-x9 FEAT2=-x9 FEAT3=-x9 ./main... ok (0 matches)
	bisect: run: FEAT1=-x9 FEAT2=-x9 FEAT3=-x9 ./main... ok (0 matches)
	bisect: target succeeds with all remaining changes enabled
	```
*/

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	changelist1Matcher *Matcher // 1 and 2 enabled, will causes fail
	changelist2Matcher *Matcher // 1 and 2 enabled, will causes fail
	changelist3Matcher *Matcher // always success

	err error
)

const changelist1 = "FEAT1"
const changelist2 = "FEAT2"
const changelist3 = "FEAT3"

var (
	h1 = Hash(changelist1)
	h2 = Hash(changelist2)
	h3 = Hash(changelist3)
)

func init() {
	changelist1Matcher, err = New(os.Getenv(changelist1))
	if err != nil {
		log.Fatalf("failed to create matcher: %v", err)
	}
	changelist2Matcher, err = New(os.Getenv(changelist2))
	if err != nil {
		log.Fatalf("failed to create matcher: %v", err)
	}
	changelist3Matcher, err = New(os.Getenv(changelist3))
	if err != nil {
		log.Fatalf("failed to create matcher: %v", err)
	}
}

func main() {
	ch := make(chan int)
	go func() {
		// lock, missing unlock
		if changelist2Matcher.ShouldEnable(h2) {
			if changelist2Matcher.ShouldReport(h2) {
				fmt.Printf("%s %s feat2() called\n", flielineno(+2), Marker(h2))
			}
			feat2()
		} else {
			fmt.Println("disable feat2()")
		}

		// lock, and unlock
		if changelist1Matcher.ShouldEnable(h1) {
			if changelist2Matcher.ShouldReport(h2) {
				fmt.Printf("%s %s feat2() called\n", flielineno(+2), Marker(h2))
			}
			feat1()
		} else {
			fmt.Println("disable feat1()")
		}

		// no lock operation
		if changelist3Matcher.ShouldEnable(h3) {
			if changelist2Matcher.ShouldReport(h2) {
				fmt.Printf("%s %s feat2() called\n", flielineno(+2), Marker(h2))
			}
			feat3()
		} else {
			fmt.Println("disable feat3()")
		}
		ch <- 1
	}()

	select {
	case <-ch:
		os.Exit(0)
	case <-time.After(time.Second):
		os.Exit(1)
	}
}

var (
	mu  sync.Mutex
	val int64
)

// feat2 doesn't release, it will cause following feat1() feat2() call deadlock
func feat1() error {
	mu.Lock()
	val++
	mu.Unlock()
	return nil
}

// feat2 doesn't release, it will cause following feat1() feat2() call deadlock
func feat2() error {
	mu.Lock()
	val--
	//mu.Unlock()
	return nil
}

func feat3() error {
	return nil
}

func flielineno(offset int) string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", file, line+offset)
}
