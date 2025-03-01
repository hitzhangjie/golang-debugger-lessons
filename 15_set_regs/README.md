测试方法：

1、首先我们准备一个测试程序，loop3.go，该程序每隔1s输出一下pid，循环由固定返回true的loop()函数控制

```go
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.LockOSThread() // 减少调试多线程程序复杂性，否则会偶现ptrace操作时程序未stopped错误

	for loop() {
		fmt.Println("pid:", os.Getpid())
		time.Sleep(time.Second)
	}
}

//go:noinline
func loop() bool {
	return true
}
```

2、按照ABI调用惯例，这里的函数调用loop()的返回值会通过RAX寄存器返回，所以我们想在loop()函数调用返回后，通过修改RAX寄存器的值来篡改返回值为false。

那我们先确定下loop()函数的返回地址，这个只要我们通过dlv调试器在loop3.go:13添加断点，然后disass，就可以确定返回地址为 0x4af15e。

确定完返回地址后我们即可detach tracee，恢复其执行。

```bash
(dlv) disass
Sending output to pager...
TEXT main.main(SB) /home/zhangjie/debugger101/golang-debugger-lessons/testdata/loop3.go
        loop3.go:10     0x4af140        493b6610                cmp rsp, qword ptr [r14+0x10]
        loop3.go:10     0x4af144        0f8601010000            jbe 0x4af24b
        loop3.go:10     0x4af14a        55                      push rbp
        loop3.go:10     0x4af14b        4889e5                  mov rbp, rsp
        loop3.go:10     0x4af14e        4883ec70                sub rsp, 0x70
        loop3.go:11     0x4af152        e8e95ef9ff              call $runtime.LockOSThread
        loop3.go:13     0x4af157        eb00                    jmp 0x4af159
=>      loop3.go:13     0x4af159*       e802010000              call $main.loop
        loop3.go:13     0x4af15e        8844241f                mov byte ptr [rsp+0x1f], al
        ...
(dlv) quit
Would you like to kill the process? [Y/n] n
```

3、如果我们不加干扰，loop3会每隔1s不停地输出pid信息。

```bash
$ ./loop3
pid: 4946
pid: 4946
pid: 4946
pid: 4946
pid: 4946
...
zhangjie🦀 testdata(master) $
```

4、现在运行我们编写的调试工具 ./15_set_regs 4946,

```bash
$ ./15_set_regs 4946
===step1===: supposing running `dlv attach pid` here
process 4946 attach succ
process 4946 stopped
tracee stopped at 476263

===step2===: supposing running `dlv> b <addr>`  and `dlv> continue` here
enter return address of loop()
0x4af15e

you entered 4af15e
process 4946 stopped

===step3===: supposing change register RAX value from true to false
before RAX=1
after RAX=0                   <= 我们篡改了返回值为0
```

```bash
...
pid: 4946
pid: 4946
pid: 4946                      <= 因为篡改了loop()的返回值为false，循环跳出，程序结束
zhangjie🦀 testdata(master) $
```

```bash
(dlv) disass
TEXT main.loop(SB) /home/zhangjie/debugger101/golang-debugger-lessons/testdata/loop3.go
        loop3.go:20     0x4af260        55              push rbp
        loop3.go:20     0x4af261        4889e5          mov rbp, rsp
=>      loop3.go:20     0x4af264*       4883ec08        sub rsp, 0x8
        loop3.go:20     0x4af268        c644240700      mov byte ptr [rsp+0x7], 0x0
        loop3.go:21     0x4af26d        c644240701      mov byte ptr [rsp+0x7], 0x1
        loop3.go:21     0x4af272        b801000000      mov eax, 0x1 <== 返回值是用eax来存的
        loop3.go:21     0x4af277        4883c408        add rsp, 0x8
        loop3.go:21     0x4af27b        5d              pop rbp
        loop3.go:21     0x4af27c        c3              ret
```

至此，通过这个实例演示了如何设置寄存器值，我们将在 [hitzhangjie/godbg](https://github.com/hitzhangjie/godbg) 中实现godbg> `set reg value` 命令来修改寄存器值。
