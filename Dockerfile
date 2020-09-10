FROM alpine:3.12.0  

LABEL maintainer="zqyangchn@gmail.com" \
      description="prometheus-alertmanager-dingtalk"

EXPOSE 8000

VOLUME /tmp

ENV DINGTALK_CONFIG_FILE="/etc/prometheus-alertmanager-dingtalk/config.yaml"

ADD .example-config/config.yaml.example /etc/prometheus-alertmanager-dingtalk/config.yaml
ADD .example-config/dingtalk.tmpl.example /etc/prometheus-alertmanager-dingtalk/dingtalk.tmpl

ADD prometheus-alertmanager-dingtalk /usr/bin/prometheus-alertmanager-dingtalk

ENTRYPOINT ["/usr/bin/prometheus-alertmanager-dingtalk"]
