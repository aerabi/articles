ARG BUILDKIT_SBOM_SCAN_CONTEXT=true
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o hello main.go

# ------------------------------
FROM scratch

COPY --from=builder /app/hello /

CMD ["/hello"]