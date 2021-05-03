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
