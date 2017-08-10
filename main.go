package main

import (
	"flag"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/nlopes/slack"
)

var (
	fileFlag      = flag.Bool("f", false, "Search Slack Files only")
	msgFlag       = flag.Bool("m", false, "Search Slack Messages only")
	q             = flag.String("q", "", "Search Query")
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

func searchAll(api *slack.Client, searchQuery string) {
	searchFiles(api, searchQuery)
	searchMessages(api, searchQuery)
}

func searchFiles(api *slack.Client, searchQuery string) {

	files, err := api.SearchFiles(searchQuery, slack.SearchParameters{Page: 1})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"page":  1,
			"query": searchQuery,
		}).Warnf("Could not search Slack Files")
	}
	for n := 1; n <= files.Paging.Pages; n++ {
		file, err := api.SearchFiles(searchQuery, slack.SearchParameters{Page: n})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"page":  n,
				"query": searchQuery,
			}).Warnf("Could not search Slack Files")
		}
		for _, f := range file.Matches {
			logrus.WithFields(logrus.Fields{
				"title":     f.Title,
				"url":       f.URLPrivate,
				"permalink": f.Permalink,
				"user":      searchUser(api, f.User),
				"channels":  searchChannel(api, f.Channels),
			}).Infof("File [Page: %d]", n)
		}

	}
}

// TODO: follow pattern in searchFiles.
func searchMessages(api *slack.Client, searchQuery string) {
	msgs, err := api.SearchMessages(searchQuery, slack.SearchParameters{Page: 1})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"page":  1,
			"query": searchQuery,
		}).Warnf("Could not search Slack Messages")
	}
	for n := 1; n <= msgs.Paging.Pages; n++ {
		msg, err := api.SearchMessages(searchQuery, slack.SearchParameters{Page: n})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"page":  n,
				"query": searchQuery,
			}).Warnf("Could not search Slack Messages")
		}
		logrus.Infof("Messages [Page: %d]\n%+v\n", n, msg)
	}
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
