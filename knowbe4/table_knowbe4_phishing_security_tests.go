package knowbe4

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// PhishingSecurityTest represents a KnowBe4 Phishing Security Test (PST).
type PhishingSecurityTest struct {
	CampaignID           int             `json:"campaign_id"`
	PstID                int             `json:"pst_id"`
	Status               string          `json:"status"`
	Name                 string          `json:"name"`
	Groups               json.RawMessage `json:"groups"`
	PhishPronePercentage float64         `json:"phish_prone_percentage"`
	StartedAt            string          `json:"started_at"`
	Duration             int             `json:"duration"`
	Categories           json.RawMessage `json:"categories"`
	Template             json.RawMessage `json:"template"`
	LandingPage          json.RawMessage `json:"landing-page"`
	ScheduledCount       int             `json:"scheduled_count"`
	DeliveredCount       int             `json:"delivered_count"`
	OpenedCount          int             `json:"opened_count"`
	ClickedCount         int             `json:"clicked_count"`
	RepliedCount         int             `json:"replied_count"`
	AttachmentOpenCount  int             `json:"attachment_open_count"`
	MacroEnabledCount    int             `json:"macro_enabled_count"`
	DataEnteredCount     int             `json:"data_entered_count"`
	QrCodeScannedCount   int             `json:"qr_code_scanned_count"`
	ReportedCount        int             `json:"reported_count"`
	BouncedCount         int             `json:"bounced_count"`
}

//// TABLE DEFINITION

func tableKnowBe4PhishingSecurityTest() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_phishing_security_test",
		Description: "Phishing security tests (PSTs) run within KnowBe4 phishing campaigns.",
		List: &plugin.ListConfig{
			Hydrate: listPhishingSecurityTests,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "campaign_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("pst_id"),
			Hydrate:    getPhishingSecurityTest,
		},
		Columns: []*plugin.Column{
			{Name: "pst_id", Type: proto.ColumnType_INT, Description: "Unique identifier for the phishing security test."},
			{Name: "campaign_id", Type: proto.ColumnType_INT, Description: "ID of the campaign this PST belongs to."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the phishing security test."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the PST (e.g. Closed, Active)."},
			{Name: "phish_prone_percentage", Type: proto.ColumnType_DOUBLE, Description: "Phish-Prone Percentage for this PST."},
			{Name: "started_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the PST was started."},
			{Name: "duration", Type: proto.ColumnType_INT, Description: "Duration of the PST in days."},
			{Name: "scheduled_count", Type: proto.ColumnType_INT, Description: "Number of emails scheduled to be sent."},
			{Name: "delivered_count", Type: proto.ColumnType_INT, Description: "Number of emails delivered."},
			{Name: "opened_count", Type: proto.ColumnType_INT, Description: "Number of emails opened."},
			{Name: "clicked_count", Type: proto.ColumnType_INT, Description: "Number of phishing links clicked."},
			{Name: "replied_count", Type: proto.ColumnType_INT, Description: "Number of emails replied to."},
			{Name: "attachment_open_count", Type: proto.ColumnType_INT, Description: "Number of attachments opened."},
			{Name: "macro_enabled_count", Type: proto.ColumnType_INT, Description: "Number of macros enabled from attachments."},
			{Name: "data_entered_count", Type: proto.ColumnType_INT, Description: "Number of users who entered data on the phishing page."},
			{Name: "qr_code_scanned_count", Type: proto.ColumnType_INT, Description: "Number of QR codes scanned."},
			{Name: "reported_count", Type: proto.ColumnType_INT, Description: "Number of phishing emails reported."},
			{Name: "bounced_count", Type: proto.ColumnType_INT, Description: "Number of emails that bounced."},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups that were targeted by this PST."},
			{Name: "categories", Type: proto.ColumnType_JSON, Description: "Template categories used in this PST."},
			{Name: "template", Type: proto.ColumnType_JSON, Description: "Phishing template used in this PST."},
			{Name: "landing_page", Type: proto.ColumnType_JSON, Description: "Landing page used in this PST."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listPhishingSecurityTests(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// If campaign_id qual is set, use the campaign-scoped endpoint.
	path := "/v1/phishing/security_tests"
	if v := d.EqualsQuals["campaign_id"]; v != nil {
		campaignID := v.GetInt64Value()
		if campaignID != 0 {
			path = fmt.Sprintf("/v1/phishing/campaigns/%d/security_tests", campaignID)
		}
	}

	for page := 1; ; page++ {
		params := map[string]string{
			"page":     fmt.Sprintf("%d", page),
			"per_page": "500",
		}

		var psts []PhishingSecurityTest
		if err := client.get(ctx, path, params, &psts); err != nil {
			return nil, fmt.Errorf("listing phishing security tests: %w", err)
		}

		for _, pst := range psts {
			d.StreamListItem(ctx, pst)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if len(psts) == 0 {
			break
		}
	}
	return nil, nil
}

func getPhishingSecurityTest(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["pst_id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var pst PhishingSecurityTest
	if err := client.get(ctx, fmt.Sprintf("/v1/phishing/security_tests/%d", id), nil, &pst); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting phishing security test %d: %w", id, err)
	}
	return pst, nil
}
