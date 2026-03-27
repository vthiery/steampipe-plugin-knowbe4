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

// User represents a KnowBe4 user.
type User struct {
	ID                   int      `json:"id"`
	EmployeeNumber       string   `json:"employee_number"`
	FirstName            string   `json:"first_name"`
	LastName             string   `json:"last_name"`
	JobTitle             string   `json:"job_title"`
	Email                string   `json:"email"`
	PhishPronePercentage float64  `json:"phish_prone_percentage"`
	PhoneNumber          string   `json:"phone_number"`
	Extension            string   `json:"extension"`
	MobilePhoneNumber    string   `json:"mobile_phone_number"`
	Location             string   `json:"location"`
	Division             string   `json:"division"`
	ManagerName          string   `json:"manager_name"`
	ManagerEmail         string   `json:"manager_email"`
	ProvisioningManaged  bool     `json:"provisioning_managed"`
	ProvisioningGUID     *string  `json:"provisioning_guid"`
	Groups               []int    `json:"groups"`
	CurrentRiskScore     float64  `json:"current_risk_score"`
	Aliases              []string `json:"aliases"`
	JoinedOn             string   `json:"joined_on"`
	LastSignIn           string   `json:"last_sign_in"`
	Status               string   `json:"status"`
	Organization         string   `json:"organization"`
	Department           string   `json:"department"`
	Language             string   `json:"language"`
	Comment              string   `json:"comment"`
	EmployeeStartDate    string   `json:"employee_start_date"`
	ArchivedAt           *string  `json:"archived_at"`
	CustomField1         *string  `json:"custom_field_1"`
	CustomField2         *string  `json:"custom_field_2"`
	CustomField3         *string  `json:"custom_field_3"`
	CustomField4         *string  `json:"custom_field_4"`
	CustomDate1          *string  `json:"custom_date_1"`
	CustomDate2          *string  `json:"custom_date_2"`
}

//// TABLE DEFINITION

func tableKnowBe4User() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_user",
		Description: "Users registered in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier for the user."},
			{Name: "employee_number", Type: proto.ColumnType_STRING, Description: "Employee number of the user."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "First name of the user."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "Last name of the user."},
			{Name: "job_title", Type: proto.ColumnType_STRING, Description: "Job title of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the user."},
			{Name: "phish_prone_percentage", Type: proto.ColumnType_DOUBLE, Description: "Phish-Prone Percentage (PPP) for the user."},
			{Name: "phone_number", Type: proto.ColumnType_STRING, Description: "Phone number of the user."},
			{Name: "extension", Type: proto.ColumnType_STRING, Description: "Phone extension of the user."},
			{Name: "mobile_phone_number", Type: proto.ColumnType_STRING, Description: "Mobile phone number of the user."},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "Office location of the user."},
			{Name: "division", Type: proto.ColumnType_STRING, Description: "Division the user belongs to."},
			{Name: "manager_name", Type: proto.ColumnType_STRING, Description: "Name of the user's manager."},
			{Name: "manager_email", Type: proto.ColumnType_STRING, Description: "Email of the user's manager."},
			{Name: "provisioning_managed", Type: proto.ColumnType_BOOL, Description: "Whether the user is managed via provisioning."},
			{Name: "provisioning_guid", Type: proto.ColumnType_STRING, Description: "GUID used for provisioning the user."},
			{Name: "current_risk_score", Type: proto.ColumnType_DOUBLE, Description: "Current risk score for the user."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the user (active or archived)."},
			{Name: "organization", Type: proto.ColumnType_STRING, Description: "Organization the user belongs to."},
			{Name: "department", Type: proto.ColumnType_STRING, Description: "Department the user belongs to."},
			{Name: "language", Type: proto.ColumnType_STRING, Description: "Preferred language of the user."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Comment or notes about the user."},
			{Name: "joined_on", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the user joined."},
			{Name: "last_sign_in", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time of the user's last sign in."},
			{Name: "employee_start_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the user started as an employee."},
			{Name: "archived_at", Type: proto.ColumnType_TIMESTAMP, Description: "Date and time the user was archived, if applicable."},
			{Name: "custom_field_1", Type: proto.ColumnType_STRING, Description: "Custom field 1 value for the user."},
			{Name: "custom_field_2", Type: proto.ColumnType_STRING, Description: "Custom field 2 value for the user."},
			{Name: "custom_field_3", Type: proto.ColumnType_STRING, Description: "Custom field 3 value for the user."},
			{Name: "custom_field_4", Type: proto.ColumnType_STRING, Description: "Custom field 4 value for the user."},
			{Name: "custom_date_1", Type: proto.ColumnType_STRING, Description: "Custom date 1 value for the user."},
			{Name: "custom_date_2", Type: proto.ColumnType_STRING, Description: "Custom date 2 value for the user."},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "List of group IDs the user belongs to."},
			{Name: "aliases", Type: proto.ColumnType_JSON, Description: "Email aliases for the user."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Email"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}
	if v := d.EqualsQualString("status"); v != "" {
		params["status"] = v
	}

	path := "/v1/users"

	params["cursor"] = "true"
	params["per_page"] = "500"
	for {
		var pageUsers []User
		nextCursor, err := client.get(ctx, path, params, &pageUsers)
		if err != nil {
			return nil, fmt.Errorf("listing users: %w", err)
		}

		for _, u := range pageUsers {
			d.StreamListItem(ctx, u)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if nextCursor == "" || len(pageUsers) == 0 {
			break
		}
		params["cursor"] = nextCursor
	}
	return nil, nil
}

func getUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var user User
	if _, err := client.get(ctx, fmt.Sprintf("/v1/users/%d", id), nil, &user); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user %d: %w", id, err)
	}
	return user, nil
}
