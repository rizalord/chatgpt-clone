# Chat GPT Clone - Auth Service

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

### Run grpc server

```bash
go run cmd/worker/main.go
```