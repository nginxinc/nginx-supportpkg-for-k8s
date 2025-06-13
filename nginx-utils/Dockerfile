FROM alpine:latest
COPY --chmod=744 --from=project api_stats.sh /root/api_stats.sh

RUN set -ex \
  && apk --update add --no-cache \
  bind-tools curl netcat-openbsd iproute2 \
  iperf tcpdump tshark bash jq \
  && rm -rf /var/cache/apk/* \
  && ln -s /usr/bin/iperf /usr/local/bin/iperf \
  && ls -altrh /usr/local/bin/iperf

# Setting User and Home
USER root
WORKDIR /root
ENV HOSTNAME=nginx-utils

CMD ["bash"]