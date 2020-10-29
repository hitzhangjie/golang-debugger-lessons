#!/bin/bash -e

docker build -t debugger.env . && \
docker run -it \
-v /Users/zhangjie/debugger101/golang-debugger-lessons:/root/debugger101 \
--name debugger.env --cap-add SYS_PTRACE --security-opt seccomp=unconfined  --rm debugger.env \
/bin/bash

