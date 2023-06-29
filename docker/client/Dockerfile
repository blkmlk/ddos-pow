FROM golang:1.18-alpine as builder

WORKDIR /app
ADD . .

RUN go build -v -o /out cmd/client/main.go

FROM alpine:3.15

COPY --from=builder /out /app

CMD ["/app"]
