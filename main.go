package main

import (
	"crypto/tls"
	"fmt"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

const serverssl = "irc.nerv.fi:6697"

const debug = false

func main() {
	channels := [...]string{"#pyfibot.test"}

	ircnick1 := "gobot"
	irccon := irc.IRC(ircnick1, "gobot")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// CONNECTED
	irccon.AddCallback("001", func(e *irc.Event) {
		// Join channels for this server
		for i := 0; i < len(channels); i++ {
			irccon.Join(channels[i])
			fmt.Printf("Joining %s\n", channels[i])
		}

	})
	// RPL_ENDOFNAMES
	irccon.AddCallback("366", func(e *irc.Event) {})
	// MOTD
	irccon.AddCallback("372", func(e *irc.Event) {})
	// END OF MOTD
	irccon.AddCallback("376", func(e *irc.Event) {})
	// RPL_NAMREPLY
	irccon.AddCallback("353", func(e *irc.Event) {
		fmt.Printf("Joined channel %s\n", e.Arguments[2])
		fmt.Printf("Users on channel: %s\n", e.Arguments[3])
	})
	irccon.AddCallback("NOTICE", func(e *irc.Event) {
		if !debug {
			return
		}
		fmt.Printf("%#v\n", e)
	})
	irccon.AddCallback("JOIN", func(e *irc.Event) {
		// autoop everyone who joins :D
		//fmt.Printf("Opping %s\n", e.Nick)
		//irccon.Mode(e.Arguments[0], fmt.Sprintf("+o %s", e.Nick))
	})

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		// TODO Command parser
		// TODO URL parser for titles
		fmt.Printf("<%s> %s\n", e.Nick, strings.Join(e.Arguments[1:], " "))
	})

	err := irccon.Connect(serverssl)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}
