课程目标：启动一个待调试进程并获取其pid和执行信息

示例：

method1:

```
go run main.go exec <path/to/prog>，如`go run main.go exec ls;
```

method2:

```
go build -o main main.go && /main exec <path/to/prog>，如 ./main exec ls;
```

由于现在demo非常简单，go run 就够用，如果源文件越来越多，我们可以 go build 完成后再执行。

- 早期的go版本是不支持go run运行一个包含多个源文件的main module的
- 交心的go版本是支持go run直接运行一个main module的
