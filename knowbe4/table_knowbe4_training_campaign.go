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

// TrainingCampaign represents a training campaign in KnowBe4.
type TrainingCampaign struct {
	CampaignID               int             `json:"campaign_id"`
	Name                     string          `json:"name"`
	Groups                   json.RawMessage `json:"groups"`
	Status                   string          `json:"status"`
	Content                  json.RawMessage `json:"content"`
	DurationType             string          `json:"duration_type"`
	StartDate                string          `json:"start_date"`
	EndDate                  string          `json:"end_date"`
	RelativeDuration         string          `json:"relative_duration"`
	AutoEnroll               bool            `json:"auto_enroll"`
	AllowMultipleEnrollments bool            `json:"allow_multiple_enrollments"`
	CompletionPercentage     float64         `json:"completion_percentage"`
}

//// TABLE DEFINITION

func tableKnowBe4TrainingCampaign() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_training_campaign",
		Description: "Training campaigns configured in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listTrainingCampaigns,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("campaign_id"),
			Hydrate:    getTrainingCampaign,
		},
		Columns: []*plugin.Column{
			{Name: "campaign_id", Type: proto.ColumnType_INT, Description: "Unique identifier for the training campaign."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the training campaign."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the campaign (e.g. Active, Completed)."},
			{Name: "duration_type", Type: proto.ColumnType_STRING, Description: "Type of duration (e.g. Specific End Date, Relative Duration)."},
			{Name: "start_date", Type: proto.ColumnType_TIMESTAMP, Description: "Start date of the training campaign."},
			{Name: "end_date", Type: proto.ColumnType_TIMESTAMP, Description: "End date of the training campaign."},
			{Name: "relative_duration", Type: proto.ColumnType_STRING, Description: "Relative duration of the campaign, if applicable."},
			{Name: "auto_enroll", Type: proto.ColumnType_BOOL, Description: "Whether new users are automatically enrolled."},
			{Name: "allow_multiple_enrollments", Type: proto.ColumnType_BOOL, Description: "Whether users can be enrolled multiple times."},
			{Name: "completion_percentage", Type: proto.ColumnType_DOUBLE, Description: "Percentage of enrolled users who have completed the campaign."},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups enrolled in this training campaign."},
			{Name: "content", Type: proto.ColumnType_JSON, Description: "Training content assigned to this campaign."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listTrainingCampaigns(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"cursor":   "true",
		"per_page": "500",
	}
	for {
		var campaigns []TrainingCampaign
		nextCursor, err := client.get(ctx, "/v1/training/campaigns", params, &campaigns)
		if err != nil {
			return nil, fmt.Errorf("listing training campaigns: %w", err)
		}

		for _, c := range campaigns {
			d.StreamListItem(ctx, c)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if nextCursor == "" || len(campaigns) == 0 {
			break
		}
		params["cursor"] = nextCursor
	}
	return nil, nil
}

func getTrainingCampaign(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["campaign_id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var campaign TrainingCampaign
	if _, err := client.get(ctx, fmt.Sprintf("/v1/training/campaigns/%d", id), nil, &campaign); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting training campaign %d: %w", id, err)
	}
	return campaign, nil
}
