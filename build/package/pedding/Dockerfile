# Environment
FROM golang:1.14 as build-env

RUN mkdir -p /opt/otus_arch_project
WORKDIR /opt/otus_arch_project
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/pedding ./cmd/pedding/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/otus_arch_project/configs/pedding.yml /etc/pedding.yml
COPY --from=build-env /opt/bin/pedding /bin/pedding

RUN mkdir -p /var/log/otus_arch_project
ENV LOGGING_FILE=/var/log/otus_arch_project/pedding.log

ENTRYPOINT ["/bin/pedding", "--config=/etc/pedding.yml"]
