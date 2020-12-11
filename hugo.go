package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/amimof/huego"
	"github.com/getlantern/systray"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type ConfigStruct struct {
	Host string `json:"host"`
	User string `json:"user"`
}

type Group struct {
	group huego.Group
	menu  *systray.MenuItem
}

type Light struct {
	light huego.Light
	menu  *systray.MenuItem
}

var applicationName = "Hugo"
var bridge *huego.Bridge
var verbose bool

func main() {
	ParseFlags()

	// Check if configuration directory exists and create it if not
	configDir := xdg.ConfigHome + "/hugo"
	configFile := configDir + "/hugo.json"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if verbose {
			log.Println("Directory " + configDir + " not found. Creating it.")
		}
		err := os.Mkdir(configDir, 0700)
		if err != nil {
			log.Fatalf("Unable to create configuration directory: %v", err)
		}
	}

	// Try loading configuration from file
	conf := ConfigStruct{}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		bridge, err = createBridgeUser(configFile, conf)
		if err != nil {
			log.Fatalf("Unable to create configuration file: %v", err)
		}
	} else {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatalf("Unable to read configuration file: %v", err)
		}
		err = json.Unmarshal(file, &conf)
		if err != nil {
			log.Fatalf("Unable to read configuration file: %v", err)
		}

		// TODO: Fix some validation to check if the values are valid
		if conf.Host == "" && conf.User == "" {
			bridge, err = createBridgeUser(configFile, conf)
			if err != nil {
				log.Fatalf("Unable to create configuration file: %v", err)
			}
		} else {
			bridge = huego.New(conf.Host, conf.User)
		}
	}

	if verbose {
		log.Printf("Logged into %v with username %v", bridge.Host, bridge.User)
	}

	systray.Run(onReady, onExit)
}

func createBridgeUser(configFile string, config ConfigStruct) (*huego.Bridge, error) {
	log.Printf("No application user found.")
	bridge, err := huego.Discover()
	if err != nil {
		return bridge, err
	}

	fmt.Printf("Creating new user. Please press the link button on top of the Hue bridge before pressing enter to continue.")
	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		return bridge, err
	}
	rand.Seed(time.Now().UnixNano())
	randSeed := rand.Intn(899999) + 100000
	appUser := fmt.Sprintf(applicationName+"-%v", randSeed)
	user, err := bridge.CreateUser(appUser)
	if err != nil {
		return bridge, err
	}

	bridge = bridge.Login(user)

	config.User = bridge.User
	config.Host = bridge.Host
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return bridge, err
	}
	err = ioutil.WriteFile(configFile, file, 0600)
	if err != nil {
		return bridge, err
	}

	return bridge, nil
}

func onReady() {
	systray.SetIcon(getIcon("bulb.ico"))
	systray.SetTitle(applicationName)
	systray.SetTooltip(applicationName)

	// Disable systray groups until the library has support for it in linux
	trayGroups := systray.AddMenuItem("Groups", "Show local groups of lamps")
	trayLights := systray.AddMenuItem("Lights", "Show local lights")
	lightsChannel := make(chan Light)
	groupsChannel := make(chan Group)

	groups, err := bridge.GetGroups()
	if err != nil {
		log.Fatalf("Could not get groups from bridge: %v", err)
		return
	}
	for i := 0; i < len(groups); i++ {
		mGroup := trayGroups.AddSubMenuItem(groups[i].Name, groups[i].Name)
		//mGroup := systray.AddMenuItem(groups[i].Name, groups[i].Name)
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

	lights, err := bridge.GetLights()
	if err != nil {
		log.Fatalf("Could not get lights from bridge: %v", err)
		return
	}
	for i := 0; i < len(lights); i++ {
		mLight := trayLights.AddSubMenuItem(lights[i].Name, lights[i].Name)
		go func(menuLight *systray.MenuItem, light huego.Light) {
			if light.State.On {
				menuLight.Check()
			}
			for {
				<-menuLight.ClickedCh
				lightsChannel <- Light{light, menuLight}
			}
		}(mLight, lights[i])
	}

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
			case light := <-lightsChannel:
				if light.light.State.On {
					err := light.light.Off()
					if err != nil {
						log.Printf("Unable to set state of %v to off: %v", light.light.Name, err)
					}
					light.menu.Uncheck()
				} else {
					err := light.light.On()
					if err != nil {
						log.Printf("Unable to set state of %v to on: %v", light.light.Name, err)
					}
					light.menu.Check()
				}
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
