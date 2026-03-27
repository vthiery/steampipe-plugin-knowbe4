package knowbe4

import (
	"context"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// TrainingPolicy represents an uploaded policy in KnowBe4.
type TrainingPolicy struct {
	PolicyID        int    `json:"policy_id"`
	ContentType     string `json:"content_type"`
	Name            string `json:"name"`
	MinimumTime     int    `json:"minimum_time"`
	DefaultLanguage string `json:"default_language"`
	Status          string `json:"status"`
}

//// TABLE DEFINITION

func tableKnowBe4TrainingPolicy() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_training_policy",
		Description: "Uploaded compliance policies in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listTrainingPolicies,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("policy_id"),
			Hydrate:    getTrainingPolicy,
		},
		Columns: []*plugin.Column{
			{Name: "policy_id", Type: proto.ColumnType_INT, Description: "Unique identifier for the policy."},
			{Name: "content_type", Type: proto.ColumnType_STRING, Description: "Content type of the policy (e.g. Uploaded Policy)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the policy."},
			{Name: "minimum_time", Type: proto.ColumnType_INT, Description: "Minimum time (in minutes) required to review the policy."},
			{Name: "default_language", Type: proto.ColumnType_STRING, Description: "Default language of the policy."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the policy (e.g. published)."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listTrainingPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"cursor":   "true",
		"per_page": "500",
	}
	for {
		var policies []TrainingPolicy
		nextCursor, err := client.get(ctx, "/v1/training/policies", params, &policies)
		if err != nil {
			return nil, fmt.Errorf("listing training policies: %w", err)
		}

		for _, p := range policies {
			d.StreamListItem(ctx, p)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if nextCursor == "" || len(policies) == 0 {
			break
		}
		params["cursor"] = nextCursor
	}
	return nil, nil
}

func getTrainingPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["policy_id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var policy TrainingPolicy
	if _, err := client.get(ctx, fmt.Sprintf("/v1/training/policies/%d", id), nil, &policy); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting training policy %d: %w", id, err)
	}
	return policy, nil
}
