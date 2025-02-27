## 课程目标：演示下如何跟踪多线程程序

### 测试方法：

1、先看看testdata/fork_noquit.c，这个程序每隔一段时间就创建一个pthread线程出来

主线程、其他线程创建出来后都会打印该线程对应的pid、tid（这里的tid就是对应的lwp的pid）

>ps: fork_noquit.c 和 fork.c 的区别就是每个线程都会不停sleep(1) 永远不会退出，这么做的目的就是我们跑这个测试用时比较久，让线程不退出可以避免我们输入线程id执行attach thread 或者 switch thread1 to thread2 时出现线程已退出导致失败的情况。

下面执行该程序等待被调试器调试：

```bash
zhangjie🦀 testdata(master) $ ./fork_noquit
process: 12368, thread: 12368
process: 12368, thread: 12369
process: 12368, thread: 12527
process: 12368, thread: 12599
process: 12368, thread: 12661
...
```

2、我们同时观察 ./21_trace_old_threads `<上述fork_noquit程序进程pid>` 的执行情况

```bash
zhangjie🦀 21_trace_old_threads(master) $ ./21_trace_old_threads 12368
===step1===: check target process existed or not

===step2===: enumerate created threads by reading /proc
threads: [12368 12369 12527 12599 12661 12725 12798 12864 12934 13004 13075]    <= created thread IDs

===step3===: supposing running `dlv> thread <n>` here
12369
process 12369 attach succ                                                       <= prompt user input and attach thread
process 12369 stopped
tracee stopped at 7f06c29cf098

===step3===: supposing running `dlv> thread <n>` here
12527
switch from thread 12369 thread 12527
process 12527 attach succ                                                       <= prompt user input and switch thread
process 12527 stopped
tracee stopped at 7f06c29cf098

===step3===: supposing running `dlv> thread <n>` here

```

3、上面我们先后输入了两个线程id，第一次输入的12369，第二次输入的时12527，我们分别看下这两次输入时线程状态变化如何

最开始没有输入时，线程状态都是 S，表示Sleep，因为线程一直在做 `while(1) {sleep(1);}` 这个操作，处于sleep状态很好理解。
```bash
$ top -H -p 12368

top - 00:54:17 up 8 days,  2:10,  2 users,  load average: 0.02, 0.06, 0.08
Threads:   7 total,   0 running,   7 sleeping,   0 stopped,   0 zombie
%Cpu(s):  0.1 us,  0.1 sy,  0.0 ni, 99.8 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
MiB Mem :  31964.6 total,  26011.4 free,   4052.5 used,   1900.7 buff/cache
MiB Swap:   8192.0 total,   8192.0 free,      0.0 used.  27333.2 avail Mem

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND
12368 zhangjie  20   0   55804    888    800 S   0.0   0.0   0:00.00 fork_noquit
12369 zhangjie  20   0   55804    888    800 S   0.0   0.0   0:00.00 fork_noquit
12527 zhangjie  20   0   55804    888    800 S   0.0   0.0   0:00.00 fork_noquit
12599 zhangjie  20   0   55804    888    800 S   0.0   0.0   0:00.00 fork_noquit
12661 zhangjie  20   0   55804    888    800 S   0.0   0.0   0:00.00 fork_noquit
12725 zhangjie  20   0   55804    888    800 S   0.0   0.0   0:00.00 fork_noquit
12798 zhangjie  20   0   55804    888    800 S   0.0   0.0   0:00.00 fork_noquit
...
```

在我们输入了12369后，线程12369的状态从 S 变成了 t，表示线程现在正在被调试器调试（traced状态）

```bash
12369 zhangjie  20   0   88588    888    800 t   0.0   0.0   0:00.00 fork_noquit
```

在我们继续输入了12527之后，调试行为从跟踪线程12369变为跟踪12527,，我们看到线程12369重新从t切换为S，而12527从S切换为t

```bash
12369 zhangjie  20   0   88588    888    800 S   0.0   0.0   0:00.00 fork_noquit
12527 zhangjie  20   0   88588    888    800 t   0.0   0.0   0:00.00 fork_noquit
```

OK，ctrl+c杀死 ./21_trace_old_threads 进程，然后我们继续观察线程的状态，会自动从t变为S，因为内核会负责善后，即在tracer退出后，将所有的tracee恢复执行。
