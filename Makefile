.PHONY: proto server client auth-client chat-client room-chat-client db-test

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/hello.proto proto/auth.proto proto/chat.proto

server:
	go run server/*.go

client:
	go run client/main.go

auth-client:
	go run cmd/auth_client/main.go

chat-client:
	go run cmd/chat_client/main.go

room-chat-client:
	go run cmd/room_chat_client/main.go

db-test:
	go run cmd/db_test/main.go