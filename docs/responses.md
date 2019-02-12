#### All commands to which SaltGopher respond:

`Commands are not case sensitive.`

User commands are first trimed from all leading and trailing white space and converted to all lowercase, than compared to list.

```go
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

	// Salt subcommands
	saltMatch, _ := regexp.MatchString("^salt *", text)

	saltJobDetailMatch, _ := regexp.MatchString("^job [a-z0-9A-Z]+ details", text)

	acceptedSaltGetMinionInfo := map[string]bool{
		"get minions info": true,
	}

	acceptedSaltListMinions := map[string]bool{
		"get minions list": true,
		"list minions":     true,
	}

	acceptedSaltListJobs := map[string]bool{
		"get jobs list": true,
		"list jobs":     true,
	}

	// Role subcommands
	setRoleMatch, _ := regexp.MatchString("^set role [a-z0-9A-Z]+ to *", text)

```

If you want to add your own commands, edit this section in *slack.go* before build.