#!/bin/bash -e

# ptrace need priviledges enabled, like `--cap-add SYS_PTRACE --security-opt seccomp=unconfined`
#docker run -it                                                              \
#-v /Users/zhangjie/debugger101/golang-debugger-lessons:/root/debugger101    \
#--name debugger.env --cap-add SYS_PTRACE --security-opt seccomp=unconfined  \
#--rm debugger.env                                                           \
#/bin/bash

# debugger need priviledges including, ptrace, etc.
docker run -it                                                              \
-v `pwd -P`:/root/debugger101                                               \
--name debugger.env --cap-add ALL                                           \
--rm debugger.env                                                           \
/bin/bash
