FROM golang:1.18-alpine as builder

WORKDIR /app
ADD ../.. .

RUN set -x \
    && export CGO_ENABLED=0 \
    && go build -v -ldflags "${LDFLAGS}" -o /server cmd/server/main.go

FROM alpine:3.15

COPY --from=builder /server /server

CMD ["/server"]
