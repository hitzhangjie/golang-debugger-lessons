测试方法：

1、运行 `while [ 1 -eq 1 ]; do echo "pid $$"; sleep 1; done

    ``bash     pid 190165     pid 190165     pid 190165     pid 190165     ...     ``

2. 此时我们希望运行 ./11_continue 190165 之后，程序会询问我们一个期望加断点的地址，但是我们这里的demo没法直接给到，需要借助其他工具来获取一个有效指令地址
3. 通过 dlv attach 190165 然后disass找一个有效地址呢？不行，dlv是面向DWARF的符号级调试器，没有调试信息时不能像指令级调试器那样工作，我们的tracee是shell程序，不行；
4. 通过 gdb attach 190165 然后disass找一个有效地址呢？可行，gdb还是功能上更加偏向底层一点，没有调试信息时，也还可以执行某些低级操作

   ```bash
   $ gdb attach 190165
   ...
   (no debugging symbols found)...done.
   0x00007f79aa703e8b in waitpid () from /lib64/libc.so.6
   Missing separate debuginfos, use: yum debuginfo-install bash-4.4.20-5.el8.x86_64

   (gdb) disass

   Dump of assembler code for function waitpid:
      0x00007f79aa703e70 <+0>:     endbr64
      0x00007f79aa703e74 <+4>:     lea    0x2cd8d5(%rip),%rax        # 0x7f79aa9d1750 <__libc_multiple_threads>
      0x00007f79aa703e7b <+11>:    mov    (%rax),%eax
      0x00007f79aa703e7d <+13>:    test   %eax,%eax
      0x00007f79aa703e7f <+15>:    jne    0x7f79aa703e98 <waitpid+40>
      0x00007f79aa703e81 <+17>:    xor    %r10d,%r10d
      0x00007f79aa703e84 <+20>:    mov    $0x3d,%eax
      0x00007f79aa703e89 <+25>:    syscall
   => 0x00007f79aa703e8b <+27>:    cmp    $0xfffffffffffff000,%rax
      0x00007f79aa703e91 <+33>:    ja     0x7f79aa703ee8 <waitpid+120>
      0x00007f79aa703e93 <+35>:    retq
      0x00007f79aa703e94 <+36>:    nopl   0x0(%rax)
      0x00007f79aa703e98 <+40>:    push   %r12
      0x00007f79aa703e9a <+42>:    mov    %edx,%r12d
      0x00007f79aa703e9d <+45>:    push   %rbp
      0x00007f79aa703e9e <+46>:    mov    %rsi,%rbp
      0x00007f79aa703ea1 <+49>:    push   %rbx
      0x00007f79aa703ea2 <+50>:    mov    %edi,%ebx
      0x00007f79aa703ea4 <+52>:    sub    $0x10,%rsp
      0x00007f79aa703ea8 <+56>:    callq  0x7f79aa620190 <__libc_enable_asynccancel>
      0x00007f79aa703ead <+61>:    xor    %r10d,%r10d
      0x00007f79aa703eb0 <+64>:    mov    %r12d,%edx
      0x00007f79aa703eb3 <+67>:    mov    %rbp,%rsi
   ```

   所以我们可以选择一个指令地址 0x00007f79aa703ea8，先用这个来进行测试
5. 继续执行我们的测试

   ```bash
   ./11_continue 190165
   ===step1===: supposing running `dlv attach pid` here
   process 190165 attach succ
   process 190165 stopped
   tracee stopped at 7f79aa703e8b

   enter a address you want to add breakpoint
   0x00007f79aa703ea8
   you entered 7f79aa703ea8

   ===step2===: supposing running `dlv> break <address>` here
   add breakpoint ok

   ===step3===: supposing running `dlv> continue` here
   process 190165 stopped
   tracee stopped at 7f79aa64a8b1
   ```

   然后我们发现tracee停在的位置并不是我们想让它停的位置7f79aa703ea8，差的很远，这是什么情况呢？

   - tracee早就开始执行了，gdb attach当时执行一瞬间的指令位置detach后立马就执行结束了，我们随便拿个指令地址不一定能被执行到；
   - 我们需要找一个循环执行到的指令位置来作为断点，or 我们tracer直接启动tracee+attach一次性完成；

6、如果要改成直接启动tracee+attach的方式，还得改代码呢，我们还是再选个有效的指令地址吧

   还是用gdb attach后disass，b printf，找到个地址：0x7fd263df1970，
   对吧，我们这个shell一直在echo，这个函数应该可以被不停地执行到，就选这个位置了 …… 测试后不行，停下来位置不符合预期。

   继续gdb找，waitpid吧，sleep入手？

```bash
   (gdb) bt
    #0  0x00007fd263e8ae70 in waitpid () from /lib64/libc.so.6
    #1  0x000055a9ad8587a9 in waitchld.isra ()
    #2  0x000055a9ad859317 in sigchld_handler ()
    #3  <signal handler called>
    #4  0x00007fd263dd18b1 in sigprocmask () from /lib64/libc.so.6
    #5  0x000055a9ad859f74 in wait_for ()
    #6  0x000055a9ad848792 in execute_command_internal ()
    #7  0x000055a9ad847e6c in execute_command_internal ()
    #8  0x000055a9ad848ae6 in execute_command ()
    #9  0x000055a9ad848bb7 in execute_while_or_until ()
    #10 0x000055a9ad846ce6 in execute_command_internal ()
    #11 0x000055a9ad848ae6 in execute_command ()
    #12 0x000055a9ad830489 in reader_loop ()
    #13 0x000055a9ad82ec5b in main ()
    ```

    通过bt显示可以看到sleep 1s会调用到waitpid，那我们就用这个呗，b waitpid，然后c，然后c，发现是可以被命中这个位置的。
    那我们就用这个位置：

   ```bash
    (gdb) b waitpid
     Breakpoint 1 at 0x7fd263e8ae70
```

   继续测试./11_continue，不行，:(
   ===step2===: supposing running `dlv> break <address>` here
   read instruction data fail: input/output error

   ps：上面两个都还好说，有个问题，tracee为什么停在了这个位置呢？7f79aa64a8b1，我们有没有在这个位置添加断点 …… 停下来不意味着就都是断点。这个后面可以在展开介绍下。
   因为地址每次会变? 这个地址不应该是线性地址吗，而且也没有开asan，地址应该不会变，weired :( TODO 后续待查

7、算了，我们手动写一个循环打印的go程序来测试吧

1. 编译testdata/forloopprint.go
2. 运行testdata/forloopprint，记下输出的pid
3. dlv attach `<pid>` 然后 b time.Sleep (Breakpoint 1 set at 0x45f70e for time.Sleep() /usr/local/go/src/runtime/time.go:178)
4. 用这个地址 0x45f70e 来作为 11_continue 测试时的输入地址
5. ./11_continue `<pid>`

   ```bash
   zhangjie🦀 11_continue(master) $ ./11_continue 226046
   ===step1===: supposing running `dlv attach pid` here
   process 226046 attach succ
   process 226046 stopped
   tracee stopped at 40332e

   enter a address you want to add breakpoint
   0x45f70e
   you entered 45f70e

   ===step2===: supposing running `dlv> break <address>` here
   add breakpoint ok

   ===step3===: supposing running `dlv> continue` here
   process 226046 stopped
   tracee stopped at 45f70f
   ```
   断点位置为45f70e，执行continue后最后停下来的位置是45f70f，刚好是目标位置patch后的下一个字节位置，符合预期。

   测试结束。
