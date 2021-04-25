package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func setUpPinnedListBox() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	lines, err := loadTextFile(pinnedFile)
	if err == nil {
		println(fmt.Sprintf("Loaded %v pinned items", len(pinnedFile)))
		for _, l := range lines {
			entry := id2entry[l]
			row, _ := gtk.ListBoxRowNew()
			vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
			hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
			vBox.PackStart(hBox, false, false, 2)

			pixbuf, _ := createPixbuf(entry.Icon, *iconSizeLarge)
			img, _ := gtk.ImageNewFromPixbuf(pixbuf)
			if err != nil {
				println(err, entry.Icon)
			}
			hBox.PackStart(img, false, false, 0)
			lbl, _ := gtk.LabelNew("")
			if entry.NameLoc != "" {
				lbl.SetText(entry.NameLoc)
			} else {
				lbl.SetText(entry.Name)
			}
			hBox.PackStart(lbl, false, false, 0)
			row.Add(vBox)
			listBox.Add(row)
		}
	} else {
		println(fmt.Sprintf("%s file not found", pinnedFile))
	}
	listBox.Connect("enter-notify-event", func() {
		cancelClose()
	})

	return listBox
}

func setUpCategoriesList() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	for _, cat := range categories {
		row, _ := gtk.ListBoxRowNew()
		vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		vBox.PackStart(hBox, false, false, 2)

		pixbuf, _ := createPixbuf(cat.Icon, *iconSizeLarge)
		img, _ := gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackStart(img, false, false, 0)

		lbl, _ := gtk.LabelNew(cat.DisplayName)
		hBox.PackStart(lbl, false, false, 0)

		pixbuf, _ = createPixbuf("pan-end-symbolic", *iconSizeSmall)
		img, _ = gtk.ImageNewFromPixbuf(pixbuf)
		hBox.PackEnd(img, false, false, 0)

		row.Add(vBox)
		listBox.Add(row)
	}
	listBox.Connect("enter-notify-event", func() {
		cancelClose()
	})
	return listBox
}

func setUpSearchEntry() *gtk.SearchEntry {
	searchEntry, _ := gtk.SearchEntryNew()
	searchEntry.Connect("enter-notify-event", func() {
		cancelClose()
	})

	return searchEntry
}

func setUpUserDirsList() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	userDirsMap := mapXdgUserDirs()

	row := setUpUserDirsListRow("folder-home-symbolic", "Home", "home", userDirsMap)
	listBox.Add(row)
	row = setUpUserDirsListRow("folder-documents-symbolic", "", "documents", userDirsMap)
	listBox.Add(row)
	row = setUpUserDirsListRow("folder-downloads-symbolic", "", "downloads", userDirsMap)
	listBox.Add(row)
	row = setUpUserDirsListRow("folder-music-symbolic", "", "music", userDirsMap)
	listBox.Add(row)
	row = setUpUserDirsListRow("folder-pictures-symbolic", "", "pictures", userDirsMap)
	listBox.Add(row)
	row = setUpUserDirsListRow("folder-videos-symbolic", "", "videos", userDirsMap)
	listBox.Add(row)

	listBox.Connect("enter-notify-event", func() {
		cancelClose()
	})

	return listBox
}

func setUpUserDirsListRow(iconName, displayName, entryName string, userDirsMap map[string]string) *gtk.ListBoxRow {
	if displayName == "" {
		parts := strings.Split(userDirsMap[entryName], "/")
		displayName = parts[(len(parts) - 1)]
	}
	row, _ := gtk.ListBoxRowNew()
	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	vBox.PackStart(hBox, false, false, 10)

	img, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
	hBox.PackStart(img, false, false, 0)
	lbl, _ := gtk.LabelNew(displayName)
	hBox.PackStart(lbl, false, false, 0)
	row.Add(vBox)
	row.SetTooltipText(userDirsMap[entryName])

	return row
}
