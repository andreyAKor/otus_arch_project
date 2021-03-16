# Environment
FROM golang:1.14 as build-env

RUN mkdir -p /opt/otus_arch_project
WORKDIR /opt/otus_arch_project
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/pending ./cmd/pending/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/otus_arch_project/configs/pending.yml /etc/pending.yml
COPY --from=build-env /opt/bin/pending /bin/pending

RUN mkdir -p /var/log/otus_arch_project
ENV LOGGING_FILE=/var/log/otus_arch_project/pending.log

ENTRYPOINT ["/bin/pending", "--config=/etc/pending.yml"]
