.PHONY: proto server auth-client chat-client room-chat-client

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth.proto proto/chat.proto

server:
	go run server/*.go

auth-client:
	go run cmd/auth_client/main.go

chat-client:
	go run cmd/chat_client/main.go

room-chat-client:
	go run cmd/room_chat_client/main.go