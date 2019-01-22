package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/nlopes/slack"
)

type jobsResponse struct {
	Jobs []map[string]jobsDetails `json:"return"`
}

type jobsDetails struct {
	Arguments  []string `json:"Arguments"`
	Function   string   `json:"Function"`
	StartTime  string   `json:"StartTime"`
	Target     string   `json:"Target"`
	TargetType string   `json:"Target-type"`
	User       string   `json:"User"`
}

type jobDetailsResponse struct {
	Job []jobDetails `json:"info"`
}

type jobDetails struct {
	ID         string            `json:"jid"`
	Function   string            `json:"Function"`
	Target     string            `json:"Target"`
	User       string            `json:"User"`
	StartTime  string            `json:"StartTime"`
	TargetType string            `json:"Target-Type"`
	Arguments  []string          `json:"Arguments"`
	Minions    []string          `json:"Minions"`
	Result     map[string]Result `json:"Result"`
}

// Result ...
type Result struct {
	Return struct {
		// PID     int    `json:"pid"`
		Retcode int `json:"retcode"`
		// Stdout  string `json:"stdout"`
		// Stderr  string `json:"stderr"`
		Return bool
	} `json:"return"`
}

func getJobsList(rtm *slack.RTM, msg *slack.MessageEvent, config botConfig) {

	sendingUser := getUserInfo(rtm, msg.User)

	if checkForSaltAdmin(sendingUser) || checkForSaltViewer(sendingUser) {
		cl, err := newSaltClient(config)
		if err != nil {
			fmt.Println("Error:", err)
		}

		resp, err := cl.getJobs("")
		if err != nil {
			fmt.Println("Error:", err)
		}
		if resp == nil {
			fmt.Printf("Client failed with nil: %v", err)
			errorMsg := "*Gopher Panic!* \nI think I might run into trouble: \n`" + err.Error() + "`"
			fmt.Println(errorMsg)
			rtm.SendMessage(rtm.NewOutgoingMessage(errorMsg, msg.Channel))

		} else {

			body, _ := ioutil.ReadAll(resp.Body)

			j := jobsResponse{}

			err = json.Unmarshal(body, &j)

			var keys []string
			for k := range j.Jobs[0] {
				keys = append(keys, k)
				fmt.Println(k)

			}
			sort.Sort(sort.Reverse(sort.StringSlice(keys)))

			fmt.Println(keys)

			jobsCount := len(j.Jobs[0])
			if jobsCount <= 6 {
				rtm.SendMessage(rtm.NewOutgoingMessage("Here are all jobs I could find", msg.Channel))

				for _, k := range keys {
					// fmt.Println(k, j.Jobs[0][k].StartTime)

					details := j.Jobs[0][k]

					attachment := slack.Attachment{
						// Text:  "`" + jobID + "`",
						Color: "#4d004d",

						Fields: []slack.AttachmentField{
							slack.AttachmentField{
								Value: "JID *" + k + "*: started " + details.StartTime,
							},
						},
					}

					rtm.PostMessage(msg.Channel, slack.MsgOptionAttachments(attachment))

				}

			} else {

				rtm.SendMessage(rtm.NewOutgoingMessage("Here are last six jobs, but there are more.\nI sent them to you in PM, to avoid channel fload.", msg.Channel))

				keysTen := keys[:6]

				for _, k := range keysTen {

					details := j.Jobs[0][k]

					attachment := slack.Attachment{
						Color: "#4d004d",
						Fields: []slack.AttachmentField{
							slack.AttachmentField{
								Value: "JID *" + k + "*: started " + details.StartTime,
							},
						},
					}

					rtm.PostMessage(msg.Channel, slack.MsgOptionAttachments(attachment))

				}

				for _, k := range keys {

					details := j.Jobs[0][k]

					attachment := slack.Attachment{
						Color: "#4d004d",
						Fields: []slack.AttachmentField{
							slack.AttachmentField{
								Value: "JID *" + k + "*: started " + details.StartTime,
							},
						},
					}

					rtm.PostMessage(sendingUser.ID, slack.MsgOptionAttachments(attachment))

				}

			}

		}

	} else {
		response := "Sorry *" + sendingUser.RealName + "*, but you don't have role to run this job.\nMost likely you will have to contact your admin to sort you out!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}

func getJobDetails(rtm *slack.RTM, msg *slack.MessageEvent, config botConfig) {

	sendingUser := getUserInfo(rtm, msg.User)

	if checkForSaltAdmin(sendingUser) || checkForSaltViewer(sendingUser) {
		cl, err := newSaltClient(config)
		if err != nil {
			fmt.Println("Error:", err)
		}

		JID := parseJobDetails(msg.Text)

		resp, err := cl.getJobs(JID)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if resp == nil {
			fmt.Printf("Client failed with nil: %v", err)
			errorMsg := "*Gopher Panic!* \nI think I might run into trouble: \n`" + err.Error() + "`"
			fmt.Println(errorMsg)
			rtm.SendMessage(rtm.NewOutgoingMessage(errorMsg, msg.Channel))

		} else {

			body, _ := ioutil.ReadAll(resp.Body)

			// A bit formating for nicer slack output:
			jsonParsed, err := gabs.ParseJSON(body)
			if err != nil {
				fmt.Println("Error:", err)
			}

			fmt.Println(jsonParsed)

			j := jobDetailsResponse{}

			err = json.Unmarshal(body, &j)

			jj := j.Job[0]

			var args string
			if len(jj.Arguments) == 0 {
				args = "None"
			} else {
				args = strings.Join(jj.Arguments, "")
			}
			fmt.Println(j.Job[0].Result)

			attachment := slack.Attachment{
				Text:  "`" + jj.ID + "`",
				Color: "#800040",

				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Value: "*Start time*: " + jj.StartTime,
					},
					slack.AttachmentField{
						Value: "*Start time*: " + strings.Join(jj.Minions, "\n"),
					},
					slack.AttachmentField{
						Value: "*Start time*: " + jj.Function,
					},
					slack.AttachmentField{
						Value: "*Start time*: " + args,
					},
					slack.AttachmentField{
						Value: "*Start time*: " + jj.User,
					},
					slack.AttachmentField{
						Value: "*Start time*: " + jj.Target,
					},
					slack.AttachmentField{
						Value: "*Start time*: " + jj.TargetType,
					},
				},
			}

			rtm.PostMessage(msg.Channel, slack.MsgOptionAttachments(attachment))

			// for jobID, details := range j.Jobs[0].Result {
			// }

		}

	} else {
		response := "Sorry *" + sendingUser.RealName + "*, but you don't have role to run this job.\nMost likely you will have to contact your admin to sort you out!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}

// parse command from slack to get JID
func parseJobDetails(t string) (jid string) {

	parsed := strings.Fields(t)

	return parsed[2]
}
