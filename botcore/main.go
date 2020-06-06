package botcore

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

const debug = false

var (
	urlRegex = regexp.MustCompile(`https?://([^ ]+)`)
)

// IRCConfig defines the settings needed to connect to a server
type IRCConfig struct {
	Nick     string
	Server   string
	Channels []string
	Verbose  bool
	Debug    bool
	TLS      bool
}

// TitleQuery used directly from titlebot
type TitleQuery struct {
	Added   int64  `json:"timestamp"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	TTL     int64  `json:"ttl"` // TTL is used to expire the item in automatically when caching
}

func handleURL(url string, e *irc.Event) {
	var query = TitleQuery{
		URL:     url,
		Channel: e.Arguments[0],
		User:    e.Source, // nick!user@host
	}

	jsonBytes, err := json.Marshal(&query)
	if err != nil {
		fmt.Errorf("Error marshaling JSON: %#v", err)
		return
	}

	// TODO: Configurable URL
	// TODO: Authentication (apikey?)
	req, err := http.NewRequest("POST", "http://localhost:8081/title", bytes.NewBuffer(jsonBytes))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Error connecting to title service: %#v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &query)
	if err != nil {
		fmt.Errorf("Unable to unmarshal JSON response: %#v", err)
		return
	}

	e.Connection.Privmsg(e.Arguments[0], query.Title)
}

// Connect to a network and join the appropriate channels
// Starts feeding messages to callback webserver
func Connect(config IRCConfig) {
	irccon := irc.IRC(config.Nick, "gobot")
	irccon.VerboseCallbackHandler = config.Verbose
	irccon.Debug = config.Debug
	irccon.UseTLS = config.TLS
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// CONNECTED
	irccon.AddCallback("001", func(e *irc.Event) {
		// Join channels for this server
		for i := 0; i < len(config.Channels); i++ {
			irccon.Join(config.Channels[i])
			fmt.Printf("Joining %s\n", config.Channels[i])
		}

	})
	// RPL_ENDOFNAMES
	irccon.AddCallback("366", func(e *irc.Event) {})
	// MOTD
	irccon.AddCallback("372", func(e *irc.Event) {})
	// END OF MOTD
	irccon.AddCallback("376", func(e *irc.Event) {})
	// RPL_NAMREPLY (List of names from channel after join is successful)
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
		// first argument is always the channel
		for _, arg := range e.Arguments[1:] {
			urls := urlRegex.FindAllString(arg, -1)
			if urls != nil {
				fmt.Printf("URLS FOUND: %#v\n", urls)
				for _, url := range urls {
					// Launch each URL handler as a goroutine
					// Yes there is a risk of DDoS or resouce exhaustion attacks,
					// but we're willing to accept that for now
					go handleURL(url, e)
				}
			}
		}
		// TODO: handle commands

		fmt.Printf("<%s> %s\n", e.Nick, strings.Join(e.Arguments[1:], " "))
	})

	err := irccon.Connect(config.Server)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}
