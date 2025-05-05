build:
	go build -o gopass && mv gopass $(HOME)/go/bin
	
test:
	go test ./... -v

remove:
	rm -rf ~/.local/gopass && rm -rf ~/.config/gopass
