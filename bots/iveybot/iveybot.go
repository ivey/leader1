package main

import (
	"os"

	"github.com/ivey/leader1"

	"fmt"
)

// func makeTweet() string {
// 	return &string{"This is a a random tweet"}
// }

func main() {
	iveybot := leader1.NewBot("iveybot",
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("OAUTH_TOKEN"),
		os.Getenv("OAUTH_SECRET"))

	iveybot.Keywords = []string{"rocksalt"} // keywords to search for

	iveybot.OnStartup = func(b *leader1.Bot) {
		b.TrainTweetCorpus("tweets/data/js/tweets")
	}

	iveybot.OnTweet = func(b *leader1.Bot, t *leader1.Tweet) {
		fmt.Println("TWEEEET: ", t.Text)
	}

	iveybot.OnFollow = func(b *leader1.Bot, u *leader1.User) {
		b.Follow(u.ScreenName)
		// fmt.Println("FOLLOWED BY: ", u.ScreenName)
	}

	iveybot.OnMessage = func(b *leader1.Bot, m *leader1.DirectMessage) {
		// if m.Sender.ScreenName == "ivey" {
		// 	// // if tweet Foo
		// 	// makeTweet(b, m.Text)
		// 	b.Reply("HI THERE", m)
		// } else {
		b.Reply(b.RandomText(), m)
		// fmt.Println("I need to reply to ", m.Sender.ScreenName)
		// }
	}

	iveybot.Start()
}
