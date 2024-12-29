# Chat GPT Clone - API Gateway

## Configuration

All configuration is in `.env` file.

### Update Wire
```shell
wire ./app/di
```

### Generate Protobuf
```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative app/model/proto/*.proto
```

## Run Application

### Run unit test

```bash
go test -v ./...
```

### Run integration test

```bash
go test -v ./test/
```

### Run grpc clients
```bash
```bash
cd ./../user-service && go run cmd/worker/main.go
cd ./../auth-service && go run cmd/worker/main.go
cd ./../chat-service && go run cmd/worker/main.go
```

### Run web server

```bash
go run cmd/web/main.go
```