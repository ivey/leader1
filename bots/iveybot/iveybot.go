package main

import (
	"os"

	"github.com/ivey/leader1"

	"fmt"
)

func main() {
	iveybot := leader1.NewBot(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("OAUTH_TOKEN"),
		os.Getenv("OAUTH_SECRET"))

	iveybot.Keywords = []string{"foo", "bar"} // keywords to search for

	iveybot.OnStartup = func(b *leader1.Bot) {
		fmt.Println("ON STARTUP %v", b)
	}

	iveybot.Start()
}
