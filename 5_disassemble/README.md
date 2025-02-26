课程目标：学习对tracee的指令数据进行反汇编

测试步骤如下：

1、随便写一个可无限运行的go程序，或者ps找个当前正在运行的工具或者服务，
   比如运行个top，然后ps看下其pid，
   或者直接 while [ 1 -eq 1 ]; do echo "pid $$"; sleep 1; done
2、执行 go run main.go <pid>
3、执行disass进行反汇编，查看输出的汇编指令
4、可以尝试修改汇编格式为gnu、plan9、intel格式，这部分支持仅包含在 hitzhangjie/godbg repo中
