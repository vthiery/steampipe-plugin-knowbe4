// Package knowbe4 provides a Steampipe plugin for querying KnowBe4 resources using SQL.
package knowbe4

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the definition of the KnowBe4 Steampipe plugin.
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-knowbe4",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"knowbe4_account":                    tableKnowBe4Account(),
			"knowbe4_account_risk_score_history": tableKnowBe4AccountRiskScoreHistory(),
			"knowbe4_user":                       tableKnowBe4User(),
			"knowbe4_user_risk_score_history":    tableKnowBe4UserRiskScoreHistory(),
			"knowbe4_group":                      tableKnowBe4Group(),
			"knowbe4_group_risk_score_history":   tableKnowBe4GroupRiskScoreHistory(),
			"knowbe4_phishing_campaign":          tableKnowBe4PhishingCampaign(),
			"knowbe4_phishing_security_test":     tableKnowBe4PhishingSecurityTest(),
			"knowbe4_phishing_recipient":         tableKnowBe4PhishingRecipient(),
			"knowbe4_training_store_purchase":    tableKnowBe4TrainingStorePurchase(),
			"knowbe4_training_policy":            tableKnowBe4TrainingPolicy(),
			"knowbe4_training_campaign":          tableKnowBe4TrainingCampaign(),
			"knowbe4_training_enrollment":        tableKnowBe4TrainingEnrollment(),
		},
	}
	return p
}
