package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/nlopes/slack"
)

func parseAddRole(msg string) (role, user string) {

	parsed := strings.Fields(msg)
	role = parsed[3]
	user = strings.Trim(parsed[5], `<>`)
	user = strings.Trim(user, `@`)

	return role, user

}

func checkIfUserExist(user string, su saltUsers) bool {

	userExist := false

	// check if user exist in config
	for i := 0; i < len(su.Users); i++ {
		if su.Users[i].ID == user {
			userExist = true
			break
		}
	}

	return userExist
}

func addNewRoleToUser(rtm *slack.RTM, msg string) string {

	// var rolesList []string
	writeJSON := false
	roleExist := false
	var resp string
	var su saltUsers

	role, user := parseAddRole(msg)

	// Make sure that user exist in roles.json
	addNewUser(user)

	// Get info about new user / user to add role
	newUser := getUserInfo(rtm, user)

	jsonFile, err := os.Open(rolesFileName)
	if err != nil {
		fmt.Println("Loading roles:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &su)

	// add new role to user
	for i := 0; i < len(su.Users); i++ {
		if su.Users[i].ID == user {
			rolesList := su.Users[i].Roles
			for _, rolePresent := range rolesList {
				if rolePresent == role {
					roleExist = true
					break
				}
			}
			if roleExist {
				resp = "User already have this role!"
			} else {
				su.Users[i].Roles = append(su.Users[i].Roles, role)
				writeJSON = true
			}
			break
		}
	}

	if writeJSON {
		// SOme conversion magic to get nice formated json to file
		rolesJSON, _ := json.Marshal(su)
		jsonParsed, err := gabs.ParseJSON(rolesJSON)
		if err != nil {
			fmt.Println("Parsing JSON:", err)
		}
		bodyFormated := jsonParsed.StringIndent("", "  ")

		err = ioutil.WriteFile(rolesFileName, []byte(bodyFormated), 0644)

		resp = "Role *" + role + "* added for user *" + newUser.RealName + "(uid:_" + user + "_)"
	}
	return resp

}

func addNewUser(user string) {

	var sut saltUsers

	jsonFile, err := os.Open(rolesFileName)
	if err != nil {
		fmt.Println("Loading roles:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &sut)

	// check if user exist in config
	userExist := checkIfUserExist(user, sut)

	if !userExist {

		sut.Users = append(sut.Users, saltUserDetails{ID: user, Roles: []string{}})

		// SOme conversion magic to get nice formated json to file
		rolesJSON, _ := json.Marshal(sut)
		jsonParsed, err := gabs.ParseJSON(rolesJSON)
		if err != nil {
			fmt.Println("Parsing JSON:", err)
		}
		bodyFormated := jsonParsed.StringIndent("", "  ")

		// write new config file for roles
		err = ioutil.WriteFile(rolesFileName, []byte(bodyFormated), 0644)
		fmt.Println("Added:", user)

	} else {
		fmt.Println("Exist:", user)

	}

}
