FROM alpine

ENV CONFIG_FILE /app/config.yaml
ENV SERVICE_ID 0

COPY ./builds/release/service /app/service
COPY ./config/example.yaml /app/config.yaml
COPY ./builds/docker/entrypoint.sh /app/entrypoint.sh
COPY ./builds/docker/timezone /etc/localtime

WORKDIR /app
ENTRYPOINT ["/app/entrypoint.sh"]