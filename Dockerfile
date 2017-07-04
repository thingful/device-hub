FROM bitnami/minideb:stretch

EXPOSE 50051

COPY tmp/build/linux-amd64/ /
ENTRYPOINT [ "./device-hub", "server"]
