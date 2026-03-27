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

// Group represents a KnowBe4 user group.
type Group struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	GroupType        string  `json:"group_type"`
	ProvisioningGUID string  `json:"provisioning_guid"`
	MemberCount      int     `json:"member_count"`
	CurrentRiskScore float64 `json:"current_risk_score"`
	Status           string  `json:"status"`
}

//// TABLE DEFINITION

func tableKnowBe4Group() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_group",
		Description: "User groups defined in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listGroups,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getGroup,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier for the group."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the group."},
			{Name: "group_type", Type: proto.ColumnType_STRING, Description: "Type of the group (e.g. console_group)."},
			{Name: "provisioning_guid", Type: proto.ColumnType_STRING, Description: "GUID used for provisioning the group."},
			{Name: "member_count", Type: proto.ColumnType_INT, Description: "Number of members in the group."},
			{Name: "current_risk_score", Type: proto.ColumnType_DOUBLE, Description: "Current risk score for the group."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the group (active or archived)."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}
	if v := d.EqualsQualString("status"); v != "" {
		params["status"] = v
	}

	params["cursor"] = "true"
	params["per_page"] = "500"
	for {
		var groups []Group
		nextCursor, err := client.get(ctx, "/v1/groups", params, &groups)
		if err != nil {
			return nil, fmt.Errorf("listing groups: %w", err)
		}

		for _, g := range groups {
			d.StreamListItem(ctx, g)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if nextCursor == "" || len(groups) == 0 {
			break
		}
		params["cursor"] = nextCursor
	}
	return nil, nil
}

func getGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var group Group
	if _, err := client.get(ctx, fmt.Sprintf("/v1/groups/%d", id), nil, &group); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting group %d: %w", id, err)
	}
	return group, nil
}
