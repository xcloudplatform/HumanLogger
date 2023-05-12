
all: event_streamer grpc_robo_uiclient grpc_robo_server websocket_robo_server

protoc:
	@echo "Generating Go files"
	cd proto && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto

grpc_robo_uiclient: protoc
	@echo "Building grpc_robo_uiclient"
	go build -o grpc_robo_uiclient \
		github.com/ClickerAI/ClickerAI/cmd/grpc_robo_uiclient

grpc_robo_server: protoc
	@echo "Building grpc_robo_server"
	go build -o client \
		github.com/ClickerAI/ClickerAI/cmd/grpc_robo_server

websocket_robo_server: 
	@echo "Building websocket_robo_server"
	go build -o client \
		github.com/ClickerAI/ClickerAI/cmd/websocket_robo_server

event_streamer: 
	@echo "Building event_streamer"
	go build -o client \
		github.com/ClickerAI/ClickerAI/cmd/event_streamer


clean:
	go clean github.com/ClickerAI/ClickerAI/...
	rm -f event_streamer grpc_robo_uiclient grpc_robo_server websocket_robo_server

.PHONY: event_streamer grpc_robo_uiclient grpc_robo_server websocket_robo_server protoc
