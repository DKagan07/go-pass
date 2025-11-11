PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

build:
	go build -o gopass
	sudo mv gopass $(BINDIR)
	
test:
	go test -p 1 -count=1 ./...

remove:
	rm -rf ~/.local/gopass && rm -rf ~/.config/gopass

uninstall:
	sudo rm -f $(BINDIR)/gopass
