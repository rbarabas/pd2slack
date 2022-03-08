package slack

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

func NewSlackClient(token string) *SlackClient {
	return &SlackClient{
		client: slack.New(token),
	}
}

type SlackClient struct {
	client *slack.Client
}

func (s *SlackClient) GetGroupIDbyName(ctx context.Context, name string) (string, bool, error) {
	groups, err := s.client.GetUserGroupsContext(ctx, slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		return "", false, errors.Wrapf(err, "Unable to get the group by name:[%s]", name)
	}
	for _, grp := range groups {
		if grp.Name == name {
			return grp.ID, true, nil
		}
	}
	return "", false, nil
}

// func (s *SlackClient) ClearGroup(id string) {

// }

// func foo() {
// 	sl := slack.New(slackToken)

// 	log.Println("Obtain user data")
// 	user, err := sl.GetUserInfo("U01V0A4SYBC")
// 	if err != nil {
// 		log.Printf("%s\n", err)
// 	} else {
// 		log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
// 	}

// 	log.Println("Obtain existing groups")
// 	groups, err := sl.GetUserGroups()
// 	if err != nil {
// 		log.Printf("%s\n", err)
// 		os.Exit(2)
// 	}

// 	log.Println("List existing groups")
// 	for _, group := range groups {
// 		log.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
// 	}

// 	log.Println("Create test group")
// 	ug := slack.UserGroup{
// 		Name:   "test-oncall",
// 		Users:  []string{"U01V0A4SYBC"},
// 		Handle: "test-oncall",
// 	}

// 	if group, err := sl.DisableUserGroup("test-oncall"); err != nil {
// 		log.Println(err)
// 		os.Exit(1)
// 	} else {
// 		log.Printf("\nDeleted user group: [%+v]", group)
// 	}

// 	group, err := sl.CreateUserGroup(ug)
// 	if err != nil {
// 		log.Printf("%s\n", err)
// 		os.Exit(3)
// 	}

// 	log.Println("Reading back test group")
// 	log.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
// }
