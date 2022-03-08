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

	log.Printf("%v\n", pdGroups)

	sl := slack.NewSlackClient(slackToken)

	groupName := "test-oncall"
	grpID, isExists, err := sl.GetGroupIDbyName(ctx, "test-oncall")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	if isExists {
		log.Printf("GroupID: [%s]", grpID)
	} else {
		log.Printf("Group:[%s] doesn't exits", groupName)

	}

	// sl := slack.New(slackToken)

	// log.Println("Obtain user data")
	// user, err := sl.GetUserInfo("U01V0A4SYBC")
	// if err != nil {
	// 	log.Printf("%s\n", err)
	// } else {
	// 	log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
	// }

	// log.Println("Obtain existing groups")
	// groups, err := sl.GetUserGroups()
	// if err != nil {
	// 	log.Printf("%s\n", err)
	// 	os.Exit(2)
	// }

	// log.Println("List existing groups")
	// for _, group := range groups {
	// 	log.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	// }

	// log.Println("Create test group")
	// ug := slack.UserGroup{
	// 	Name:   "test-oncall",
	// 	Users:  []string{"U01V0A4SYBC"},
	// 	Handle: "test-oncall",
	// }

	// if group, err := sl.DisableUserGroup("test-oncall"); err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// } else {
	// 	log.Printf("\nDeleted user group: [%+v]", group)
	// }

	// group, err := sl.CreateUserGroup(ug)
	// if err != nil {
	// 	log.Printf("%s\n", err)
	// 	os.Exit(3)
	// }

	// log.Println("Reading back test group")
	// log.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
}
