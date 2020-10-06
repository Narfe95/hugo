package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/gotk3/gotk3/gtk"
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
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mChange := systray.AddMenuItem("Change Me", "Change Me")
	mChecked := systray.AddMenuItem("Unchecked", "Check Me")
	mEnabled := systray.AddMenuItem("Enabled", "Enabled")
	// Sets the icon of a menu item. Only available on Mac.
	mEnabled.SetTemplateIcon(icon.Data, icon.Data)

	systray.AddMenuItem("Ignored", "Ignored")

	subMenuTop := systray.AddMenuItem("SubMenu", "SubMenu Test (top)")
	subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenu - Level 2", "SubMenu Test (middle)")
	subMenuBottom := subMenuMiddle.AddSubMenuItem("SubMenu - Level 3", "SubMenu Test (bottom)")
	subMenuBottom2 := subMenuMiddle.AddSubMenuItem("Panic!", "SubMenu Test (bottom)")

	//mUrl := systray.AddMenuItem("Open UI", "my home")

	systray.AddSeparator()
	shown := true
	toggle := func() {
		if shown {
			subMenuBottom.Check()
			subMenuBottom2.Hide()
			mQuit.Hide()
			mEnabled.Hide()
			shown = false
		} else {
			subMenuBottom.Uncheck()
			subMenuBottom2.Show()
			mQuit.Show()
			mEnabled.Show()
			shown = true
		}
	}

	for {
		select {
		case <-mChange.ClickedCh:
			mChange.SetTitle("I've Changed")
		case <-mChecked.ClickedCh:
			if mChecked.Checked() {
				mChecked.Uncheck()
				mChecked.SetTitle("Unchecked")
			} else {
				mChecked.Check()
				mChecked.SetTitle("Checked")
			}
		case <-mEnabled.ClickedCh:
			mEnabled.SetTitle("Disabled")
			mEnabled.Disable()
		//case <-mUrl.ClickedCh:
		//	open.Run("https://www.getlantern.org")
		case <-subMenuBottom2.ClickedCh:
			panic("panic button pressed")
		case <-subMenuBottom.ClickedCh:
			toggle()
		//case <-mToggle.ClickedCh:
		//	toggle()
		case <-mQuit.ClickedCh:
			systray.Quit()
			fmt.Println("Quit2 now...")
			return
		}
	}
}

func onExit() {
	systray.Quit()
}
