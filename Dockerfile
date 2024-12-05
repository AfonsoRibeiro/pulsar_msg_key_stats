FROM golang:1.22-alpine AS buildstage

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY src/ ./src/

RUN go build -C src/main -o /app/pulsar_id_stats

FROM alpine:latest

WORKDIR /app

COPY --from=buildstage /app/pulsar_id_stats ./pulsar_id_stats

ENTRYPOINT [ "./pulsar_id_stats" ]