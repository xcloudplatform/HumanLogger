
all: client server

protoc:
	@echo "Generating Go files"
	cd src/proto && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto

server: protoc
	@echo "Building server"
	go build -o server \
		github.com/ClickerAI/ClickerAI/src/server

client: protoc
	@echo "Building client"
	go build -o client \
		github.com/ClickerAI/ClickerAI/src/client

clean:
	go clean github.com/ClickerAI/ClickerAI/...
	rm -f server client

.PHONY: client server protoc
