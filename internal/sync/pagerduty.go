package sync

import (
	"context"
	"regexp"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type pagerDutyClient struct {
	client *pagerduty.Client
}

//limitSchedules sets the maximum number of schedules to fetch in a single query (max = 100).
const limitSchedules uint = 100

//groupExp is a pattern that matches #slack:<group-name>.
var groupExp *regexp.Regexp

func init() {
	groupExp = regexp.MustCompile("#slack:([a-zA-Z0-9_-]+)")
}

func NewPagerDutyClient(client *pagerduty.Client) *pagerDutyClient {
	return &pagerDutyClient{client: client}
}

type User struct {
	Email    string
	Rotation string
}

type Schedule struct {
	ID    string
	Group string
	Name  string
}

//GetGroups returns a mapping of Slack group names to User objects that are in that group at this time.
func (c *pagerDutyClient) GetGroups(ctx context.Context) (map[string][]User, error) {
	now := time.Now().UTC()
	// ListOnCallUsersWithContext requires a time-slice of at least 1 second, so we'll run with that
	// If you don't pass any time-slice, you basically get everyone who is eligible to be on-call at any time
	opts := pagerduty.ListOnCallUsersOptions{
		Since: now.Format(time.RFC3339),
		Until: now.Add(time.Second).Format(time.RFC3339),
	}

	schedules, err := c.FindSchedules(ctx)
	if err != nil {
		return nil, err
	}

	var groups = make(map[string][]User, 0)

	for _, s := range schedules {
		users, err := c.client.ListOnCallUsersWithContext(ctx, s.ID, opts)
		if err != nil {
			return nil, err
		}

		if groups[s.Group] == nil {
			groups[s.Group] = make([]User, 0)
		}

		for _, u := range users {
			groups[s.Group] = append(groups[s.Group], User{
				Email:    u.Email,
				Rotation: s.Name,
			})
		}
	}

	return groups, nil
}

//FindSchedules finds all PagerDuty schedules that have a #slack:<group-name> tag in their description.
func (c *pagerDutyClient) FindSchedules(ctx context.Context) ([]Schedule, error) {
	opts := pagerduty.ListSchedulesOptions{
		Limit: limitSchedules,
	}

	var matched = make([]Schedule, 0)

	// ListSchedulesWithContext is paginated, and returns at most 100 Schedules per call
	more := true
	for more {
		schedules, err := c.client.ListSchedulesWithContext(ctx, opts)
		if err != nil {
			return nil, err
		}

		for _, s := range schedules.Schedules {
			// Grab the schedule if it contains #slack:<group-name> in the Description
			if res := groupExp.FindAllStringSubmatch(s.Description, -1); len(res) > 0 {
				group := res[0][1]

				matched = append(matched, Schedule{
					ID:    s.ID,
					Group: group,
					Name:  s.Name,
				})
			}
		}

		// Update the offset and possibly query for the next page
		more = schedules.More
		opts.Offset += opts.Limit
	}

	return matched, nil
}
