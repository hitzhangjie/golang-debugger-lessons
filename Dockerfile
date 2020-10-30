FROM hitzhangjie/linux101

USER ROOT

RUN mkdir /root/debugger101 && \
	yum install -y libdwarf-tools.x86_64

WORKDIR /root/debugger101
