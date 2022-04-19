package main

import (
	"context"
	"os"
	"pd2slack/internal/config"
	"pd2slack/internal/log"
	"pd2slack/internal/slack"
	"pd2slack/internal/sync"

	"github.com/PagerDuty/go-pagerduty"
)

func main() {

	config, err := config.Get()
	if err != nil {
		log.Errorf("unable to obtain configuration: %v:", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pd := sync.NewPagerDutyClient(pagerduty.NewClient(config.PagerDutyToken))
	pdGroups, err := pd.GetGroups(context.TODO())
	if err != nil {
		log.Panicf("Unable to instantiate PagerDuty client: %v:", err)
	}

	log.Infof("%+v\n", pdGroups)

	sl := slack.NewSlackClient(config.SlackToken)

	for key, pdUsers := range pdGroups {
		log.Infof("PD Group ID: %s", key)
		slackGroupName := key
		slackGroupID, groupExists, err := sl.GetGroupIDbyName(ctx, slackGroupName)
		if err != nil {
			log.Errorf("Error: %s", err)
		}

		if groupExists {
			log.Infof("Pre existing GroupID: [%s]", slackGroupID)
		}

		if !groupExists {
			slackGroupID, err = sl.CreateGroup(ctx, slackGroupName)
			if err != nil {
				log.Errorf("unable to create user group: [%s], error:[%s]", slackGroupName, err)
				continue
			}
			log.Infof("Created new group. Name:[%s], ID:[%s]", slackGroupName, slackGroupID)
		}

		emails := []string{}
		for _, user := range pdUsers {
			emails = append(emails, user.Email)
		}

		// Empty schedule
		if len(emails) == 0 {
			log.Infof("No emails found for slack group %s (ID:%s)", slackGroupName, slackGroupID)
			os.Exit(0)
		}

		err = sl.AddMembersToGroup(ctx, slackGroupID, emails...)
		if err != nil {
			log.Errorf("Error: %s", err)
			continue
		}
	}
}
