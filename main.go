package main

import (
	"fmt"
	"os"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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

	var opts pagerduty.ListEscalationPoliciesOptions
	pd := pagerduty.NewClient(pdToken)
	eps, err := pd.ListEscalationPolicies(opts)
	if err != nil {
		panic(err)
	}

	for _, p := range eps.EscalationPolicies {
		m := make(map[interface{}]interface{})
		if err := yaml.Unmarshal([]byte(p.Description), &m); err != nil {
			continue
		}
		if m["slackgroup"] != nil {
			fmt.Printf("%s (%s)\n", m["slackgroup"], p.Name)
		}
	}

	sl := slack.New(slackToken)

	fmt.Println("Obtain user data")
	user, err := sl.GetUserInfo("Robert Barabas")
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
		ID:    "test-oncall",
		Users: []string{"Robert Barabas", "Jake Edgington", "Chetan Gowda"},
	}

	sl.CreateUserGroup(ug)
}
