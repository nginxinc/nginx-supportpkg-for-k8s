FROM alpine:latest

RUN apk --update add bind-tools curl netcat-openbsd iproute2 iperf bash && rm -rf /var/cache/apk/* \
  && ln -s /usr/bin/iperf /usr/local/bin/iperf \
  && ls -altrh /usr/local/bin/iperf

ADD https://github.com/coredns/coredns/releases/download/v1.11.4/coredns_1.11.4_linux_amd64.tgz /coredns.tgz
RUN tar -xzvf /coredns.tgz && rm -f /coredns.tgz

CMD ["bash"]