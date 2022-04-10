package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"pd2slack/internal/config"
	"pd2slack/internal/slack"
	"pd2slack/internal/sync"

	"github.com/PagerDuty/go-pagerduty"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func initLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Cannot initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar = logger.Sugar()

}

func main() {
	initLogger()

	config, err := config.Get()
	if err != nil {
		sugar.Errorf("unable to obtain configuration: %v:", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pd := sync.NewPagerDutyClient(pagerduty.NewClient(config.PagerDutyToken))
	pdGroups, err := pd.GetGroups(context.TODO())
	if err != nil {
		msg := fmt.Sprintf("unable to instantiate PagerDuty client: %v:", err)
		sugar.Error(msg)
		os.Exit(0)
	}

	sugar.Infof("%+v\n", pdGroups)

	sl := slack.NewSlackClient(config.SlackToken)

	for key, pdUsers := range pdGroups {
		sugar.Info(fmt.Sprintf("PD Group ID: %s", key))
		slackGroupName := key
		slackGroupID, groupExists, err := sl.GetGroupIDbyName(ctx, slackGroupName)
		if err != nil {
			sugar.Errorf("Error: %s", err)
		}

		if groupExists {
			sugar.Infof("Pre existing GroupID: [%s]", slackGroupID)
		}

		if !groupExists {
			slackGroupID, err = sl.CreateGroup(ctx, slackGroupName)
			if err != nil {
				sugar.Errorf("unable to create user group: [%s], error:[%s]", slackGroupName, err)
				continue
			}
			sugar.Infof("Created new group. Name:[%s], ID:[%s]", slackGroupName, slackGroupID)
		}

		emails := []string{}
		for _, user := range pdUsers {
			emails = append(emails, user.Email)
		}

		// Empty schedule
		if len(emails) == 0 {
			sugar.Infof("No emails found for slack group %s (ID:%s)", slackGroupName, slackGroupID)
			os.Exit(0)
		}

		err = sl.AddMembersToGroup(ctx, slackGroupID, emails...)
		if err != nil {
			sugar.Errorf("Error: %s", err)
			continue
		}
	}
}
