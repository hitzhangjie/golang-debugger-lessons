课程目标：attach到一个运行进程

示例：`go run main.go attach <pid>`

我们可以bash启动一个命令，让其一直运行，然后获取其pid，并让godbg attach将其挂住，观察程序的暂停、恢复执行。

比如，我们在bash里面先执行以下命令，它会每隔一秒打印一下当前的pid，以及计数器：

```bash
$ while [ 1 -eq 1 ]; do t=`date`; echo "$t pid: $$"; sleep 1; done

Sat Nov 14 14:29:04 UTC 2020 pid: 1311
Sat Nov 14 14:29:06 UTC 2020 pid: 1311
Sat Nov 14 14:29:07 UTC 2020 pid: 1311
Sat Nov 14 14:29:08 UTC 2020 pid: 1311
Sat Nov 14 14:29:09 UTC 2020 pid: 1311
Sat Nov 14 14:29:10 UTC 2020 pid: 1311
Sat Nov 14 14:29:11 UTC 2020 pid: 1311
Sat Nov 14 14:29:12 UTC 2020 pid: 1311
Sat Nov 14 14:29:13 UTC 2020 pid: 1311
Sat Nov 14 14:29:14 UTC 2020 pid: 1311  ==> 14s
^C
```

然后我们执行命令：
```bash
$ go run main.go attach 1311

process 1311 attach succ

process 1311 wait succ, status:4991, rusage:{{12 607026} {4 42304} 43580 0 0 0 375739 348 0 68224 35656 0 0 0 29245 153787}

we're doing some debugging...           ==> 这里sleep 10s
```

执行完上述命令后，回来看shell命令的输出情况，可见其被挂起了，等了10s之后又继续恢复执行，说明detach之后又可以继续执行。

```
Sat Nov 14 14:29:04 UTC 2020 pid: 1311
Sat Nov 14 14:29:06 UTC 2020 pid: 1311
Sat Nov 14 14:29:07 UTC 2020 pid: 1311
Sat Nov 14 14:29:08 UTC 2020 pid: 1311
Sat Nov 14 14:29:09 UTC 2020 pid: 1311
Sat Nov 14 14:29:10 UTC 2020 pid: 1311
Sat Nov 14 14:29:11 UTC 2020 pid: 1311
Sat Nov 14 14:29:12 UTC 2020 pid: 1311
Sat Nov 14 14:29:13 UTC 2020 pid: 1311
Sat Nov 14 14:29:14 UTC 2020 pid: 1311  ==> 14s attached and stopped

Sat Nov 14 14:29:24 UTC 2020 pid: 1311  ==> 24s detached and continued
Sat Nov 14 14:29:25 UTC 2020 pid: 1311
Sat Nov 14 14:29:26 UTC 2020 pid: 1311
Sat Nov 14 14:29:27 UTC 2020 pid: 1311
Sat Nov 14 14:29:28 UTC 2020 pid: 1311
Sat Nov 14 14:29:29 UTC 2020 pid: 1311
^C
```

然后我们再看下我们调试器的输出，可见其attach、暂停、detach逻辑，都是正常的。

```bash
$ go run main.go attach 1311

process 1311 attach succ

process 1311 wait succ, status:4991, rusage:{{12 607026} {4 42304} 43580 0 0 0 375739 348 0 68224 35656 0 0 0 29245 153787}

we're doing some debugging...
process 1311 detach succ
```
