FROM scratch

HEALTHCHECK --retries=10 CMD healthcheck

ENV ZONEINFO /zoneinfo.zip
COPY zoneinfo.zip zoneinfo.zip

ENTRYPOINT [ "/bin/sh" ]

COPY bin/notifier /bin/sh
