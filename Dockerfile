FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN go mod download && go build -o /app cmd/server/main.go

FROM alpine:3.19
ENV PGHOST=db PGPORT=5432 PGUSER=app PGPASS=secret PGDB=train
COPY --from=builder /app /app
CMD ["/app"]

ENV YANDEX_RASP_API_KEY=${YANDEX_RASP_API_KEY}
ENV YANDEX_RASP_TIMEOUT=${YANDEX_RASP_TIMEOUT}