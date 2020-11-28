FROM hitzhangjie/linux101

USER root

ENV	GOPROXY=https://goproxy.cn,direct

RUN mkdir /root/debugger101 && \
	yum install -y libdwarf-tools.x86_64 && \
	yum install -y binutils.x86_64

WORKDIR /root/debugger101
