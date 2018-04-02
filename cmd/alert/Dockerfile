FROM scratch

HEALTHCHECK --retries=10 CMD [ "/notifier", "-c" ]

ENV ZONEINFO /zoneinfo.zip
COPY zoneinfo.zip zoneinfo.zip
COPY pkg/notifier/email.html /pkg/notifier/email.html

ENTRYPOINT [ "/notifier" ]

COPY bin/notifier /notifier
