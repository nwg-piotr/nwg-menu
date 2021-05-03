get:
	go mod download github.com/gotk3/gotk3
	go mod download github.com/gotk3/gotk3/gdk
	go mod download github.com/gotk3/gotk3/glib
	go mod download github.com/dlasky/gotk3-layershell/layershell
	go mod download github.com/joshuarubin/go-sway
	go mod download github.com/allan-simon/go-singleinstance

build:
	go build -o bin/nwg-menu *.go

install:
	mkdir -p /usr/share/nwg-menu
	cp -r desktop-directories /usr/share/nwg-menu
	cp menu-start.css /usr/share/nwg-menu
	cp bin/nwg-menu /usr/bin

uninstall:
	rm /usr/bin/nwg-menu

run:
	go run *.go
