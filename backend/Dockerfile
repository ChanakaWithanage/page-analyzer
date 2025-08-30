FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o page-analyzer ./cmd/web

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/page-analyzer .

EXPOSE 8080

ENV PORT=8080
ENV FETCH_TIMEOUT_SECONDS=20
ENV FETCH_MAX_REDIRECTS=5
ENV FETCH_MAX_BYTES=4194304

CMD ["./page-analyzer"]
