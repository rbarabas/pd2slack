package slack

import (
	"context"
	"log"
	"strings"

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

func (s *SlackClient) CreateGroup(ctx context.Context, name string) (string, error) {
	ug := slack.UserGroup{
		Name:   name,
		Handle: name,
	}
	group, err := s.client.CreateUserGroupContext(ctx, ug)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to get the group by name:[%s]", name)
	}
	return group.ID, nil

}

func (s *SlackClient) AddMembersToGroup(ctx context.Context, groupID string, emailAddresses ...string) error {
	userIDs := []string{}
	for _, email := range emailAddresses {
		user, err := s.client.GetUserByEmailContext(ctx, email)
		if err != nil {
			return errors.Wrapf(err, "Unable to fetch member info for email:[%s]", email)
		}
		log.Printf("User info: [%+v]", user)
		userIDs = append(userIDs, user.ID)
	}

	joinedUserIDs := strings.Join(userIDs, ",")
	_, err := s.client.UpdateUserGroupMembersContext(ctx, groupID, joinedUserIDs)
	if err != nil {
		return errors.Wrapf(err, "Unable to add members:[%s] to user group:[%s]", emailAddresses, groupID)
	}

	return nil
}
