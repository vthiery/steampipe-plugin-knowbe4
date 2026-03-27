package knowbe4

import (
	"context"
	"errors"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// PhishingRecipient represents a recipient in a KnowBe4 Phishing Security Test.
type PhishingRecipient struct {
	RecipientID        int         `json:"recipient_id"`
	PstID              int         `json:"pst_id"`
	User               interface{} `json:"user"`
	Template           interface{} `json:"template"`
	ScheduledAt        string      `json:"scheduled_at"`
	DeliveredAt        *string     `json:"delivered_at"`
	OpenedAt           *string     `json:"opened_at"`
	ClickedAt          *string     `json:"clicked_at"`
	RepliedAt          *string     `json:"replied_at"`
	AttachmentOpenedAt *string     `json:"attachment_opened_at"`
	MacroEnabledAt     *string     `json:"macro_enabled_at"`
	DataEnteredAt      *string     `json:"data_entered_at"`
	QrCodeScanned      *string     `json:"qr_code_scanned"`
	ReportedAt         *string     `json:"reported_at"`
	BouncedAt          *string     `json:"bounced_at"`
	IP                 string      `json:"ip"`
	IPLocation         string      `json:"ip_location"`
	Browser            string      `json:"browser"`
	BrowserVersion     string      `json:"browser_version"`
	Os                 string      `json:"os"`
}

//// TABLE DEFINITION

func tableKnowBe4PhishingRecipient() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_phishing_recipient",
		Description: "Phishing security test recipient results in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listPhishingRecipients,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "pst_id", Require: plugin.Required},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "pst_id", Require: plugin.Required},
				{Name: "recipient_id", Require: plugin.Required},
			},
			Hydrate: getPhishingRecipient,
		},
		Columns: []*plugin.Column{
			{Name: "recipient_id", Type: proto.ColumnType_INT, Description: "Unique identifier for the recipient result."},
			{Name: "pst_id", Type: proto.ColumnType_INT, Description: "ID of the Phishing Security Test this result belongs to."},
			{Name: "ip", Type: proto.ColumnType_STRING, Description: "IP address from which the phishing link was clicked."},
			{Name: "ip_location", Type: proto.ColumnType_STRING, Description: "Geographic location of the IP address."},
			{Name: "browser", Type: proto.ColumnType_STRING, Description: "Browser used when clicking the phishing link."},
			{Name: "browser_version", Type: proto.ColumnType_STRING, Description: "Version of the browser used."},
			{Name: "os", Type: proto.ColumnType_STRING, Description: "Operating system of the device that clicked the phishing link."},
			{Name: "scheduled_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing email was scheduled."},
			{Name: "delivered_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing email was delivered."},
			{Name: "opened_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing email was opened."},
			{Name: "clicked_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing link was clicked."},
			{Name: "replied_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing email was replied to."},
			{Name: "attachment_opened_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing attachment was opened."},
			{Name: "macro_enabled_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time a macro was enabled from the attachment."},
			{Name: "data_entered_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time data was entered on the phishing page."},
			{Name: "qr_code_scanned", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the QR code was scanned."},
			{Name: "reported_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing email was reported."},
			{Name: "bounced_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the phishing email bounced."},
			{Name: "user", Type: proto.ColumnType_JSON, Description: "User details for the recipient."},
			{Name: "template", Type: proto.ColumnType_JSON, Description: "Phishing template details used for this recipient."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listPhishingRecipients(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	pstID := d.EqualsQuals["pst_id"].GetInt64Value()
	if pstID == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"cursor":   "true",
		"per_page": "500",
	}
	for {
		var recipients []PhishingRecipient
		nextCursor, err := client.get(ctx, fmt.Sprintf("/v1/phishing/security_tests/%d/recipients", pstID), params, &recipients)
		if err != nil {
			return nil, fmt.Errorf("listing phishing recipients for PST %d: %w", pstID, err)
		}

		for _, r := range recipients {
			d.StreamListItem(ctx, r)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if nextCursor == "" || len(recipients) == 0 {
			break
		}
		params["cursor"] = nextCursor
	}
	return nil, nil
}

func getPhishingRecipient(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	pstID := d.EqualsQuals["pst_id"].GetInt64Value()
	recipientID := d.EqualsQuals["recipient_id"].GetInt64Value()
	if pstID == 0 || recipientID == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var recipient PhishingRecipient
	if _, err := client.get(ctx, fmt.Sprintf("/v1/phishing/security_tests/%d/recipients/%d", pstID, recipientID), nil, &recipient); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting phishing recipient %d (PST %d): %w", recipientID, pstID, err)
	}
	return recipient, nil
}
