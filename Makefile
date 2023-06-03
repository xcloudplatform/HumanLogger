
all: humanlogger

humanlogger: 
	go build -o humanlogger \
		humanlogger/cmd/humanlogger


clean:
	go clean humanlogger/...
	rm -f humanlogger

.PHONY: humanlogger
