package knowbe4

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// PhishingCampaign represents a KnowBe4 phishing campaign.
type PhishingCampaign struct {
	CampaignID               int             `json:"campaign_id"`
	Name                     string          `json:"name"`
	Groups                   json.RawMessage `json:"groups"`
	LastPhishPronePercentage float64         `json:"last_phish_prone_percentage"`
	LastRun                  string          `json:"last_run"`
	Status                   string          `json:"status"`
	Hidden                   bool            `json:"hidden"`
	SendDuration             string          `json:"send_duration"`
	TrackDuration            string          `json:"track_duration"`
	Frequency                string          `json:"frequency"`
	DifficultyFilter         json.RawMessage `json:"difficulty_filter"`
	CreateDate               string          `json:"create_date"`
	PstsCount                int             `json:"psts_count"`
	Psts                     json.RawMessage `json:"psts"`
}

//// TABLE DEFINITION

func tableKnowBe4PhishingCampaign() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_phishing_campaign",
		Description: "Phishing campaigns configured in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listPhishingCampaigns,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("campaign_id"),
			Hydrate:    getPhishingCampaign,
		},
		Columns: []*plugin.Column{
			{Name: "campaign_id", Type: proto.ColumnType_INT, Description: "Unique identifier for the phishing campaign."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the phishing campaign."},
			{Name: "last_phish_prone_percentage", Type: proto.ColumnType_DOUBLE, Description: "Most recent Phish-Prone Percentage for this campaign."},
			{Name: "last_run", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the campaign last ran."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the campaign (e.g. Closed, Active)."},
			{Name: "hidden", Type: proto.ColumnType_BOOL, Description: "Whether the campaign is hidden."},
			{Name: "send_duration", Type: proto.ColumnType_STRING, Description: "Duration over which phishing emails are sent."},
			{Name: "track_duration", Type: proto.ColumnType_STRING, Description: "Duration over which responses are tracked."},
			{Name: "frequency", Type: proto.ColumnType_STRING, Description: "Frequency of the campaign (e.g. One Time, Monthly)."},
			{Name: "create_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the campaign was created."},
			{Name: "psts_count", Type: proto.ColumnType_INT, Description: "Number of phishing security tests in this campaign."},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups targeted by the campaign."},
			{Name: "difficulty_filter", Type: proto.ColumnType_JSON, Description: "Difficulty filter settings for the campaign."},
			{Name: "psts", Type: proto.ColumnType_JSON, Description: "List of phishing security tests in this campaign."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listPhishingCampaigns(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}
	for page := 1; ; page++ {
		params["page"] = fmt.Sprintf("%d", page)
		params["per_page"] = "500"

		var campaigns []PhishingCampaign
		if err := client.get(ctx, "/v1/phishing/campaigns", params, &campaigns); err != nil {
			return nil, fmt.Errorf("listing phishing campaigns: %w", err)
		}

		for _, c := range campaigns {
			d.StreamListItem(ctx, c)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if len(campaigns) == 0 {
			break
		}
	}
	return nil, nil
}

func getPhishingCampaign(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["campaign_id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var campaign PhishingCampaign
	if err := client.get(ctx, fmt.Sprintf("/v1/phishing/campaigns/%d", id), nil, &campaign); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting phishing campaign %d: %w", id, err)
	}
	return campaign, nil
}
