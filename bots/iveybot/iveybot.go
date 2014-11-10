package main

import (
	"os"
	"strings"

	"github.com/ivey/leader1"

	"fmt"
)

func main() {
	iveybot := leader1.NewBot("iveybot",
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("OAUTH_TOKEN"),
		os.Getenv("OAUTH_SECRET"))

	iveybot.Keywords = []string{"gweezlebur"}

	iveybot.OnStartup = func(b *leader1.Bot) {
		b.TrainTweetCorpus("tweets/data/js/tweets")
	}

	iveybot.OnTweet = func(b *leader1.Bot, t *leader1.Tweet) {
		words := strings.Fields(t.Text)
		if words[0] == "@iveybot" {

		}
	}

	iveybot.OnFollow = func(b *leader1.Bot, u *leader1.User) {
		b.Follow(u.ScreenName)
	}

	iveybot.OnMessage = func(b *leader1.Bot, m *leader1.DirectMessage) {
		words := strings.Fields(m.Text)
		if m.Sender.ScreenName == "ivey" {
			if words[0] == "tweet" {
				if len(words) > 1 {
					b.SendTweet(b.SeededRandomText(words[1]))
				} else {
					b.SendTweet(b.RandomText())
				}
				return
			}
		}
		b.Reply(b.RandomText(), m)
	}

	iveybot.Start()
}
