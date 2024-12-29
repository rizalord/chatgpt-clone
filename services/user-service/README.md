# Chat GPT Clone - User Service

## Configuration

All configuration is in `.env` file.

## Database Migration

All database migration is in `db/migrations` folder.

### Create Migration

```shell
migrate create -ext sql -dir db/migrations create_table_xxx
```

### Run Migration

```shell
migrate -database "postgres://postgres@localhost:5432/chatgpt_clone?sslmode=disable" -path db/migrations up
migrate -database "postgres://postgres@localhost:5432/chatgpt_clone?sslmode=disable" -path db/migrations down
```

### Run Seeder
    
```shell
go run db/seeder/seeder.go
```

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