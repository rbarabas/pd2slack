package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	SlackToken     string `mapstructure:"SLACK_TOKEN"`
	PagerDutyToken string `mapstructure:"PAGERDUTY_TOKEN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("pd2slack")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}

func Get() (config Config, err error) {
	c, err := LoadConfig(".")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Override if ENV variables
	if pdEnv := os.Getenv("PAGERDUTY_TOKEN"); pdEnv != "" {
		c.PagerDutyToken = pdEnv
	}

	if slEnv := os.Getenv("SLACK_TOKEN"); slEnv != "" {
		c.SlackToken = slEnv
	}

	if c.PagerDutyToken == "" {
		err = fmt.Errorf("no PagerDuty authorization token found in configuration")
		return
	}

	if c.SlackToken == "" {
		err = fmt.Errorf("no Slack authorization token found in configuration")
		return
	}

	return c, err
}
