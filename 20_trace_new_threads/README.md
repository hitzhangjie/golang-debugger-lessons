课程目标：演示下如何跟踪多线程程序

测试方法：

1、先看看testdata/fork.c，这个程序每隔一段时间就创建一个pthread线程出来

主线程、其他线程创建出来后都会打印该线程对应的pid、tid（这里的tid就是对应的lwp的pid）

```
zhangjie🦀 testdata(master) $ ./fork 
process: 35573, thread: 35573
process: 35573, thread: 35574
process: 35573, thread: 35716
process: 35573, thread: 35853
process: 35573, thread: 35944
process: 35573, thread: 36086
process: 35573, thread: 36192
process: 35573, thread: 36295
process: 35573, thread: 36398
...
```

2、我们同时观察 ./20_trace_new_threads `<上述fork程序进程pid> 的执行情况`

```
zhangjie🦀 20_trace_new_threads(master) $ ./20_trace_new_threads 35573
===step1===: supposing running `dlv attach pid` here
process 35573 attach succ

process 35573 stopped

tracee stopped at 7f318346f098
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 35716
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 35853
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 35944
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 35944
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 35944
tracee stopped, tracee pid:35573, status: trace/breakpoint trap1, trapcause is clone: true
eventmsg: new thread lwp pid: 36086
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 36192
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 36295
tracee stopped, tracee pid:35573, status: trace/breakpoint trap, trapcause is clone: true
eventmsg: new thread lwp pid: 36398
..
```

3、20_trace_new_threads 每隔一段时间都会打印一个event msg: `<new thread LWP pid>`

结论就是，我们通过显示设置PtraceSetOptions(pid, syscall.PTRACE_O_TRACECLONE)后，恢复tracee执行，这样tracee执行起来后，当执行到clone系统调用时，就会触发一个TRAP，内核会给tracer发送一个SIGTRAP来通知tracee运行状态变化。然后tracer就可以检查对应的status数据，来判断是否是对应的clone事件。

如果是clone事件，我们可以继续通过syscall.PtraceGetEventMsg(...)来获取新clone出来的线程的LWP的pid。

检查是不是clone事件呢，参考 man 2 ptrace手册对选项PTRACE_O_TRACECLONE的介绍部分，有解释clone状况下的status值如何编码。

4、另外设置了选项PTRACE_O_TRACECLONE之后，新线程会自动被trace，所以新线程也会被暂停执行，此时如果希望新线程恢复执行，我们需要显示将其syscall.PtraceDetach或者执行syscall.PtraceContinue操作来让新线程恢复执行。
