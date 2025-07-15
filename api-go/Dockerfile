FROM golang:1.23.9-alpine3.21

# disable cgo to avoid gcc requirement bug
ENV CGO_ENABLED=0

RUN apk --no-cache add git entr tini make

WORKDIR /app

ENTRYPOINT ["/sbin/tini", "-g", "--"]
CMD ["go", "run", "cmd/main.go"]
# CMD ["tail", "-f", "/dev/null"]
