get:
	go get github.com/gotk3/gotk3@289cfb6dbf32de11dd2a392e86de4a144ac6be48
	go get github.com/gotk3/gotk3/gdk
	go get github.com/gotk3/gotk3/glib
	go get github.com/dlasky/gotk3-layershell/layershell
	go get github.com/joshuarubin/go-sway
	go get github.com/allan-simon/go-singleinstance

build:
	go build -o bin/nwg-panel-plugin-menu *.go

install:
	cp bin/nwg-panel-plugin-menu /usr/bin

uninstall:
	rm /usr/bin/nwg-panel-plugin-menu

run:
	go run *.go