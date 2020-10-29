#!/bin/bash -e

# rebuild the image
docker build -t debugger.env .

# ptrace need priviledges enabled, like `--cap-add SYS_PTRACE --security-opt seccomp=unconfined`
docker run -it 															   \
-v /Users/zhangjie/debugger101/golang-debugger-lessons:/root/debugger101   \
--name debugger.env --cap-add SYS_PTRACE --security-opt seccomp=unconfined \
--rm debugger.env 														   \
/bin/bash

