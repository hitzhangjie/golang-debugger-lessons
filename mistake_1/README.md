这个错误与 2_process_attach 相关，最初写了一版，没有意识到 ptrace对tracee后续接受请求必须来自同一个tracer 的限制。

所以一开始没有使用 runtime.LockOSThread，然后执行 syscall.PtraceDetach操作时会报错，当时为了不让它报错，自己手写了函数 ptraceDetachDetach(tid, sig int)来尝试绕过这个错误，当时为了绕过这个错误，发现控制sig=1 or 0是有效果的。

尽管可以绕过，但是实际上这种用法也是错误的。最主要的还是要解决 ptrace后续请求必须来自同一个tracer 的要求。
