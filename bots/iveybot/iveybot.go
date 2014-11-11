package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ivey/leader1"
)

var interestingWords []string
var grumpyUsers []string
var talkedTo map[string]time.Time

func interesting(words []string) bool {
	for _, word := range words {
		for _, match := range interestingWords {
			if word == match {
				return true
			}
		}
	}
	return false
}

func d100(dc int) bool {
	if rand.Intn(100)+1 > dc {
		return true
	}
	return false
}

func main() {
	interestingWords = []string{"gweezlebur", "leagueoflegends", "ebooks", "ivey", "riot"} // TODO take as a flag/env
	grumpyUsers = []string{}
	twoMinutes := time.Duration(2) * time.Minute
	tenMinutes := time.Duration(10) * time.Minute
	twoHours := time.Duration(2) * time.Hour

	iveybot := leader1.NewBot("iveybot",
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("OAUTH_TOKEN"),
		os.Getenv("OAUTH_SECRET"))

	iveybot.Keywords = []string{"gweezlebur"} // TODO: cli flag/env

	iveybot.OnStartup = func(b *leader1.Bot) {
		b.TrainTweetCorpus("tweets/data/js/tweets") // TODO: cli/flag + serialization
	}

	iveybot.OnTweet = func(b *leader1.Bot, t *leader1.Tweet) {
		for _, grump := range grumpyUsers { // TODO: move to leader1
			if t.User.ScreenName == grump {
				return
			}
		}

		// TODO: move retweet check to leader1
		if t.RetweetedStatus != nil || t.Text[0:2] == "RT" { // No RTs
			if d100(90) { // small chance of following new people
				b.Follow(t.User.ScreenName)
			}
			return
		}

		words := strings.Fields(t.Text)
		if interesting(words) {
			b.API.Favorite(t.Id)
			b.Follow(t.User.ScreenName)

			if d100(75) {
				b.API.Retweet(t.Id, true)
			}

			if d100(50) && time.Since(talkedTo[t.User.ScreenName]) > twoHours {
				b.Reply(b.SeededRandomText(words[0]), t)
			}
		} else {
			if d100(95) {
				b.API.Favorite(t.Id)
			}

			if d100(95) {
				b.API.Retweet(t.Id, true)
			}

			if d100(95) && time.Since(talkedTo[t.User.ScreenName]) > twoHours {
				b.Reply(b.SeededRandomText(words[0]), t)
			}
		}
		talkedTo[t.User.ScreenName] = time.Now()
	}

	iveybot.OnMention = func(b *leader1.Bot, t *leader1.Tweet) {
		words := strings.Fields(t.Text)
		if interesting(words) {
			b.API.Favorite(t.Id)
		}

		if t.User.ScreenName == "blippyradar" && time.Since(talkedTo[t.User.ScreenName]) > twoMinutes {
			b.Reply(b.SeededRandomText(words[0]), t) // always time for blippyradar
		} else {
			if time.Since(talkedTo[t.User.ScreenName]) > tenMinutes {
				b.Reply(b.SeededRandomText(words[0]), t)
			}
		}
		talkedTo[t.User.ScreenName] = time.Now()
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
		b.PrivateReply(b.RandomText(), m)
	}

	iveybot.Start()
}
