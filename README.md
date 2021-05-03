# nwg-menu

This code has initially been intended to provide the MenuStart plugin to
[nwg-panel](https://github.com/nwg-piotr/nwg-panel). Finally it's going to work standalone as well. You'll only need
to pass some command line arguments.

The `nwg-menu` command displays the system menu with simplified [freedesktop main categories](https://specifications.freedesktop.org/menu-spec/latest/apa.html) (8 instead of 13). It also provides the search entry,
to look for installed application on the basis of .desktop files, and for files in XDG user directories.

You may pin applications by right-clicking them. Pinned items will appear above the categories list. Right-click
a pinned item to unpin it. The pinned items cache is shared with the `nwggrid` command from
[nwg-launchers](https://github.com/nwg-piotr/nwg-launchers).

In the bottom-right corner of the window you'll also see a set of buttons: logout, lock screen, restart and shutdown.
The commands attached to them may be defined in the nwg-panel settings or given as the arguments.

![00.png](https://scrot.cloud/images/2021/05/03/00.png)

To use the menu standalone (e.g. with another panel/bar or with a key binding), take a look at arguments:

```text
$ nwg-menu -h
Usage of nwg-menu:
  -cmd-lock string
    	screen lock command (default "swaylock -f -c 000000")
  -cmd-logout string
    	logout command (default "swaymsg exit")
  -cmd-restart string
    	reboot command (default "systemctl reboot")
  -cmd-shutdown string
    	shutdown command (default "systemctl -i poweroff")
  -d	auto-hiDe: close window when left
  -fm string
    	File Manager (default "thunar")
  -ha string
    	Horizontal Alignment: "left" or "right" (default "left")
  -height int
    	window height
  -isl int
    	Icon Size Large (default 32)
  -iss int
    	Icon Size Small (default 16)
  -lang string
    	force lang, e.g. "en", "pl"
  -mb int
    	Margin Bottom
  -ml int
    	Margin Left
  -mr int
    	Margin Right
  -mt int
    	Margin Top
  -o string
    	name of the Output to display the menu on
  -padding uint
    	vertical item padding (default 2)
  -s string
    	Styling: css file name (default "menu-start.css")
  -term string
    	Terminal emulator (default "alacritty")
  -v	display Version information
  -va string
    	Vertical Alignment: "bottom" or "top" (default "bottom")
  -width int
    	window width
```

## Installation

### Requirements

- go 1.16
- gtk3
- gtk-layer-shell

### Steps

1. Clone the repository, cd into it.
2. Install necessary golang libraries with `make get`. First time it may take awhile, be patient.
3. `make build`
4. `sudo make install`

## Running

Plugin integration and the config GUI will be available in the nwg-panel 0.3.1 release. For now you can start the menu
from the command line / key binding. On sway, if you provide the output name, the window will be automatically scaled
to the output height * 0.6. This may look bad on vertically oriented displays: use `-width` / `height` arguments
where necessary.

## Compatibility

1. **sway** for sure
2. **other Wayland compositors** - probably

## Styling

Edit `~/.config/nwg-dock/menu-start.css` to your taste.

## Credits

This program uses some great libraries:

- [gotk3](https://github.com/gotk3/gotk3) Copyright (c) 2013-2014 Conformal Systems LLC,
Copyright (c) 2015-2018 gotk3 contributors
- [gotk3-layershell](https://github.com/dlasky/gotk3-layershell) by [@dlasky](https://github.com/dlasky/gotk3-layershell/commits?author=dlasky) - many thanks for writing this software, and for patience with my requests!
- [go-sway](https://github.com/joshuarubin/go-sway) Copyright (c) 2019 Joshua Rubin
- [go-singleinstance](github.com/allan-simon/go-singleinstance) Copyright (c) 2015 Allan Simon

The nwg-shell logo (which is also the menu button graphics) created by [SGSE](https://github.com/sgse), licensed
under the terms of the Creative Commons license [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/deed.en).
