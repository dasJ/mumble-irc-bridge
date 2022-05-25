package main

import (
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/lrstanley/girc"
	"github.com/microcosm-cc/bluemonday"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleutil"
	_ "layeh.com/gumble/opus"
)

const (
	MUMBLE_SERVER  = "lassul.us:64738"
	MUMBLE_CHANNEL = "nixos"
	MUMBLE_USER    = "irc"
	IRC_SERVER     = "irc.eu.hackint.org"
	IRC_NICK       = "krumble"
	IRC_CHANNEL    = "#krebs"
)

func handleMumble(mumbleMsgChan chan string, ircMsgChan chan string) {
	config := gumble.NewConfig()
	config.Username = MUMBLE_USER
	blue := bluemonday.StrictPolicy()
	disconnectChannel := make(chan bool)

	config.Attach(gumbleutil.Listener{
		TextMessage: func(e *gumble.TextMessageEvent) {
			msg := blue.Sanitize(e.Message)
			if msg == "" {
				return
			}
			msg = html.UnescapeString(msg)
			if msg == "" {
				return
			}
			if strings.HasPrefix(msg, " ") {
				return
			}
			mumbleMsgChan <- fmt.Sprintf("[%s] %s\n", e.Sender.Name, msg)
		},
		Disconnect: func(e *gumble.DisconnectEvent) {
			disconnectChannel <- true
		},
	})
	for {
		fmt.Printf("Connecting to %s...\n", MUMBLE_SERVER)
		client, err := gumble.Dial(MUMBLE_SERVER, config)
		if err != nil {
			panic(err)
		}
		channel := client.Channels.Find(MUMBLE_CHANNEL)
		if channel == nil {
			panic("Channel not found")
		}
		client.Self.Move(channel)
		client.Self.SetSelfDeafened(true)
		client.Self.SetSelfMuted(true)

		go func() {
			for {
				msg := <-ircMsgChan
				if msg == "" {
					break
				}
				channel.Send(msg, false)
			}
		}()

		<-disconnectChannel
		ircMsgChan <- ""
	}
}

func handleIRC(mumbleMsgChan chan string, ircMsgChan chan string) {
	client := girc.New(girc.Config{
		Server: IRC_SERVER,
		Port:   6697,
		SSL:    true,
		Nick:   IRC_NICK,
		User:   IRC_NICK,
		Name:   "Krebs Mumble Bridge",
		//Debug:  os.Stdout,
	})

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(IRC_CHANNEL)
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		// ???
		if e.Source == nil {
			return
		}

		user := e.Source.Name
		// Do not forward our own mesages
		if user == IRC_NICK {
			return
		}
		if relayedNick, ok := e.Tags.Get("draft/relaymsg"); ok && relayedNick == IRC_NICK {
			return
		}

		// Generate HTML if necessary
		msg := e.Last()
		re := regexp.MustCompile(`(https?://[^\s]+)`)
		msg = re.ReplaceAllString(msg, "<a href=\"$1\">$1</a>")
		if msg == "" {
			return
		}
		if strings.HasPrefix(msg, " ") {
			return
		}

		ircMsgChan <- fmt.Sprintf("[%s] %s", user, girc.StripRaw(msg))
	})

	// An example of how you would add reconnect logic.
	go func() {
		for {
			if err := client.Connect(); err != nil {
				log.Printf("irc error: %s", err)

				log.Println("reconnecting to irc in 10 seconds...")
				time.Sleep(10 * time.Second)
			} else {
				return
			}
		}
	}()

	for {
		msg := <-mumbleMsgChan
		client.Cmd.Message(IRC_CHANNEL, girc.TrimFmt(msg))
	}
}

func main() {
	mumbleMsgChan := make(chan string)
	ircMsgChan := make(chan string)
	go handleMumble(mumbleMsgChan, ircMsgChan)
	go handleIRC(mumbleMsgChan, ircMsgChan)

	select {}
}
