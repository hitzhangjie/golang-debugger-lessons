逐条指令指令singlestep是利用的处理器自身+内核支持，并不需要手动添加断点，这里提供个demo简单演示下。

测试方法：

1、还是用shell命令启动一个被测试进程吧，run `while [ 1 -eq 1 ]; do echo "pid $$"; sleep 1; done`

    ```bash
    pid 168732
    pid 168732
    pid 168732
    ```

2、运行 ./10_step 168732

    ```bash
    pid 168732
    pid 168732
    pid 168732
    pid 168732 <= 假定此时运行了 ./10_step 168732
    pid 168732 <= +10s
    pid 168732 <= +20s
    ```

singlestep没10微秒运行一次，大约要十几秒钟才能看到shell命令echo一次。

另外 hitzhangjie/godbg 中的实现详见：https://github.com/hitzhangjie/godbg/blob/master/cmd/debug/step.go
