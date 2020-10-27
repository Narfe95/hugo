package main

import (
	"bufio"
	"flag"
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

type Group struct {
	group huego.Group
	menu  *systray.MenuItem
}

var applicationName = "Go-hue!"

var bridge *huego.Bridge

var verbose bool

func main() {
	ParseFlags()
	db, err := scribble.New("./", nil)
	if err != nil {
		log.Fatalf("A problem occurred when initializing scribble db: %v", err)
		return
	}

	login := Login{}
	err = db.Read("hue", "login", &login)
	if err != nil {
		log.Fatalf("An error occurred when reading the database: %v", err)
		return
	}

	if login == (Login{}) {
		log.Printf("No user found in database. Creating new user.")
		bridge, err := huego.Discover()
		if err != nil {
			log.Fatalf("Could not discover bridge: %v", err)
			return
		}
		fmt.Printf("No user found, creating new. Please press the link button on top of the Hue bridge before pressing enter to continue.")
		_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			log.Fatalf("An error occurred when reading input from stdin: %v", err)
		}
		rand.Seed(time.Now().UnixNano())
		randSeed := rand.Intn(899999) + 100000
		appUser := fmt.Sprintf("go-hue-%v", randSeed)
		user, err := bridge.CreateUser(appUser)
		if err != nil {
			log.Fatalf("Unable to create user %v, %v", appUser, err)
			return
		}
		bridge = bridge.Login(user)
		err = db.Write("hue", "login", Login{Hostname: bridge.Host, User: bridge.User})
		if err != nil {
			log.Fatalf("Unable to write host information to file: %v", err)
			return
		}
	} else {
		bridge = huego.New(login.Hostname, login.User)
	}
	if verbose {
		log.Printf("Logged into %v with username %v", bridge.Host, bridge.User)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("bulb.ico"))
	systray.SetTitle(applicationName)
	systray.SetTooltip(applicationName)

	// Disable systray groups until the library has support for it in linux
	//trayGroups := systray.AddMenuItem("Groups", "Show local groups of lamps")
	//trayLights := systray.AddMenuItem("Lights", "Show local lights")
	//lightsChannel := make(chan huego.Light)
	groupsChannel := make(chan Group)

	groups, err := bridge.GetGroups()
	if err != nil {
		log.Fatalf("Could not get groups from bridge: %v", err)
		return
	}
	for i := 0; i < len(groups); i++ {
		//mGroup := trayGroups.AddSubMenuItem(groups[i].Name, groups[i].Name)
		mGroup := systray.AddMenuItem(groups[i].Name, groups[i].Name)
		if groups[i].State.On {
			mGroup.Check()
		}
		go func(menuGroup *systray.MenuItem, group huego.Group) {
			for {
				<-menuGroup.ClickedCh
				groupsChannel <- Group{group, menuGroup}
			}
		}(mGroup, groups[i])
	}

	//lights, err := bridge.GetLights()
	//if err != nil {
	//	log.Fatalf("Could not get lights from bridge: %v", err)
	//	return
	//}
	//for i := 0; i < len(lights); i++ {
	//	mLight := trayLights.AddSubMenuItem(lights[i].Name, lights[i].Name)
	//	go func(menuLight *systray.MenuItem, light huego.Light) {
	//		if light.State.On {
	//			menuLight.Check()
	//		}
	//		for {
	//			<-menuLight.ClickedCh
	//			lightsChannel <- light
	//		}
	//	}(mLight, lights[i])
	//}

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case group := <-groupsChannel:
				if group.group.State.On {
					err := group.group.Off()
					if err != nil {
						log.Printf("Unable to set state of %v to off: %v", group.group.Name, err)
					}
					group.menu.Uncheck()
				} else {
					err := group.group.On()
					if err != nil {
						log.Printf("Unable to set state of %v to on: %v", group.group.Name, err)
					}
					group.menu.Check()
				}
			//case light := <-lightsChannel:
			//	fmt.Printf("Light %v clicked.", light.Name)
			case <-mQuit.ClickedCh:
				log.Println("Quitting...")
				systray.Quit()
				return
			}
		}
	}()
}

func ParseFlags() {
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.Parse()
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Fatalf("Could not read icon file: %v", err)
	}
	return b
}

func onExit() {
	systray.Quit()
}
