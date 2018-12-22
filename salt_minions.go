package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nlopes/slack"
)

type minionsResponse struct {
	Minions []map[string]minion `json:"return"`
}

// Some fields from minions are not used atm cause its ton of informations to add to slack channel
// when lot of minions are connected to master
// Will need to optimize this output later
type minion struct {
	// SSDs            []string `json:"SSDs"`
	// Biosreleasedate string   `json:"biosreleasedate"`
	// Biosversion     string   `json:"biosversion"`
	CPUFlags []string `json:"cpu_flags"`
	CPUModel string   `json:"cpu_model"`
	Cpuarch  string   `json:"cpuarch"`
	Disks    []string `json:"disks"`
	// DNS             struct {
	// 	Domain         string        `json:"domain"`
	// 	IP4Nameservers []string      `json:"ip4_nameservers"`
	// 	IP6Nameservers []interface{} `json:"ip6_nameservers"`
	// 	Nameservers    []string      `json:"nameservers"`
	// 	Options        []interface{} `json:"options"`
	// 	Search         []string      `json:"search"`
	// 	Sortlist       []interface{} `json:"sortlist"`
	// } `json:"dns"`
	Domain string `json:"domain"`
	Fqdn   string `json:"fqdn"`
	// FqdnIP4 []string      `json:"fqdn_ip4"`
	// FqdnIP6 []interface{} `json:"fqdn_ip6"`
	// Gid     int           `json:"gid"`
	// Gpus    []struct {
	// 	Model  string `json:"model"`
	// 	Vendor string `json:"vendor"`
	// } `json:"gpus"`
	// Groupname        string `json:"groupname"`
	Host string `json:"host"`
	// HwaddrInterfaces struct {
	// 	Docker0 string `json:"docker0"`
	// 	Ens5    string `json:"ens5"`
	// 	Lo      string `json:"lo"`
	// } `json:"hwaddr_interfaces"`
	ID string `json:"id"`
	// Init          string `json:"init"`
	// IP4Gw         string `json:"ip4_gw"`
	// IP4Interfaces struct {
	// 	Docker0 []string `json:"docker0"`
	// 	Ens5    []string `json:"ens5"`
	// 	Lo      []string `json:"lo"`
	// } `json:"ip4_interfaces"`
	// IP6Gw         bool `json:"ip6_gw"`
	// IP6Interfaces struct {
	// 	Docker0 []interface{} `json:"docker0"`
	// 	Ens5    []string      `json:"ens5"`
	// 	Lo      []string      `json:"lo"`
	// } `json:"ip6_interfaces"`
	// IPGw         bool `json:"ip_gw"`
	// IPInterfaces struct {
	// 	Docker0 []string `json:"docker0"`
	// 	Ens5    []string `json:"ens5"`
	// 	Lo      []string `json:"lo"`
	// } `json:"ip_interfaces"`
	Ipv4 []string `json:"ipv4"`
	// Ipv6          []string `json:"ipv6"`
	// Kernel        string   `json:"kernel"`
	// Kernelrelease string   `json:"kernelrelease"`
	// Kernelversion string   `json:"kernelversion"`
	// LocaleInfo    struct {
	// 	Defaultencoding  string `json:"defaultencoding"`
	// 	Defaultlanguage  string `json:"defaultlanguage"`
	// 	Detectedencoding string `json:"detectedencoding"`
	// } `json:"locale_info"`
	// Localhost             string        `json:"localhost"`
	// LsbDistribCodename    string        `json:"lsb_distrib_codename"`
	// LsbDistribDescription string        `json:"lsb_distrib_description"`
	// LsbDistribID          string        `json:"lsb_distrib_id"`
	// LsbDistribRelease     string        `json:"lsb_distrib_release"`
	// MachineID             string        `json:"machine_id"`
	// Manufacturer          string        `json:"manufacturer"`
	Master string `json:"master"`
	// Mdadm                 []interface{} `json:"mdadm"`
	MemTotal int `json:"mem_total"`
	// Nodename              string        `json:"nodename"`
	NumCpus  int    `json:"num_cpus"`
	NumGpus  int    `json:"num_gpus"`
	Os       string `json:"os"`
	OsFamily string `json:"os_family"`
	Osarch   string `json:"osarch"`
	// Oscodename            string        `json:"oscodename"`
	// Osfinger              string        `json:"osfinger"`
	Osfullname     string `json:"osfullname"`
	Osmajorrelease int    `json:"osmajorrelease"`
	Osrelease      string `json:"osrelease"`
	// OsreleaseInfo         []int         `json:"osrelease_info"`
	// Path                  string        `json:"path"`
	// Pid                   int           `json:"pid"`
	// Productname           string        `json:"productname"`
	// Ps                    string        `json:"ps"`
	Pythonexecutable string        `json:"pythonexecutable"`
	Pythonpath       []string      `json:"pythonpath"`
	Pythonversion    []interface{} `json:"pythonversion"`
	Saltpath         string        `json:"saltpath"`
	Saltversion      string        `json:"saltversion"`
	// Saltversioninfo []int  `json:"saltversioninfo"`
	// Serialnumber          string        `json:"serialnumber"`
	// ServerID              int           `json:"server_id"`
	Shell string `json:"shell"`
	// SwapTotal             int           `json:"swap_total"`
	// Systemd               struct {
	// 	Features string `json:"features"`
	// 	Version  string `json:"version"`
	// } `json:"systemd"`
	// UID             int    `json:"uid"`
	Username string `json:"username"`
	// UUID            string `json:"uuid"`
	// Virtual         string `json:"virtual"`
	// ZfsFeatureFlags bool   `json:"zfs_feature_flags"`
	// ZfsSupport      bool   `json:"zfs_support"`
	// Zmqversion      string `json:"zmqversion"`
}

func getMinionInfo(rtm *slack.RTM, msg *slack.MessageEvent, config botConfig) {

	sendingUser := getUserInfo(rtm, msg.User)

	if checkForSaltAdmin(sendingUser) || checkForSaltViewer(sendingUser) {
		cl, err := newSaltClient(config)
		if err != nil {
			fmt.Println("Error:", err)
		}

		resp, err := cl.getMinionsInfo("")
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

			m := minionsResponse{}

			err = json.Unmarshal(body, &m)

			rtm.SendMessage(rtm.NewOutgoingMessage("I collected informations about your minions,\nbut it's lot of informations so I sent them to you in PM.", msg.Channel))

			for minionID, details := range m.Minions[0] {

				cputotal := fmt.Sprintf("%v", details.NumCpus)
				gputotal := fmt.Sprintf("%v", details.NumGpus)
				memtotal := fmt.Sprintf("%v", details.MemTotal)

				pythonversion := interfaceToString(".", details.Pythonversion)

				attachment := slack.Attachment{
					Text:  "`" + minionID + "`",
					Color: "#3399ff",

					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Value: "*Version*: " + details.Saltversion,
						},

						slack.AttachmentField{
							Value: "*Domain*: " + details.Domain,
						},
						slack.AttachmentField{
							Value: "*FQDN*: " + details.Fqdn,
						},
						slack.AttachmentField{
							Value: "*Host*: " + details.Host,
						},
						slack.AttachmentField{
							Value: "*ID*: " + details.ID,
						},
						slack.AttachmentField{
							Value: "*IPv4*: \n  -  " + strings.Join(details.Ipv4, "\n  -  "),
						},
						slack.AttachmentField{
							Value: "*Master*: " + details.Master,
						},
						slack.AttachmentField{
							Value: "*OS*: " + details.Os,
						},
						slack.AttachmentField{
							Value: "*OS Family*: " + details.OsFamily,
						},
						slack.AttachmentField{
							Value: "*OS Arch*: " + details.Osarch,
						},
						slack.AttachmentField{
							Value: "*OS Full name*: " + details.Osfullname,
						},
						slack.AttachmentField{
							Value: "*OS Major release*: " + string(details.Osmajorrelease),
						},
						slack.AttachmentField{
							Value: "*OS Release*: " + details.Osrelease,
						},
						slack.AttachmentField{
							Value: "*CPU Flags*: \n - " + strings.Join(details.CPUFlags, "\n - "),
						},
						slack.AttachmentField{
							Value: "*CPU Model*: " + details.CPUModel,
						},
						slack.AttachmentField{
							Value: "*CPU Arch*: " + details.Cpuarch,
						},
						// slack.AttachmentField{
						// 	Value: "*Disks*: \n - " + strings.Join(details.Disks, "\n - "),
						// },
						slack.AttachmentField{
							Value: "*Memory tatal*: " + memtotal,
						},
						slack.AttachmentField{
							Value: "*CPUs*: " + cputotal,
						},
						slack.AttachmentField{
							Value: "*GPUs*: " + gputotal,
						},
						slack.AttachmentField{
							Value: "*Python exec*: " + details.Pythonexecutable,
						},
						slack.AttachmentField{
							Value: "*Python Path*:\n - " + strings.Join(details.Pythonpath, "\n - "),
						},
						slack.AttachmentField{
							Value: "*Python Version*: " + pythonversion,
						},
						slack.AttachmentField{
							Value: "*Salt path*: " + details.Saltpath,
						},
						slack.AttachmentField{
							Value: "*Shell*: " + details.Shell,
						},
						slack.AttachmentField{
							Value: "*Username*: " + details.Username,
						},
					},
				}

				rtm.PostMessage(sendingUser.ID, slack.MsgOptionAttachments(attachment))

			}

		}

	} else {
		response := "Sorry *" + sendingUser.RealName + "*, but you don't have role to run this job.\nMost likely you will have to contact your admin to sort you out!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}

func listMinions(rtm *slack.RTM, msg *slack.MessageEvent, config botConfig) {

	sendingUser := getUserInfo(rtm, msg.User)

	if checkForSaltAdmin(sendingUser) || checkForSaltViewer(sendingUser) {
		cl, err := newSaltClient(config)
		if err != nil {
			fmt.Println("Error:", err)
		}

		resp, err := cl.getMinionsInfo("")
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

			m := minionsResponse{}

			err = json.Unmarshal(body, &m)

			rtm.SendMessage(rtm.NewOutgoingMessage("Here are your minions.", msg.Channel))

			for minionID, details := range m.Minions[0] {

				pythonversion := interfaceToString(".", details.Pythonversion)

				attachment := slack.Attachment{
					Text:  "`" + minionID + "`",
					Color: "#36a64f",

					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Value: "*Version*: " + details.Saltversion,
						},
						slack.AttachmentField{
							Value: "*FQDN*: " + details.Fqdn,
						},
						slack.AttachmentField{
							Value: "*IPv4*: \n  -  " + strings.Join(details.Ipv4, "\n  -  "),
						},
						slack.AttachmentField{
							Value: "*OS Details*: " + details.Os + "||" + details.OsFamily + "||" + details.Osrelease,
						},
						slack.AttachmentField{
							Value: "*Python Version*: " + pythonversion,
						},
					},
				}

				rtm.PostMessage(msg.Channel, slack.MsgOptionAttachments(attachment))
			}

		}

	} else {
		response := "Sorry *" + sendingUser.RealName + "*, but you don't have role to run this job.\nMost likely you will have to contact your admin to sort you out!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}
