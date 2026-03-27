package knowbe4

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// UserRiskScoreEntry represents a risk score entry for a specific user.
type UserRiskScoreEntry struct {
	UserID    int     `json:"user_id"`
	RiskScore float64 `json:"risk_score"`
	Date      string  `json:"date"`
}

//// TABLE DEFINITION

func tableKnowBe4UserRiskScoreHistory() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_user_risk_score_history",
		Description: "Historical risk score records for each user in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listUserRiskScoreHistory,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_id", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			{Name: "user_id", Type: proto.ColumnType_INT, Description: "ID of the user this risk score history belongs to."},
			{Name: "risk_score", Type: proto.ColumnType_DOUBLE, Description: "Risk score at the recorded date."},
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the risk score was recorded."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listUserRiskScoreHistory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	userID := d.EqualsQuals["user_id"].GetInt64Value()
	if userID == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var entries []RiskScoreEntry
	if _, err := client.get(ctx, fmt.Sprintf("/v1/users/%d/risk_score_history", userID), map[string]string{"full": "true"}, &entries); err != nil {
		return nil, fmt.Errorf("getting user %d risk score history: %w", userID, err)
	}

	for _, e := range entries {
		d.StreamListItem(ctx, UserRiskScoreEntry{
			UserID:    int(userID),
			RiskScore: e.RiskScore,
			Date:      e.Date,
		})
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
