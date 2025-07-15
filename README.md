# OpenTelemetry Docker

This repository contains an example to run OpenTelemetry in Docker.

## How to use it?

1. Start services:

```sh
docker compose up
```

2. Use URLs:

- [Jaeger](http://localhost:16686/search)
- [Prometheus](http://localhost:9090/graph)

3. Make a request to the Golang API to see results in Jaeger and Prometheus:

```sh
curl http://localhost:8080/
```
