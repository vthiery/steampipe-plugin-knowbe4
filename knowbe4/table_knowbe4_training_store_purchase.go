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

// TrainingStorePurchase represents a purchased training module from the KnowBe4 store.
type TrainingStorePurchase struct {
	StorePurchasedID int     `json:"store_purchase_id"`
	ContentType      string  `json:"content_type"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Type             string  `json:"type"`
	Duration         int     `json:"duration"`
	Retired          bool    `json:"retired"`
	RetirementDate   *string `json:"retirement_date"`
	PublishDate      string  `json:"publish_date"`
	Publisher        string  `json:"publisher"`
	PurchaseDate     string  `json:"purchase_date"`
	PolicyURL        *string `json:"policy_url"`
}

//// TABLE DEFINITION

func tableKnowBe4TrainingStorePurchase() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_training_store_purchase",
		Description: "Training modules purchased from the KnowBe4 content store.",
		List: &plugin.ListConfig{
			Hydrate: listTrainingStorePurchases,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("store_purchased_id"),
			Hydrate:    getTrainingStorePurchase,
		},
		Columns: []*plugin.Column{
			{Name: "store_purchased_id", Type: proto.ColumnType_INT, Description: "Unique identifier for the store purchase."},
			{Name: "content_type", Type: proto.ColumnType_STRING, Description: "Content type of the purchase (e.g. Store Purchase)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the training module."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the training module."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the training content (e.g. Training Module, Video)."},
			{Name: "duration", Type: proto.ColumnType_INT, Description: "Duration of the training content in minutes."},
			{Name: "retired", Type: proto.ColumnType_BOOL, Description: "Whether the training content has been retired."},
			{Name: "retirement_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date on which the content was or will be retired."},
			{Name: "publish_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the content was published."},
			{Name: "publisher", Type: proto.ColumnType_STRING, Description: "Publisher of the training content."},
			{Name: "purchase_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the content was purchased."},
			{Name: "policy_url", Type: proto.ColumnType_STRING, Description: "URL of the policy associated with this content, if applicable."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listTrainingStorePurchases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"cursor":   "true",
		"per_page": "500",
	}
	for {
		var purchases []TrainingStorePurchase
		nextCursor, err := client.get(ctx, "/v1/training/store_purchases", params, &purchases)
		if err != nil {
			return nil, fmt.Errorf("listing training store purchases: %w", err)
		}

		for _, p := range purchases {
			d.StreamListItem(ctx, p)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if nextCursor == "" || len(purchases) == 0 {
			break
		}
		params["cursor"] = nextCursor
	}
	return nil, nil
}

func getTrainingStorePurchase(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["store_purchased_id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var purchase TrainingStorePurchase
	if _, err := client.get(ctx, fmt.Sprintf("/v1/training/store_purchases/%d", id), nil, &purchase); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting training store purchase %d: %w", id, err)
	}
	return purchase, nil
}
