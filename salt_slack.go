package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/nlopes/slack"
)

type saltJobPostData struct {
	Function  string
	Target    string
	Arguments string
	Kwargs    string
}

func checkForSaltAdmin(user userInfo) bool {
	adminFound := false

	for _, role := range returnRolesForUser(user) {
		if role == adminRoleName {
			adminFound = true
			break
		}
	}
	return adminFound
}

func checkForSaltViewer(user userInfo) bool {
	viewerFound := false

	for _, role := range returnRolesForUser(user) {
		if role == viewRoleName {
			viewerFound = true
			break
		}
	}
	return viewerFound
}

// Respond to new salt related command
func newSaltResponse(rtm *slack.RTM, msg *slack.MessageEvent, config botConfig) {

	sendingUser := getUserInfo(rtm, msg.User)

	if checkForSaltAdmin(sendingUser) {

		cl, err := newSaltClient(config)
		if err != nil {
			fmt.Println("Error:", err)
		}

		text := msg.Text

		tgt, mod, arg, kwargs := parseSaltRequest(text)

		job := saltJobPostData{
			Target:    tgt,
			Function:  mod,
			Arguments: arg,
			Kwargs:    kwargs,
		}

		resp, err := cl.executeJob(job)
		if err != nil {
			fmt.Println("Error:", err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		// A bit formating for nicer slack output:
		jsonParsed, err := gabs.ParseJSON(body)
		bodyFormated := jsonParsed.StringIndent("", "  ")
		if len(arg) == 0 {
			arg = "-none-"
		}
		if len(kwargs) == 0 {
			kwargs = "{ none }"
		}
		// Build slack massage and add result as attachment
		attachment := slack.Attachment{
			Pretext: "Your execution results:",
			Text:    "Job: `" + mod + "` on `" + tgt + "` with arguments: `" + arg + "` and kwargs: ```" + kwargs + "```",

			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Value: bodyFormated,
				},
			},
		}

		rtm.PostMessage(msg.Channel, slack.MsgOptionAttachments(attachment))

	} else {
		response := "Sorry *" + sendingUser.RealName + "*, but you don't have role to run this job.\nMost likely you will have to contact your admin to sort you out!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}

}

// parse command from slack to get target, arguments, module and kwargs
func parseSaltRequest(t string) (target, module, arguments, kwargs string) {

	re := regexp.MustCompile(" {")

	splitForKwargs := re.Split(t, -1)

	parsed := strings.Fields(splitForKwargs[0])

	if len(splitForKwargs) >= 2 {
		fmt.Println("Kwargs: {" + splitForKwargs[1])
		kwargs = "{" + splitForKwargs[1]
	}

	target = strings.Trim(parsed[2], `'"“”‘’`)
	module = parsed[3]

	if len(parsed) > 4 {
		arguments = parsed[4]
	} else if len(parsed) <= 3 {
		arguments = ""
	}

	return target, module, arguments, kwargs
}
