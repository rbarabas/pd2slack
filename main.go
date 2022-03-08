package main

import (
	"context"
	"fmt"
	"os"
	"pd2slack/internal/sync"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("pd2slack.conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
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

	fmt.Printf("%v\n", pdGroups)

	sl := slack.New(slackToken)

	fmt.Println("Obtain user data")
	user, err := sl.GetUserInfo("U01V0A4SYBC")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
	}

	fmt.Println("Obtain existing groups")
	groups, err := sl.GetUserGroups()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(2)
	}

	fmt.Println("List existing groups")
	for _, group := range groups {
		fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
	}

	fmt.Println("Create test group")
	ug := slack.UserGroup{
		Name:  "test-oncall",
		Users: []string{"U01V0A4SYBC"},
	}

	group, err := sl.CreateUserGroup(ug)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(3)
	}

	fmt.Println("Reading back test group")
	fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
}
