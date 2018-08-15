FROM golang:1.10 as builder

ENV APP_NAME funds
ENV WORKDIR ${GOPATH}/src/github.com/ViBiOh/funds

WORKDIR ${WORKDIR}
COPY ./ ${WORKDIR}/

RUN make ${APP_NAME}-api \
 && mkdir -p /app \
 && curl -s -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && cp bin/${APP_NAME}-api /app/

FROM scratch

ENV APP_NAME funds
HEALTHCHECK --retries=10 CMD [ "/funds-api", "-url", "https://localhost:1080/health" ]

EXPOSE 1080
ENTRYPOINT [ "/funds-api" ]

COPY --from=builder /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/${APP_NAME}-api /
