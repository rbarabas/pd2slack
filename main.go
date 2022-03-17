package main

import (
	"context"
	"log"
	"os"
	"pd2slack/internal/config"
	"pd2slack/internal/slack"
	"pd2slack/internal/sync"

	"github.com/PagerDuty/go-pagerduty"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("pd2slack: ")

	config, err := config.Get()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()

	pd := sync.NewPagerDutyClient(pagerduty.NewClient(config.PagerDutyToken))
	pdGroups, err := pd.GetGroups(context.TODO())
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", pdGroups)

	sl := slack.NewSlackClient(config.SlackToken)

	for key, pdUsers := range pdGroups {
		log.Println("PD Group ID", key)
		slackGroupName := key
		slackGroupID, groupExists, err := sl.GetGroupIDbyName(ctx, slackGroupName)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		if groupExists {
			log.Printf("Pre existing GroupID: [%s]", slackGroupID)
		}

		if !groupExists {
			slackGroupID, err = sl.CreateGroup(ctx, slackGroupName)
			if err != nil {
				log.Fatalf("unable to create user group: [%s], error:[%s]", slackGroupName, err)
			}
			log.Printf("Created new group. Name:[%s], ID:[%s]", slackGroupName, slackGroupID)
		}

		emails := []string{}
		for _, user := range pdUsers {
			emails = append(emails, user.Email)
		}

		// Empty schedule
		if len(emails) == 0 {
			log.Printf("No emails found for slack group %s (ID:%s)", slackGroupName, slackGroupID)
			os.Exit(0)
		}

		err = sl.AddMembersToGroup(ctx, slackGroupID, emails...)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
}
