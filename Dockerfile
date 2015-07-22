FROM        sdurrheimer/alpine-glibc
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

RUN apk add --update -t build-deps make git \
    && git clone git://github.com/jimmidyson/prometheus.git /app \
    && cd /app \
    && git checkout kubernetes-discovery \
    && make \
    && cp prometheus promtool /bin/ \
    && mkdir -p /etc/prometheus \
    && mv ./documentation/examples/prometheus-kubernetes.yml /etc/prometheus/prometheus.yml \
    && mv ./console_libraries/ ./consoles/ /etc/prometheus/ \
    && apk del --purge build-deps git make \
    && cd / \
    && rm -rf /app /var/cache/apk/*

EXPOSE     9090
VOLUME     [ "/prometheus" ]
WORKDIR    /prometheus
ENTRYPOINT [ "/bin/prometheus" ]
CMD        [ "-config.file=/etc/prometheus/prometheus.yml", \
             "-storage.local.path=/prometheus", \
             "-web.console.libraries=/etc/prometheus/console_libraries", \
             "-web.console.templates=/etc/prometheus/consoles" ]
