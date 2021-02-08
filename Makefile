gen-pb:
	protoc -I proto/ proto/uploader_service.proto --go_out=plugins=grpc:proto/pb

client:
	go run cmd/client/main.go

server:
	go run cmd/server/main.go

