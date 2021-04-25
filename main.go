package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/allan-simon/go-singleinstance"
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const version = "0.0.1"

var categoryNames = [...]string{
	"utility",
	"development",
	"game",
	"graphics",
	"internet-and-network",
	"office",
	"audio-video",
	"system-tools",
	"other",
}

var (
	appDirs                   []string
	configDirectory           string
	pinnedFile                string
	pinned                    []string
	leftBox                   *gtk.Box
	rightBox                  *gtk.Box
	src                       glib.SourceHandle
	refresh                   bool // we will use this to trigger rebuilding mainBox
	imgSizeScaled             int
	currentWsNum, targetWsNum int64
	win                       *gtk.Window
)

type category struct {
	Name        string
	DisplayName string
	Icon        string
}

var categories []category

type desktopEntry struct {
	DesktopID string
	Name      string
	NameLoc   string
	Icon      string
	Exec      string
	Terminal  bool
	NoDisplay bool
}

var desktopEntries []desktopEntry

// Flags
var cssFileName = flag.String("s", "menu-start.css", "Styling: css file name")
var targetOutput = flag.String("o", "", "name of Output to display the menu on")
var displayVersion = flag.Bool("v", false, "display Version information")
var autohide = flag.Bool("d", false, "auto-hiDe: close window when left")
var position = flag.String("p", "bottom", "Position: \"bottom\" or \"top\"")
var alignment = flag.String("a", "left", "Alignment: \"left\" or \"right\"")
var imgSize = flag.Int("i", 48, "Icon size")
var marginTop = flag.Int("mt", 0, "Margin Top")
var marginLeft = flag.Int("ml", 0, "Margin Left")
var marginRight = flag.Int("mr", 0, "Margin Right")
var marginBottom = flag.Int("mb", 0, "Margin Bottom")
var lang = flag.String("lang", "", "force lang, e.g. \"en\", \"pl\"")

func main() {
	flag.Parse()

	if *displayVersion {
		fmt.Printf("nwg-panel-plugin-menu version %s\n", version)
		os.Exit(0)
	}

	// Gentle SIGTERM handler thanks to reiki4040 https://gist.github.com/reiki4040/be3705f307d3cd136e85
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		for {
			s := <-signalChan
			if s == syscall.SIGTERM {
				println("SIGTERM received, bye bye!")
				gtk.MainQuit()
			}
		}
	}()

	// We want the same key/mouse binding to turn the dock off: kill the running instance and exit.
	lockFilePath := fmt.Sprintf("%s/nwg-panel-plugin-menu.lock", tempDir())
	lockFile, err := singleinstance.CreateLockFile(lockFilePath)
	if err != nil {
		pid, err := readTextFile(lockFilePath)
		if err == nil {
			i, err := strconv.Atoi(pid)
			if err == nil {
				/*if !*autohide {
					println("Running instance found, sending SIGTERM and exiting...")
					syscall.Kill(i, syscall.SIGTERM)
				} else {
					println("Already running")
				}*/
				println("Running instance found, sending SIGTERM and exiting...")
				syscall.Kill(i, syscall.SIGTERM)
			}
		}
		os.Exit(0)
	}
	defer lockFile.Close()

	// LANGUAGE
	if *lang == "" && os.Getenv("LANG") != "" {
		*lang = strings.Split(os.Getenv("LANG"), ".")[0]
	}
	println(fmt.Sprintf("lang: %s", *lang))

	// ENVIRONMENT
	configDirectory = configDir()

	cacheDirectory := cacheDir()
	if cacheDirectory == "" {
		log.Panic("Couldn't determine cache directory location")
	}

	// DATA
	pinnedFile = filepath.Join(cacheDirectory, "nwg-dock-pinned")
	cssFile := filepath.Join(configDirectory, *cssFileName)

	appDirs = getAppDirs()

	setUpCategories()
	print("Categories: ")
	for _, cat := range categories {
		print(fmt.Sprintf("%s (%s) ", cat.DisplayName, cat.Name))
	}
	println()

	desktopFiles := listDesktopFiles()
	println(fmt.Sprintf("Found %v desktop files", len(desktopFiles)))

	parseDesktopFiles(desktopFiles)

	// USER INTERFACE
	gtk.Init(nil)

	cssProvider, _ := gtk.CssProviderNew()

	err = cssProvider.LoadFromPath(cssFile)
	if err != nil {
		fmt.Printf("%s file not found, using GTK styling\n", cssFile)
	} else {
		println(fmt.Sprintf("Using style: %s\n", cssFile))
		screen, _ := gdk.ScreenGetDefault()
		gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	layershell.InitForWindow(win)

	//screenWidth := 0
	screenHeight := 0

	var output2mon map[string]*gdk.Monitor
	if *targetOutput != "" {
		// We want to assign layershell to a monitor, but we only know the output name!
		output2mon, err = mapOutputs()
		if err == nil {
			monitor := output2mon[*targetOutput]
			layershell.SetMonitor(win, monitor)

			geometry := monitor.GetGeometry()
			//screenWidth = geometry.GetWidth()
			screenHeight = geometry.GetHeight()

		} else {
			println(err)
		}
	}

	if *position == "bottom" {
		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_BOTTOM, true)
	} else {
		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_TOP, true)
	}

	if *alignment == "left" {
		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_LEFT, true)
	} else {
		layershell.SetAnchor(win, layershell.LAYER_SHELL_EDGE_RIGHT, true)
	}

	layershell.SetLayer(win, layershell.LAYER_SHELL_LAYER_TOP)

	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_TOP, *marginTop)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_LEFT, *marginLeft)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_RIGHT, *marginRight)
	layershell.SetMargin(win, layershell.LAYER_SHELL_EDGE_BOTTOM, *marginBottom)

	layershell.SetKeyboardMode(win, layershell.LAYER_SHELL_KEYBOARD_MODE_EXCLUSIVE)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Close the window on leave, but not immediately, to avoid accidental closes
	win.Connect("leave-notify-event", func() {
		if *autohide {
			src, err = glib.TimeoutAdd(uint(1000), func() bool {
				/*win.Hide()
				src = 0*/
				gtk.MainQuit()
				return false
			})
		}
	})

	win.Connect("enter-notify-event", func() {
		cancelClose()
	})

	win.SetProperty("name", "menu-start-window")

	outerBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	outerBox.SetProperty("name", "box")
	win.Add(outerBox)

	alignmentBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	outerBox.PackStart(alignmentBox, true, true, 10)

	leftBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	alignmentBox.PackStart(leftBox, true, true, 10)

	leftVBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	leftBox.PackStart(leftVBox, true, true, 0)

	categoriesListBox := setUpCategoriesList()
	leftVBox.PackStart(categoriesListBox, false, false, 0)

	searchEntry := setUpSearchEntry()
	leftVBox.PackEnd(searchEntry, false, false, 6)

	rightBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	alignmentBox.PackStart(rightBox, true, true, 0)

	l, _ := gtk.LabelNew("Right box not yet ready. Be patient.")
	rightBox.PackStart(l, true, false, 10)

	win.SetSizeRequest(0, screenHeight/2)

	win.ShowAll()
	searchEntry.GrabFocus()

	gtk.Main()
}
