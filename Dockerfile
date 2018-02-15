FROM scratch

HEALTHCHECK --retries=10 CMD https://localhost:1080/health

EXPOSE 1080
ENTRYPOINT [ "/bin/sh" ]

COPY bin/funds /bin/sh