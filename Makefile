PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

build:
	go build -o gopass
	sudo mv gopass $(BINDIR)
	
test:
	go test ./... -v

remove:
	rm -rf ~/.local/gopass && rm -rf ~/.config/gopass
