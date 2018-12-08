package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

const (
	configFileName = "config/config.json"
	version        = "0.2.2"
	rolesFileName  = "config/roles.json"
	adminRoleName  = "saltadmin"
)

type botConfig struct {
	Salt struct {
		URL           string `json:"url"`
		Port          int    `json:"port"`
		User          string `json:"user"`
		Password      string `json:"password"`
		Eauth         string `json:"eauth"`
		SSLSkipVerify bool   `json:"SSLSkipVerify"`
	} `json:"salt"`
	Slack struct {
		Token string `json:"token"`
	} `json:"slack"`
}

func loadConfiguration(filename string) (botConfig, error) {

	var config botConfig

	configFile, err := os.Open(configFileName)
	if err != nil {
		fmt.Printf("Loading config: %v", err)

	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	return config, err

}

func main() {
	fmt.Println("Salt Gopher starting on:", time.Now().Format("2006.01.02 15:04:05"))
	config, _ := loadConfiguration("config.json")
	fmt.Print("**************************************************** \n")

	// DEVELOP SECTION

	// DEVELOP SECTION

	api := slack.New(config.Slack.Token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					respond(rtm, ev, prefix, config)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break

			default:
				//Take no action
			}
		}
	}

}
