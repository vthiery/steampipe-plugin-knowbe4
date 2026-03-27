package knowbe4

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// GroupRiskScoreEntry represents a risk score entry for a specific group.
type GroupRiskScoreEntry struct {
	GroupID   int     `json:"group_id"`
	RiskScore float64 `json:"risk_score"`
	Date      string  `json:"date"`
}

//// TABLE DEFINITION

func tableKnowBe4GroupRiskScoreHistory() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_group_risk_score_history",
		Description: "Historical risk score records for each group in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listGroupRiskScoreHistory,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "group_id", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			{Name: "group_id", Type: proto.ColumnType_INT, Description: "ID of the group this risk score history belongs to."},
			{Name: "risk_score", Type: proto.ColumnType_DOUBLE, Description: "Risk score at the recorded date."},
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the risk score was recorded."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listGroupRiskScoreHistory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	groupID := d.EqualsQuals["group_id"].GetInt64Value()
	if groupID == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var entries []RiskScoreEntry
	if _, err := client.get(ctx, fmt.Sprintf("/v1/groups/%d/risk_score_history", groupID), map[string]string{"full": "true"}, &entries); err != nil {
		return nil, fmt.Errorf("getting group %d risk score history: %w", groupID, err)
	}

	for _, e := range entries {
		d.StreamListItem(ctx, GroupRiskScoreEntry{
			GroupID:   int(groupID),
			RiskScore: e.RiskScore,
			Date:      e.Date,
		})
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
