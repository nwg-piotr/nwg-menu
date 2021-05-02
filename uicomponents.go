package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func setUpPinnedListBox() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	lines, err := loadTextFile(pinnedFile)
	if err == nil {
		println(fmt.Sprintf("Loaded %v pinned items", len(lines)))
		for _, l := range lines {
			entry := id2entry[l]

			row, _ := gtk.ListBoxRowNew()
			row.SetSelectable(false)
			row.SetCanFocus(false)
			vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

			// We need gtk.EventBox to detect mouse event
			eventBox, _ := gtk.EventBoxNew()
			hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
			eventBox.Add(hBox)
			vBox.PackStart(eventBox, false, false, *itemPadding)

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

			row.Connect("activate", func() {
				launch(entry.Exec, entry.Terminal)
			})

			eventBox.Connect("button-release-event", func(row *gtk.ListBoxRow, e *gdk.Event) bool {
				btnEvent := gdk.EventButtonNewFromEvent(e)
				if btnEvent.Button() == 1 {
					launch(entry.Exec, entry.Terminal)
					return true
				} else if btnEvent.Button() == 3 {
					println("Unpin ", entry.DesktopID)
					return true
				}
				return false
			})

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

func setUpCategoriesListBox() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	for _, cat := range categories {
		if isSupposedToShowUp(cat.Name) {
			row, _ := gtk.ListBoxRowNew()
			row.SetCanFocus(false)
			row.SetSelectable(false)
			vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
			eventBox, _ := gtk.EventBoxNew()
			hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
			eventBox.Add(hBox)
			vBox.PackStart(eventBox, false, false, *itemPadding)

			connectCategoryListBox(cat.Name, eventBox, row)

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
	}
	listBox.Connect("enter-notify-event", func() {
		cancelClose()
	})
	return listBox
}

func isSupposedToShowUp(catName string) bool {
	result := catName == "utility" && notEmpty(listUtility) ||
		catName == "development" && notEmpty(listDevelopment) ||
		catName == "game" && notEmpty(listGame) ||
		catName == "graphics" && notEmpty(listGraphics) ||
		catName == "internet-and-network" && notEmpty(listInternetAndNetwork) ||
		catName == "office" && notEmpty(listOffice) ||
		catName == "audio-video" && notEmpty(listAudioVideo) ||
		catName == "system-tools" && notEmpty(listSystemTools) ||
		catName == "other" && notEmpty(listOther)

	return result
}

func notEmpty(listCategory []string) bool {
	if len(listCategory) == 0 {
		return false
	}
	for _, desktopID := range listCategory {
		entry := id2entry[desktopID]
		if entry.NoDisplay == false {
			return true
		}
	}
	return false
}

func connectCategoryListBox(catName string, eventBox *gtk.EventBox, row *gtk.ListBoxRow) {
	var listCategory []string

	switch catName {
	case "utility":
		listCategory = listUtility
	case "development":
		listCategory = listDevelopment
	case "game":
		listCategory = listGame
	case "graphics":
		listCategory = listGraphics
	case "internet-and-network":
		listCategory = listInternetAndNetwork
	case "office":
		listCategory = listOffice
	case "audio-video":
		listCategory = listAudioVideo
	case "system-tools":
		listCategory = listSystemTools
	default:
		listCategory = listOther
	}

	eventBox.Connect("button-release-event", func(eb *gtk.EventBox, e *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(e)
		if btnEvent.Button() == 1 {
			searchEntry.SetText("")
			clearSearchResult()
			row.SetSelectable(true)
			row.SetCanFocus(false)
			categoriesListBox.SelectRow(row)
			listBox := setUpCategoryListBox(listCategory)
			if resultWindow != nil {
				resultWindow.Destroy()
			}
			resultWindow, _ = gtk.ScrolledWindowNew(nil, nil)
			resultWindow.Connect("enter-notify-event", func() {
				cancelClose()
			})
			resultWrapper.PackStart(resultWindow, true, true, 0)
			resultWindow.Add(listBox)

			userDirsListBox.Hide()
			resultWindow.ShowAll()

			return true
		}
		return false
	})
}

func setUpBackButton() *gtk.Box {
	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	vBox.PackStart(hBox, false, false, 0)
	button, _ := gtk.ButtonNew()
	button.SetCanFocus(false)
	pixbuf, _ := createPixbuf("arrow-left", *iconSizeLarge)
	image, _ := gtk.ImageNewFromPixbuf(pixbuf)
	button.SetImage(image)
	button.SetAlwaysShowImage(true)
	button.Connect("enter-notify-event", func() {
		cancelClose()
	})
	button.Connect("clicked", func(btn *gtk.Button) {
		clearSearchResult()
		searchEntry.GrabFocus()
		searchEntry.SetText("")
	})
	hBox.PackEnd(button, false, true, 0)

	sep, _ := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	sep.SetCanFocus(false)
	vBox.Add(sep)

	return vBox
}

func setUpCategoryListBox(listCategory []string) *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()

	for _, desktopID := range listCategory {
		entry := id2entry[desktopID]
		name := entry.NameLoc
		if name == "" {
			name = entry.Name
		}
		if len(name) > 30 {
			name = fmt.Sprintf("%s...", name[:27])
		}
		if !entry.NoDisplay {
			row, _ := gtk.ListBoxRowNew()
			row.SetSelectable(false)
			vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
			eventBox, _ := gtk.EventBoxNew()
			hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
			eventBox.Add(hBox)
			vBox.PackStart(eventBox, false, false, *itemPadding)

			eventBox.Connect("button-release-event", func(row *gtk.ListBoxRow, e *gdk.Event) bool {
				btnEvent := gdk.EventButtonNewFromEvent(e)
				if btnEvent.Button() == 1 {
					launch(entry.Exec, entry.Terminal)
					return true
				}
				return false
			})

			pixbuf, _ := createPixbuf(entry.Icon, *iconSizeLarge)
			img, _ := gtk.ImageNewFromPixbuf(pixbuf)
			hBox.PackStart(img, false, false, 0)

			lbl, _ := gtk.LabelNew(name)
			hBox.PackStart(lbl, false, false, 0)

			row.Add(vBox)
			listBox.Add(row)
		}
	}
	backButton.Show()
	return listBox
}

func setUpCategorySearchResult(searchPhrase string) *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()

	resultWindow, _ = gtk.ScrolledWindowNew(nil, nil)
	resultWindow.Connect("enter-notify-event", func() {
		cancelClose()
	})
	resultWrapper.PackStart(resultWindow, true, true, 0)

	counter := 0
	for _, entry := range desktopEntries {
		if len(searchPhrase) == 1 && counter > 9 {
			break
		} else if len(searchPhrase) == 2 && counter > 14 {
			break
		}
		if !entry.NoDisplay && (strings.Contains(strings.ToLower(entry.NameLoc), strings.ToLower(searchPhrase)) ||
			strings.Contains(strings.ToLower(entry.CommentLoc), strings.ToLower(searchPhrase)) ||
			strings.Contains(strings.ToLower(entry.Comment), strings.ToLower(searchPhrase))) {

			counter++

			row, _ := gtk.ListBoxRowNew()
			row.SetSelectable(false)
			vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
			eventBox, _ := gtk.EventBoxNew()
			hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
			eventBox.Add(hBox)
			vBox.PackStart(eventBox, false, false, *itemPadding)

			exec := entry.Exec
			term := entry.Terminal
			row.Connect("activate", func() {
				launch(exec, term)
			})
			eventBox.Connect("button-release-event", func(row *gtk.EventBox, e *gdk.Event) bool {
				btnEvent := gdk.EventButtonNewFromEvent(e)
				if btnEvent.Button() == 1 {
					launch(exec, term)
					return true
				}
				return false
			})

			pixbuf, _ := createPixbuf(entry.Icon, *iconSizeLarge)
			img, _ := gtk.ImageNewFromPixbuf(pixbuf)
			hBox.PackStart(img, false, false, 0)

			lbl, _ := gtk.LabelNew(entry.NameLoc)
			hBox.PackStart(lbl, false, false, 0)

			row.Add(vBox)
			listBox.Add(row)

		}
	}
	resultWindow.Add(listBox)
	resultWindow.ShowAll()
	return listBox
}

func setUpFileSearchResult() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	if fileSearchResultWindow != nil {
		fileSearchResultWindow.Destroy()
	}
	fileSearchResultWindow, _ = gtk.ScrolledWindowNew(nil, nil)
	fileSearchResultWindow.Connect("enter-notify-event", func() {
		cancelClose()
	})
	resultWrapper.PackStart(fileSearchResultWindow, true, true, 0)

	fileSearchResultWindow.Add(listBox)
	fileSearchResultWindow.ShowAll()

	return listBox
}

func walk(path string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}
	if !d.IsDir() {
		parts := strings.Split(path, "/")
		fileName := parts[len(parts)-1]
		if strings.Contains(strings.ToLower(fileName), strings.ToLower(phrase)) {
			fileSearchResults[fileName] = path
		}
	}
	return nil
}

func setUpSearchEntry() *gtk.SearchEntry {
	searchEntry, _ := gtk.SearchEntryNew()
	searchEntry.Connect("enter-notify-event", func() {
		cancelClose()
	})
	searchEntry.Connect("search-changed", func() {
		phrase, _ = searchEntry.GetText()
		if len(phrase) > 0 {
			userDirsListBox.Hide()
			backButton.Show()

			if resultWindow != nil {
				resultWindow.Destroy()
			}
			resultListBox = setUpCategorySearchResult(phrase)
			if resultListBox.GetChildren().Length() == 0 {
				resultWindow.Hide()
			}

			if len(phrase) > 2 {
				if fileSearchResultWindow != nil {
					fileSearchResultWindow.Destroy()
				}
				fileSearchResultListBox = setUpFileSearchResult()
				for key := range userDirsMap {
					if key != "home" {
						fileSearchResults = make(map[string]string)
						if len(fileSearchResults) == 0 {
							fileSearchResultListBox.Show()
						}
						filepath.WalkDir(userDirsMap[key], walk)
						searchUserDir(key)
					}
				}
				if fileSearchResultListBox.GetChildren().Length() == 0 {
					fileSearchResultWindow.Hide()
				}
			} else {
				if fileSearchResultWindow != nil {
					fileSearchResultWindow.Destroy()
				}
			}

		} else {
			clearSearchResult()
			userDirsListBox.ShowAll()
		}
	})
	searchEntry.Connect("focus-in-event", func() {
		searchEntry.SetText("")
	})

	return searchEntry
}

func searchUserDir(dir string) {
	fileSearchResults = make(map[string]string)
	filepath.WalkDir(userDirsMap[dir], walk)
	if len(fileSearchResults) > 0 {
		row := setUpUserDirsListRow(fmt.Sprintf("folder-%s-symbolic", dir), "", dir, userDirsMap)
		fileSearchResultListBox.Add(row)
		fileSearchResultListBox.ShowAll()

		for fileName, path := range fileSearchResults {
			row := setUpUserFileSearchResultRow(fileName, path)
			fileSearchResultListBox.Add(row)
		}
		fileSearchResultListBox.ShowAll()
	}
}

func setUpUserDirsList() *gtk.ListBox {
	listBox, _ := gtk.ListBoxNew()
	userDirsMap = mapXdgUserDirs()

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
	row.SetCanFocus(false)
	row.SetSelectable(false)
	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	eventBox, _ := gtk.EventBoxNew()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	eventBox.Add(hBox)
	vBox.PackStart(eventBox, false, false, *itemPadding*3)

	img, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
	hBox.PackStart(img, false, false, 0)
	lbl, _ := gtk.LabelNew(displayName)
	hBox.PackStart(lbl, false, false, 0)
	row.Add(vBox)

	row.Connect("activate", func() {
		launch(fmt.Sprintf("%s %s", *fileManager, userDirsMap[entryName]), false)
	})

	eventBox.Connect("button-release-event", func(row *gtk.ListBoxRow, e *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(e)
		if btnEvent.Button() == 1 {
			launch(fmt.Sprintf("%s %s", *fileManager, userDirsMap[entryName]), false)
			return true
		}
		return false
	})

	return row
}

func setUpUserFileSearchResultRow(fileName, filePath string) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	row.SetCanFocus(false)
	row.SetSelectable(false)
	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	eventBox, _ := gtk.EventBoxNew()
	hBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	eventBox.Add(hBox)
	vBox.PackStart(eventBox, false, false, *itemPadding)

	if len(fileName) > 37 {
		fileName = fmt.Sprintf("%s...", fileName[:35])
	}
	lbl, _ := gtk.LabelNew(fileName)
	hBox.PackStart(lbl, false, false, 0)
	row.Add(vBox)

	row.Connect("activate", func() {
		open(filePath)
	})

	eventBox.Connect("button-release-event", func(row *gtk.ListBoxRow, e *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(e)
		if btnEvent.Button() == 1 {
			open(filePath)
			return true
		}
		return false
	})

	return row
}

func setUpButtonBox() *gtk.EventBox {
	eventBox, _ := gtk.EventBoxNew()
	wrapperHbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	wrapperHbox.PackStart(box, true, true, 10)
	eventBox.Add(wrapperHbox)

	btn, _ := gtk.ButtonNew()
	pixbuf, _ := createPixbuf("system-log-out", *iconSizeLarge)
	img, _ := gtk.ImageNewFromPixbuf(pixbuf)
	btn.SetImage(img)
	btn.SetCanFocus(false)
	box.PackStart(btn, true, true, 6)
	btn.Connect("clicked", func() {
		launch(*cmdLogout, false)
	})

	btn, _ = gtk.ButtonNew()
	pixbuf, _ = createPixbuf("system-lock-screen", *iconSizeLarge)
	img, _ = gtk.ImageNewFromPixbuf(pixbuf)
	btn.SetImage(img)
	btn.SetCanFocus(false)
	box.PackStart(btn, true, true, 6)
	btn.Connect("clicked", func() {
		launch(*cmdLock, false)
	})

	btn, _ = gtk.ButtonNew()
	pixbuf, _ = createPixbuf("system-reboot", *iconSizeLarge)
	img, _ = gtk.ImageNewFromPixbuf(pixbuf)
	btn.SetImage(img)
	btn.SetCanFocus(false)
	box.PackStart(btn, true, true, 6)
	btn.Connect("clicked", func() {
		launch(*cmdRestart, false)
	})

	btn, _ = gtk.ButtonNew()
	pixbuf, _ = createPixbuf("system-shutdown", *iconSizeLarge)
	img, _ = gtk.ImageNewFromPixbuf(pixbuf)
	btn.SetImage(img)
	btn.SetCanFocus(false)
	box.PackStart(btn, true, true, 6)
	btn.Connect("clicked", func() {
		launch(*cmdShutdown, false)
	})

	eventBox.Connect("enter-notify-event", func() {
		cancelClose()
	})

	return eventBox
}

func clearSearchResult() {
	if resultWindow != nil {
		resultWindow.Destroy()
	}
	if fileSearchResultWindow != nil {
		fileSearchResultWindow.Destroy()
	}
	if userDirsListBox != nil {
		userDirsListBox.ShowAll()
	}
	if categoriesListBox != nil {
		sr := categoriesListBox.GetSelectedRow()
		if sr != nil {
			categoriesListBox.GetSelectedRow().SetSelectable(false)
		}
		categoriesListBox.UnselectAll()
	}
	backButton.Hide()
	//searchEntry.SetText("")
	//searchEntry.GrabFocus()
}
