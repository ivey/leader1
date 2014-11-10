package main

import (
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ivey/leader1"

	"fmt"
)

func main() {
	iveybot := leader1.NewBot(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("OAUTH_TOKEN"),
		os.Getenv("OAUTH_SECRET"))

	// iveybot.Keywords = []string{"gamergate"} // keywords to search for

	iveybot.OnStartup = func(b *leader1.Bot) {
		fmt.Println("TODO: implement tweet storage")
	}

	iveybot.OnTweet = func(b *leader1.Bot, t *anaconda.Tweet) {
		fmt.Println("TWEEEET: ", t.Text)
	}

	iveybot.OnFollow = func(b *leader1.Bot, u *anaconda.User) {
		fmt.Println("FOLLOWED BY: ", u.ScreenName)
	}

	iveybot.OnMessage = func(b *leader1.Bot, m *anaconda.DirectMessage) {
		if m.Sender.ScreenName == "ivey" {
			// // if tweet Foo
			// makeTweet(b, m.Text)
			fmt.Println("I WILL OBEY")
		} else {
			fmt.Println("I need to reply to ", m.Sender.ScreenName)
		}
	}

	iveybot.Start()
}
