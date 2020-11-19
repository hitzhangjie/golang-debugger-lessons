FROM hitzhangjie/linux101

USER root

RUN mkdir /root/debugger101 && \
	yum install -y libdwarf-tools.x86_64 && \
	go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /root/debugger101
