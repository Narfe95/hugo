package main

import (
	"bufio"
	"fmt"
	"github.com/amimof/huego"
	"github.com/getlantern/systray"
	"github.com/nanobox-io/golang-scribble"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type Login struct {
	Hostname string
	User     string
}

var applicationName = "Go-hue!"

var mQuit *systray.MenuItem
var db *scribble.Driver
var bridge *huego.Bridge

func main() {
	db, err := scribble.New("./", nil)
	if err != nil {
		log.Println("A problem occurred when initializing scribble db: ", err)
		return
	}

	login := Login{}
	err = db.Read("hue", "login", &login)
	if err != nil {
		log.Println("An error occurred when reading the database ", err)
	}

	if login == (Login{}) {
		log.Printf("No user found in database. Creating new user.")
		bridge, err := huego.Discover()
		if err != nil {
			log.Printf("Could not discover bridge: %v", err)
			return
		}
		fmt.Printf("No user found, creating new. Please press the link button on top of the Hue bridge before pressing enter to continue.")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		rand.Seed(time.Now().UnixNano())
		randSeed := rand.Intn(899999) + 100000
		appUser := fmt.Sprintf("go-hue-%v", randSeed)
		user, err := bridge.CreateUser(appUser)
		if err != nil {
			log.Printf("Unable to create user %v, %v", appUser, err)
			return
		}
		bridge = bridge.Login(user)
		err = db.Write("hue", "login", Login{Hostname: bridge.Host, User: bridge.User})
		if err != nil {
			log.Printf("Unable to write host information to file: %v", err)
			return
		}
	} else {
		bridge = huego.New(login.Hostname, login.User)
	}
	log.Printf("Logged into %v with username %v", bridge.Host, bridge.User)

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("bulb.ico"))
	systray.SetTitle(applicationName)
	systray.SetTooltip(applicationName)

	trayGroups := systray.AddMenuItem("Groups", "Show local groups of lamps")
	trayLights := systray.AddMenuItem("Lights", "Show local lights")
	var groupSubitems []*systray.MenuItem
	var lightSubitems []*systray.MenuItem

	groups, err := bridge.GetGroups()
	if err != nil {
		log.Printf("Could not get groups from bridge: %v", err)
	}
	for _, group := range groups {
		groupSubitems = append(groupSubitems, trayGroups.AddSubMenuItem(group.Name, group.Name))
	}

	lights, err := bridge.GetLights()
	if err != nil {
		log.Printf("Could not get lights from bridge: %v", err)
	}
	for _, light := range lights {
		lightSubitems = append(lightSubitems, trayLights.AddSubMenuItem(light.Name, light.Name))
	}

	systray.AddSeparator()
	mQuit = systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for _, light := range lightSubitems {
			<-light.ClickedCh
		}
	}()
	for {

		<-mQuit.ClickedCh
		log.Println("Quitting...")
		return
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
