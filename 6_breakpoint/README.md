课程目标：启动一个进程或attach到指定进程，然后能够对进程设置断点

1. 先运行一个用来测试的被调试进程：`while [ 1 -eq 1 ]; do echo "pid $$"; sleep 1; done`
2. 先通过 `go run main.go attach <pid>` attach到上述shell进程；
3. 然后检查tracee的PC，通过指令patch修改其下条待执行指令的前1字节指令数据为0xCC;
4. 然后尝试恢复tracee的执行，验证其是否会暂停执行，从而判断断点是否生效；
