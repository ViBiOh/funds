FROM scratch

HEALTHCHECK --retries=10 CMD [ "/notifier", "-c" ]

ENV ZONEINFO /zoneinfo.zip
COPY zoneinfo.zip zoneinfo.zip

ENTRYPOINT [ "/notifier" ]

COPY bin/notifier /notifier
