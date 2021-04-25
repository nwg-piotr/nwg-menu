package main

import "github.com/gotk3/gotk3/gtk"

func setUpCategoriesList() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	for _, cat := range categories {
		row, _ := gtk.ListBoxRowNew()
		vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		vBox.PackStart(hBox, false, false, 6)
		img, _ := gtk.ImageNewFromIconName(cat.Icon, gtk.ICON_SIZE_DND)
		hBox.PackStart(img, false, false, 0)
		lbl, _ := gtk.LabelNew(cat.DisplayName)
		hBox.PackStart(lbl, false, false, 0)
		img, _ = gtk.ImageNewFromIconName("pan-end-symbolic", gtk.ICON_SIZE_MENU)
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
