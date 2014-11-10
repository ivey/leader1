package leader1

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/garyburd/twitterstream"
	"github.com/ivey/anaconda"
)

type Bot struct {
	Username         string
	Consumer, Access *oauth.Credentials
	API              *anaconda.TwitterApi
	Keywords         []string

	OnStartup func(*Bot)
	OnFollow  func(*Bot, *User)
	OnTweet   func(*Bot, *Tweet)
	OnMessage func(*Bot, *DirectMessage)
}

type Tweet struct {
	*anaconda.Tweet
}

type User struct {
	*anaconda.User
}

type DirectMessage struct {
	*anaconda.DirectMessage
}

type StreamEvent struct {
	Friends       []int64
	Event         string
	Source        *anaconda.User
	DirectMessage *anaconda.DirectMessage `json:"direct_message"`
}

func NewBot(username string, consumerKey string, consumerSecret string, accessKey string, accessSecret string) *Bot {
	bot := &Bot{Username: username}
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
		item, err := ts.Next()
		if err != nil {
			log.Print(err)
			continue
		}

		// var debug interface{}
		// err = json.Unmarshal(item, &debug)
		// if err == nil {
		// 	log.Print("DEBUG: ", debug)
		// }

		se := &StreamEvent{}
		err = json.Unmarshal(item, se)
		if err == nil {
			if se.Friends != nil {
				log.Print("friends update")
				continue
			}
			if se.DirectMessage != nil {
				if se.DirectMessage.Sender.ScreenName == b.Username {
					continue
				}
				if b.OnMessage != nil {
					b.OnMessage(b, &DirectMessage{se.DirectMessage})
				}
				continue
			}
			if se.Event != "" {
				if se.Event == "follow" {
					if b.OnFollow != nil {
						b.OnFollow(b, &User{se.Source})
					}
					continue
				}
			}
			log.Print("unhandled stream event: ", se)
			continue
		}

		t := &anaconda.Tweet{}
		err = json.Unmarshal(item, t)
		if err == nil {
			if b.OnTweet != nil {
				b.OnTweet(b, &Tweet{t})
			}
			continue
		}

		var u interface{}
		err = json.Unmarshal(item, &u)
		if err != nil {
			log.Print("WARNING!!! Error unmarshalling object in stream: ", err)
			continue
		}
		log.Print("WARNING!!! Unhandled object in stream: ", u)
	}
	log.Print(ts.Err)
}

func (b *Bot) Follow(username string) {
	user, err := b.API.PostFriendshipsCreateToScreenName(username)
	if err != nil {
		log.Print("WARNING: couldn't follow ", username, ": ", err)
	}
	log.Print("FOLLOWED ", user)
}

func (b *Bot) Reply(text string, m *DirectMessage) {
	log.Print("Replying to DM")
	b.API.PostDMToUserId(text, m.Sender.Id)
}
