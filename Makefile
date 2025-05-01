build:
	go build -o gopass && mv gopass $(HOME)/go/bin
	
test:
	go test ./... -v
