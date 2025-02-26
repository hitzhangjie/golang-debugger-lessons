该实例内容比较简单，仅仅是在break命令实现时记录下已经添加的断点，然后执行clear命令时将上述断点清除。

清除断点使用的系统调用和清除断点使用的系统调用是相同的，只不过一个是清除（用旧数据覆盖0xCC），一个是用0xCC覆盖旧数据。

这部分内容详见：https://github.com/hitzhangjie/godbg/blob/master/cmd/debug/clear.go
