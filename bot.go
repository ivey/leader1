package leader1

import (
	"log"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/garyburd/twitterstream"
)

type Bot struct {
	Username         string
	Consumer, Access *oauth.Credentials
	API              *anaconda.TwitterApi
	Keywords         []string

	OnStartup func(*Bot)
}

func NewBot(consumerKey string, consumerSecret string, accessKey string, accessSecret string) *Bot {
	bot := &Bot{}
	bot.Consumer = &oauth.Credentials{Token: consumerKey, Secret: consumerSecret}
	bot.Access = &oauth.Credentials{Token: accessKey, Secret: accessSecret}
	return bot
}

func (b *Bot) Start() {
	anaconda.SetConsumerKey(b.Consumer.Token)
	anaconda.SetConsumerSecret(b.Consumer.Secret)
	b.API = anaconda.NewTwitterApi(b.Access.Token, b.Access.Secret)

	if b.OnStartup != nil {
		b.OnStartup(b)
	}

	oauthClient := oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}
	oauthClient.Credentials.Token = b.Consumer.Token
	oauthClient.Credentials.Secret = b.Consumer.Secret

	params := url.Values{}
	if b.Keywords != nil {
		params.Add("track", strings.Join(b.Keywords, ", "))
	}

	ts, err := twitterstream.Open(
		&oauthClient,
		b.Access,
		"https://userstream.twitter.com/1.1/user.json",
		params)
	if err != nil {
		log.Fatal(err)
	}
	defer ts.Close()
	// Loop until stream has a permanent error.
	for ts.Err() == nil {
		var t interface{}
		if err := ts.UnmarshalNext(&t); err != nil {
			log.Fatal(err)
		}
		log.Print(t)
	}
	log.Print(ts.Err)
}
