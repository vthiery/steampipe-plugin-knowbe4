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

// EnrollmentUser represents the user embedded in a training enrollment.
type EnrollmentUser struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	EmployeeNumber int    `json:"employee_number"`
}

// TrainingEnrollment represents a training enrollment record in KnowBe4.
type TrainingEnrollment struct {
	EnrollmentID       int            `json:"enrollment_id"`
	ContentType        string         `json:"content_type"`
	ModuleName         string         `json:"module_name"`
	User               EnrollmentUser `json:"user"`
	CampaignID         int            `json:"campaign_id"`
	CampaignName       string         `json:"campaign_name"`
	EnrollmentDate     string         `json:"enrollment_date"`
	StartDate          string         `json:"start_date"`
	CompletionDate     string         `json:"completion_date"`
	Status             string         `json:"status"`
	TimeSpent          int            `json:"time_spent"`
	PolicyAcknowledged bool           `json:"policy_acknowledged"`
	StorePurchaseID    int            `json:"store_purchase_id"`
}

//// TABLE DEFINITION

func tableKnowBe4TrainingEnrollment() *plugin.Table {
	return &plugin.Table{
		Name:        "knowbe4_training_enrollment",
		Description: "Training enrollment records for users in the KnowBe4 account.",
		List: &plugin.ListConfig{
			Hydrate: listTrainingEnrollments,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "campaign_id", Require: plugin.Optional},
				{Name: "user_id", Require: plugin.Optional},
				{Name: "store_purchase_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("enrollment_id"),
			Hydrate:    getTrainingEnrollment,
		},
		Columns: []*plugin.Column{
			{Name: "enrollment_id", Type: proto.ColumnType_INT, Description: "Unique identifier for the training enrollment."},
			{Name: "content_type", Type: proto.ColumnType_STRING, Description: "Type of content (e.g. Training Module, Uploaded Policy)."},
			{Name: "module_name", Type: proto.ColumnType_STRING, Description: "Name of the training module or policy."},
			{Name: "campaign_name", Type: proto.ColumnType_STRING, Description: "Name of the training campaign this enrollment belongs to."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Enrollment status (e.g. Passed, In Progress, Not Started)."},
			{Name: "time_spent", Type: proto.ColumnType_INT, Description: "Time spent on the training in seconds."},
			{Name: "policy_acknowledged", Type: proto.ColumnType_BOOL, Description: "Whether the user acknowledged the policy."},
			{Name: "enrollment_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the user was enrolled."},
			{Name: "start_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the user started the training."},
			{Name: "completion_date", Type: proto.ColumnType_TIMESTAMP, Description: "Date the user completed the training."},
			{Name: "campaign_id", Type: proto.ColumnType_INT, Description: "ID of the training campaign this enrollment belongs to."},
			{Name: "user_id", Type: proto.ColumnType_INT, Transform: transform.FromField("User.ID"), Description: "ID of the enrolled user."},
			{Name: "store_purchase_id", Type: proto.ColumnType_INT, Description: "Store purchase (module) ID for this enrollment."},
			// Nested user object
			{Name: "user", Type: proto.ColumnType_JSON, Description: "User details for the enrolled user."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("ModuleName"), Description: "The display name of the resource."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listTrainingEnrollments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	baseParams := map[string]string{
		"include_campaign_id":       "true",
		"include_store_purchase_id": "true",
		"include_employee_number":   "true",
	}
	if v := d.EqualsQuals["campaign_id"]; v != nil {
		if id := v.GetInt64Value(); id != 0 {
			baseParams["campaign_id"] = fmt.Sprintf("%d", id)
		}
	}
	if v := d.EqualsQuals["user_id"]; v != nil {
		if id := v.GetInt64Value(); id != 0 {
			baseParams["user_id"] = fmt.Sprintf("%d", id)
		}
	}
	if v := d.EqualsQuals["store_purchase_id"]; v != nil {
		if id := v.GetInt64Value(); id != 0 {
			baseParams["store_purchase_id"] = fmt.Sprintf("%d", id)
		}
	}

	for page := 1; ; page++ {
		params := make(map[string]string, len(baseParams)+2)
		for k, v := range baseParams {
			params[k] = v
		}
		params["page"] = fmt.Sprintf("%d", page)
		params["per_page"] = "500"

		var enrollments []TrainingEnrollment
		if err := client.get(ctx, "/v1/training/enrollments", params, &enrollments); err != nil {
			return nil, fmt.Errorf("listing training enrollments: %w", err)
		}

		for _, e := range enrollments {
			d.StreamListItem(ctx, e)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if len(enrollments) == 0 {
			break
		}
	}
	return nil, nil
}

func getTrainingEnrollment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["enrollment_id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var enrollment TrainingEnrollment
	if err := client.get(ctx, fmt.Sprintf("/v1/training/enrollments/%d", id), nil, &enrollment); err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting training enrollment %d: %w", id, err)
	}
	return enrollment, nil
}
