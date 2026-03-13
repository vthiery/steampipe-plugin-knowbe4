package knowbe4

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// AccountAdmin represents an admin user in the account.
type AccountAdmin struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Account represents the KnowBe4 account.
type Account struct {
	Name                string         `json:"name"`
	Type                string         `json:"type"`
	Domains             []string       `json:"domains"`
	Admins              []AccountAdmin `json:"admins"`
	SubscriptionLevel   string         `json:"subscription_level"`
	SubscriptionEndDate string         `json:"subscription_end_date"`
	NumberOfSeats       int            `json:"number_of_seats"`
	CurrentRiskScore    float64        `json:"current_risk_score"`
}

//// TABLE DEFINITION

func tableKnowBe4Account() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_account",
		Description: "KnowBe4 account details including subscription level and current risk score.",
		List: &plugin.ListConfig{
			Hydrate: listAccount,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the KnowBe4 account."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the account (e.g. paid)."},
			{Name: "subscription_level", Type: proto.ColumnType_STRING, Description: "Subscription level (e.g. Diamond, Platinum)."},
			{Name: "subscription_end_date", Type: proto.ColumnType_STRING, Description: "Date on which the subscription ends."},
			{Name: "number_of_seats", Type: proto.ColumnType_INT, Description: "Number of licensed seats in the account."},
			{Name: "current_risk_score", Type: proto.ColumnType_DOUBLE, Description: "Current account-level risk score."},
			{Name: "domains", Type: proto.ColumnType_JSON, Description: "List of domains associated with the account."},
			{Name: "admins", Type: proto.ColumnType_JSON, Description: "List of admin users for the account."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := client.get(ctx, "/v1/account", nil, &account); err != nil {
		return nil, fmt.Errorf("getting account: %w", err)
	}

	d.StreamListItem(ctx, account)
	return nil, nil
}
