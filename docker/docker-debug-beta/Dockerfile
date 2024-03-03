FROM golang:1.17.2-alpine3.14 AS builder

WORKDIR /app

COPY main.go .

RUN go build -o hello main.go

# ------------------------------
FROM scratch

COPY --from=builder /app/hello /

CMD ["/hello"]