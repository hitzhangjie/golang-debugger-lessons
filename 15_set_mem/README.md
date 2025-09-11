测试方法：

1、首先我们准备了一个测试程序 testdata/loop.go
   这个程序通过一个for循环每隔1s打印当前进程的pid，循环控制变量loop默认为true

2、我们先构建并运行这个程序

   ```bash
   $ cd ../testdata && make
   $ ./loop

    pid: 49701
    pid: 49701
    pid: 49701
    pid: 49701
    pid: 49701
    ...
   ```
3、然后我们借助dlv来观察变量loop的内存位置
   ```bash
   $ dlv attach 49701

   (dlv) b loop.go:11
    Breakpoint 1 set at 0x4af0f9 for main.main() ./debugger101/golang-debugger-lessons/testdata/loop.go:11
    (dlv) c
    > [Breakpoint 1] main.main() ./debugger101/golang-debugger-lessons/testdata/loop.go:11 (hits goroutine(1):1 total:1) (PC: 0x4af0f9)
         6:         "time"
         7: )
         8:
         9: func main() {
        10:         loop := true
    =>  11:         for loop {
        12:                 fmt.Println("pid:", os.Getpid())
        13:                 time.Sleep(time.Second)
        14:         }
        15: }
    (dlv) p &loop
    (*bool)(0xc0000caf17)
    (dlv) x 0xc0000caf17
    0xc0000caf17:   0x01
    ...
   ```

3、然后我们让dlv进程退出恢复loop的执行
   ```bash
   (dlv) quit
   Would you like to kill the process? [Y/n] n
   ```
4、然后我们执行自己的程序
   ```bash
   $ ./14_set_mem 49701
    ===step1===: supposing running `dlv attach pid` here
    process 49701 attach succ
    process 49701 stopped
    tracee stopped at 476203

    enter a address you want to modify data         <= input address of variable `loop`
    0xc0000caf17
    you entered c0000caf17

    enter a value you want to change to             <= input false of variable `loop`
    0x00
    you entered 0

    we'll set *(c0000caf17) = 0                     <= do loop=false

    ===step2===: supposing running `dlv> set *addr = 0xaf` here     <= do loop=false succ
    change data from 1 to 0 succ
   ```

   And we can watch the tracee output, it soon quit for for-condition `loop=true` not met.

   ```bash
    pid: 49701
    pid: 49701
    pid: 49701                       <= the tracee exit successfully for `loop=false`
    zhangjie🦀 testdata(master) $
   ```

