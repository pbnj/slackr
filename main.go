package main

import (
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
	open "github.com/petermbenjamin/go-open"
)

var (
	fileFlag      = flag.Bool("f", false, "Search Slack Files only")
	msgFlag       = flag.Bool("m", false, "Search Slack Messages only")
	q             = flag.String("q", "", "Search Query")
	openFlag      = flag.Bool("open", false, "Open URLs in browser")
	slackAPIToken = os.Getenv("SLACK_API_TOKEN")
)

func main() {
	flag.Parse()

	if slackAPIToken == "" {
		logrus.Fatalf("Slack API token cannot be blank.")
	}

	if *q == "" {
		logrus.Fatalf("Query (-q) cannot be blank.")
	}

	api := slack.New(slackAPIToken)

	if *fileFlag {
		searchFiles(api, *q)
	}

	if *msgFlag {
		searchMessages(api, *q)
	}
}

func searchFiles(api *slack.Client, searchQuery string) {
	var wg sync.WaitGroup
	files, err := api.SearchFiles(searchQuery, slack.SearchParameters{Page: 1})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"page":  1,
			"query": searchQuery,
		}).Warnf("Could not search Slack Files")
	}
	for i := 1; i <= files.Paging.Pages; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			file, err := api.SearchFiles(searchQuery, slack.SearchParameters{Page: n})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"page":  n,
					"query": searchQuery,
				}).Warnf("Could not search Slack Files")
			}
			for _, f := range file.Matches {
				if *openFlag {
					open.Open(f.Permalink)
				} else {
					logrus.WithFields(logrus.Fields{
						"title":     f.Title,
						"permalink": f.Permalink,
						"user":      searchUser(api, f.User),
						"channels":  searchChannel(api, f.Channels),
					}).Infof("File %d", n)
				}
			}
		}(i)
	}
	wg.Wait()
}

func searchMessages(api *slack.Client, searchQuery string) {
	var wg sync.WaitGroup
	msgs, err := api.SearchMessages(searchQuery, slack.SearchParameters{Page: 1})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"page":  1,
			"query": searchQuery,
		}).Warnf("Could not search Slack Messages")
	}
	for i := 1; i <= msgs.Paging.Pages; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			msg, err := api.SearchMessages(searchQuery, slack.SearchParameters{Page: n})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"page":  n,
					"query": searchQuery,
				}).Warnf("Could not search Slack Messages")
			}
			for _, m := range msg.Matches {
				if *openFlag {
					open.Open(m.Permalink)
				} else {
					logrus.WithFields(logrus.Fields{
						"text":      m.Text,
						"permalink": m.Permalink,
						"user":      m.Username,
						"channel":   m.Channel.Name,
					}).Infof("Message [Page: %d]", n)
				}
			}
		}(i)
	}
	wg.Wait()
}

func searchUser(api *slack.Client, userID string) string {
	user, err := api.GetUserInfo(userID)
	if err != nil {
		logrus.Warnf("Could not search User %s", userID)
		return ""
	}
	return user.Name
}

func searchChannel(api *slack.Client, channelIDs []string) string {
	channelNames := []string{}
	for _, channelID := range channelIDs {
		channel, err := api.GetChannelInfo(channelID)
		if err != nil {
			logrus.Warnf("Could not search Channel %s", channelID)
			channelNames = append(channelNames, "")
		}
		channelNames = append(channelNames, channel.Name)
	}
	return strings.Join(channelNames, ",")
}
