# WALLET MICROSERVICE

## Overview

GO application for store funds for each user.

The microservice runs a GRPC server as a way of internal communication with other microservices for faster performance.
Also runs an HTTP server to monitor the server.

Endpoints:

- /health (GET)

### How to run on localhost

To be able to run in localhost you must have to have a postgres database running locally and execute the init.sql located in scripts folder.

```sh
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres

go run cmd/server/main.go
go run cmd/client/main.go
```

### Test

To be able to run the tests you must have to have a postgres database running locally and execute the init.sql located in scripts folder.

Windows

```sh
sh scripts/test.sh
```

Linux

```sh
chmod +x scripts/test.sh
sh scripts/test.sh
```

### Docker

Best option to run everything in one command

```sh
docker-compose -f build/docker-compose.yaml up --build -d
go run cmd/client/main.go
```
