FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/questboard ./cmd/server

FROM debian:bookworm-slim

WORKDIR /app

ENV PORT=8080
ENV QUESTBOARD_STORAGE=sqlite
ENV QUESTBOARD_DB_FILE=/data/questboard.db

RUN mkdir -p /data

COPY --from=builder /out/questboard /app/questboard
COPY web /app/web

VOLUME ["/data"]
EXPOSE 8080

CMD ["/app/questboard"]
