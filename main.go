package main

import (
	"context"
	"log"
	"os"
	"pd2slack/internal/slack"
	"pd2slack/internal/sync"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("pd2slack.conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	log.SetFlags(log.Lshortfile)
	log.SetPrefix("pd2slack: ")

	ctx := context.Background()

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	pdToken := viper.GetString("pagerduty.authtoken")
	if pdToken == "" {
		panic("no PagerDuty authorization token found in configuration")
	}

	slackToken := viper.GetString("slack.authtoken")
	if slackToken == "" {
		panic("no Slack authorization token found in configuration")
	}

	pd := sync.NewPagerDutyClient(pagerduty.NewClient(pdToken))
	pdGroups, err := pd.GetGroups(context.TODO())
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", pdGroups)

	sl := slack.NewSlackClient(slackToken)

	for key, _ := range pdGroups {
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

	}
}
