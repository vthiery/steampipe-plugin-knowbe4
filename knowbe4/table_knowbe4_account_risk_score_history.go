package knowbe4

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// RiskScoreEntry represents a single risk score data point.
type RiskScoreEntry struct {
	RiskScore float64 `json:"risk_score"`
	Date      string  `json:"date"`
}

//// TABLE DEFINITION

func tableKnowBe4AccountRiskScoreHistory() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_account_risk_score_history",
		Description: "Historical risk score records for the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listAccountRiskScoreHistory,
		},
		Columns: []*plugin.Column{
			{Name: "risk_score", Type: proto.ColumnType_DOUBLE, Description: "Risk score at the recorded date."},
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the risk score was recorded."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAccountRiskScoreHistory(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var entries []RiskScoreEntry
	if err := client.get(ctx, "/v1/account/risk_score_history", map[string]string{"full": "true"}, &entries); err != nil {
		return nil, fmt.Errorf("getting account risk score history: %w", err)
	}

	for _, entry := range entries {
		d.StreamListItem(ctx, entry)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
