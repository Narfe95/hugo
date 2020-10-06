package main

import (
	"flag"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

var Window *gtk.Window
var Label *gtk.Label
var Entry *gtk.Entry
var VPanel *gtk.Paned
var Button *gtk.Button

func main() {
	//daemon := ParseFlags()
	systray.Run(onReady, onExit)
	//gtk.Init(nil)
	//
	//// Window
	//Window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	//if err != nil {
	//	log.Fatal("Unable to create window:", err)
	//}
	//Window.SetTitle("Go-Hue!")
	//Window.Connect("destroy", func() {
	//	gtk.MainQuit()
	//})
	//Window.SetDefaultSize(800, 600)
	//
	//// Label
	//Label, err = gtk.LabelNew("Test")
	//if err != nil {
	//	log.Fatal("Unable to create label:", err)
	//}
	//
	//// Entry
	//Entry, err = gtk.EntryNew()
	//if err != nil {
	//	log.Fatal("Unable to create entry:", err)
	//}
	//
	//Button, err = gtk.ButtonNew()
	//if err != nil {
	//	log.Fatal("Unable to create button:", err)
	//}
	//Button.SetLabel("Run")
	//Button.SetSizeRequest(120, 18)
	//Button.Connect("clicked", button_onclick)
	//
	//// horizontal panel: hPanel
	//VPanel, err = gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	//if err != nil {
	//	log.Fatal("Unable to create panel:", err)
	//}
	//VPanel.Add1(Entry)
	//VPanel.Add2(Button)
	//
	//Window.Add(VPanel)
	//
	//if !daemon {
	//	Window.ShowAll()
	//}
	//
	//gtk.Main()
}

func ParseFlags() bool {
	var daemon bool
	flag.BoolVar(&daemon, "d", false, "Daemonize instead of showing the full window.")
	flag.Parse()
	return daemon
}

func button_onclick() {
	dialog, err := gtk.DialogNew()
	if err != nil {
		log.Fatal("Unable to create dialog:", err)
	}
	dialogText, err := Entry.GetText()
	if err != nil {
		log.Fatal("Unable to get text from entry:", err)
	}
	dialog.SetTitle(dialogText)
	dialog.SetDefaultSize(200, 100)
	dialog.Show()
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome")
	_ = systray.AddMenuItem("Quit", "Quit the whole app")
	//go func() {
	//	<-quitItem.ClickedCh
	//}
}

func onExit() {
	systray.Quit()
}
