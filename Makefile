PREFIX ?= /usr
DESTDIR ?=

get:
	go get github.com/gotk3/gotk3
	go get github.com/gotk3/gotk3/gdk
	go get github.com/gotk3/gotk3/glib
	go get github.com/dlasky/gotk3-layershell/layershell
	go get github.com/joshuarubin/go-sway
	go get github.com/allan-simon/go-singleinstance

build:
	go build -o bin/nwg-menu *.go

install:
	mkdir -p $(DESTDIR)$(PREFIX)/share/nwg-menu
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -r desktop-directories $(DESTDIR)$(PREFIX)/share/nwg-menu
	cp menu-start.css $(DESTDIR)$(PREFIX)/share/nwg-menu
	cp bin/nwg-menu $(DESTDIR)$(PREFIX)/bin/nwg-menu

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/nwg-menu
	rm -fr $(DESTDIR)$(PREFIX)/share/nwg-menu

run:
	go run *.go
