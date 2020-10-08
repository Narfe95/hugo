package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/getlantern/systray"
	"io/ioutil"
	"log"
)

var applicationName = "Go-hue!"
var daemon bool

var application fyne.App
var window fyne.Window
var label *widget.Label
var button *widget.Button

var mViewWindow *systray.MenuItem
var mQuit *systray.MenuItem

var winC chan bool

func main() {
	ParseFlags()
	winC = make(chan bool)
	run := true
	go func() {
		systray.Run(onReady, onExit)
	}()
	for run {
		run := <-winC
		if run {
			fyneWindow()
		} else {
			systray.Quit()
		}
	}
}

func fyneWindow() {
	// Create application
	application = app.New()
	// Create window from application
	window = application.NewWindow(applicationName)
	window.Resize(fyne.NewSize(400, 200))
	// Create a label
	label = widget.NewLabel("test")

	// Create a button
	button = widget.NewButton("Testar", buttonOnClick)

	// Set the window content to contain a vertical box which in turn contains the label the button
	window.SetContent(widget.NewVBox(label, button))
	window.ShowAndRun()
}

func buttonOnClick() {
	label.SetText("asdf")
}

func ParseFlags() {
	flag.BoolVar(&daemon, "d", false, "Daemonize instead of showing the full window.")
	flag.Parse()
}

func onReady() {
	systray.SetIcon(getIcon("bulb.ico"))
	systray.SetTitle(applicationName)
	systray.SetTooltip(applicationName)

	mViewWindow = systray.AddMenuItem("Show window", "Show the program window")
	systray.AddSeparator()
	mQuit = systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		case <-mViewWindow.ClickedCh:
			log.Println("Opening main window")
			winC <- true
		case <-mQuit.ClickedCh:
			log.Println("Quitting...")
			winC <- false
			return
		}
	}
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func onExit() {
	systray.Quit()
}
