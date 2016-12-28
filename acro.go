package main

import (
	"fmt"
	"os/user"
	"path/filepath"
)

var root string
var term *Term
var settings *Settings

func main() {
	term = createTerm(0)

	usr, err := user.Current()
	if err != nil {
		term.writeString(fmt.Sprintf("Couldn't get the current user: %s", err))
	} else if len(usr.HomeDir) == 0 {
		term.writeString(fmt.Sprintf("Couldn't find home folder for %v", usr))
	}

	settingsFile := filepath.Join(usr.HomeDir, ".acromantula", "settings.yml")
	settings, err = initSettings(settingsFile)
	if err != nil {
		term.writeString(fmt.Sprintf("Error loading settings from %s, : %s", settingsFile, err))
	}
	term.setPrompt(settings.Settings["prompt"] + " >> ")
	term.writeString("Hit Ctrl+D to quit\n")

	defer term.restoreTerm()

	for {
		tokens, err := term.readline()

		if err != nil {
			term.writeString(fmt.Sprintf("\nExiting....%v\n", err))
			break
		}

		if len(tokens) == 0 || len(tokens[0]) == 0 {
			continue
		}

		switch tokens[0] {
		case "header":
			handleHeaders(tokens)
		case "headers":
			handleHeaders(tokens)
		case "set":
			handleSet(tokens)
		case "settings":
			handleSet(tokens)
		case "get":
			handleGet(tokens)
		default:
			term.writeString(fmt.Sprintf("Unknown command, %v\r\n", tokens[0]))
		}
	}
}

func setRoot(str string) {
	root = str
}

func handleSet(tokens []string) {

	if len(tokens) == 1 {
		for key, value := range settings.Settings {
			term.writeString(fmt.Sprintf("  %v => %v\n", key, value))
		}
		return
	}

	if len(tokens) < 3 {
		term.writeString(fmt.Sprintf("%v needs a value\n", tokens[1]))
	} else {
		settings.Settings[tokens[1]] = tokens[2]
		term.setPrompt(tokens[2] + " >> ")
	}
}

func handleGet(tokens []string) {

	url := settings.Settings["root"]
	if len(url) == 0 && len(tokens) < 2 {
		term.writeString("No root or URL specified.\n")
		return
	}

	if len(tokens) > 1 {
		url = url + tokens[1]
	}

	performGet(term, url, settings)
}

func handleHeaders(tokens []string) {

	//
	// In the case of just 'headers'
	//
	if len(tokens) == 1 {
		term.writeString("Headers\n")
		for k, v := range settings.Headers {
			if k == "Authorization" {
				v = "****************"
			}
			term.writeString(fmt.Sprintf("%v => %v\n", k, v))

		}
		return
	}

	switch tokens[1] {
	case "set":
		if len(tokens) < 4 {
			term.writeString(fmt.Sprintf("%v is missing a value\r\n", tokens[2]))
		} else {
			settings.Headers[tokens[2]] = tokens[3]
		}
	default:
		term.writeString(fmt.Sprintf("Unknown option %v\r\n", tokens[1]))
	}

}
