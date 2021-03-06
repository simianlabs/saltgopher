package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

type userInfo struct {
	RealName string
	ID       string
	Email    string
	Admin    bool
}

type saltUsers struct {
	Users []saltUserDetails `json:"users"`
}
type saltUserDetails struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}

func returnRolesForUser(user userInfo) []string {
	var rolesList []string

	jsonFile, err := os.Open(rolesFileName)
	if err != nil {
		fmt.Println("Loading roles:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var su saltUsers

	json.Unmarshal(byteValue, &su)

	for i := 0; i < len(su.Users); i++ {
		if su.Users[i].ID == user.ID {
			rolesList = su.Users[i].Roles
		}
	}
	return rolesList
}

func checkForGopherAdmin(user userInfo) bool {
	adminFound := false

	for _, role := range returnRolesForUser(user) {
		if role == "gopheradmin" {
			adminFound = true
			break
		}
	}
	return adminFound
}
func getUserInfo(rtm *slack.RTM, u string) userInfo {

	user, err := rtm.GetUserInfo(u)
	if err != nil {
		fmt.Printf("%s\n", err)

	}

	userInfo := userInfo{
		ID:       user.ID,
		RealName: user.Profile.RealName,
		Email:    user.Profile.Email,
		Admin:    user.IsAdmin,
	}

	return userInfo

}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string, config botConfig) {

	sendingUser := getUserInfo(rtm, msg.User)

	responceHoldOn := []string{
		"Let me check that for you!",
		"I will put your Salt master to work!",
		"Hold on!",
		"Sit tight while I am getting stuff done!",
		"I need just a little moment to do it ....",
		"I am on it!",
		"Ok, if I have to do it...",
		"Happy to gopher for you!",
	}
	r := rand.Int() % len(responceHoldOn)

	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	// General response
	acceptedHelp := map[string]bool{
		"help me!":     true,
		"help":         true,
		"i need help!": true,
		"help!":        true,
	}

	acceptedVersion := map[string]bool{
		"version":               true,
		"show me version":       true,
		"What is your version?": true,
	}

	acceptedRoles := map[string]bool{
		"my roles":              true,
		"show my roles":         true,
		"what roles do i have":  true,
		"what roles do i have?": true,
		"what are my roles?":    true,
		"what are my roles":     true,
	}

	acceptedWhoAreYou := map[string]bool{
		"who are you?": true,
	}

	// Testing and debugging call
	acceptedTestMsg := map[string]bool{
		"testmsg": true,
	}

	saltMatch, _ := regexp.MatchString("^salt *", text)
	setRoleMatch, _ := regexp.MatchString("^set role [a-z0-9A-Z]+ to *", text)
	fmt.Println("Salt match", saltMatch)
	fmt.Println("Role match", setRoleMatch)

	if saltMatch {

		rtm.SendMessage(rtm.NewOutgoingMessage(responceHoldOn[r], msg.Channel))
		newSaltResponse(rtm, msg, config)
		// saltMatch = false

	} else if setRoleMatch {

		if checkForGopherAdmin(sendingUser) {
			response = addNewRoleToUser(rtm, msg.Text)
		} else {
			response = "You are not *Gopher admin*!\nI can't allow you to this, unfortunately."
		}
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
		// setRoleMatch = false

	} else if acceptedRoles[text] {

		response = "Your current roles are: " + strings.Join(returnRolesForUser(sendingUser), " , ")
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))

	} else if acceptedVersion[text] {

		response = "My current version is: `" + version + "`"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))

	} else if acceptedWhoAreYou[text] {

		response = "Nice to meet you " + sendingUser.RealName + "!\nI am *Salt Gopher*,\n\nSimple chat bot created by <https://www.simianlabs.io|Simian Labs> to help you interact with your <https://www.saltstack.com|SaltStack> infrastructure."
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))

	} else if acceptedTestMsg[text] {

		response = "<http://www.foo.com|www.foo.com>"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))

		returnRolesForUser(sendingUser)

	} else if acceptedHelp[text] {

		botHelp(rtm, msg)

	} else {

		response = "Can't help you with that."
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))

	}

}

func botHelp(rtm *slack.RTM, msg *slack.MessageEvent) {

	attachment := slack.Attachment{
		Pretext: "This are few things I can help you with:",
		Text:    " _Commands are not case sensitive._",

		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "-----------------------------",
				// Value: "Print current version of bot",
			},
			slack.AttachmentField{
				Title: "\"Version\" || Alt: \"Show me version\"",
				Value: "Print current version of bot",
			},
			slack.AttachmentField{
				Title: "\"salt “target” module.function argument {kwargs}\"",
				Value: "Send job to SaltMaster via salt-api and wait for result.",
			},
			slack.AttachmentField{
				Title: "\"set role ROLENAME to @USER\"",
				Value: "Add role to user. (Only *adminsalt* and *gopheradmin* are in use!)",
			},
			slack.AttachmentField{
				Title: "\"Show my roles\"",
				Value: "Print out which roles you have currently assigned.",
			},
		},
	}

	rtm.PostMessage(msg.Channel, slack.MsgOptionAttachments(attachment))
}
