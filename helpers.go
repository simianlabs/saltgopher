package main

import (
	"fmt"
	"strings"
)

func interfaceToString(join string, text interface{}) string {

	fmted := fmt.Sprintf("%v", text)
	stringSliced := strings.Fields(fmted)
	stringJoined := strings.Join(stringSliced, join)

	return stringJoined
}
