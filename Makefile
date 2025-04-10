.PHONY: proto server client auth-client simple-auth-client simple-login-client auth-hello-client db-test db-query db-schema db-alter

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/hello.proto proto/auth.proto

server:
	go run server/*.go

client:
	go run client/main.go

auth-client:
	go run cmd/auth_client/main.go

simple-auth-client:
	go run cmd/simple_auth_client/main.go

simple-login-client:
	go run cmd/simple_login_client/main.go

auth-hello-client:
	go run cmd/auth_hello_client/main.go

db-test:
	go run cmd/db_test/main.go

db-query:
	go run cmd/db_query/main.go

db-schema:
	go run cmd/db_schema/main.go

db-alter:
	go run cmd/db_alter/main.go