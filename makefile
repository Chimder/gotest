swag:
	swag init -g cmd/api/main.go -o ./docs

generate:
	cd proto && protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./**/*.proto

protoui:
	grpcui -plaintext localhost:50051

gpweb:
	pgweb --url ${DB_URL}