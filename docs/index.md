# KnowBe4 Plugin for Steampipe

Use SQL to query users, phishing campaigns, training enrollments, and more from [KnowBe4](https://www.knowbe4.com).

## Installation

Clone and build the plugin:

```sh
git clone https://github.com/vthiery/steampipe-plugin-knowbe4.git
cd steampipe-plugin-knowbe4
mkdir -p ~/.steampipe/plugins/local/knowbe4
go build -o ~/.steampipe/plugins/local/knowbe4/knowbe4.plugin .
```

## Configuration

Copy the sample config:

```sh
cp config/knowbe4.spc ~/.steampipe/config/knowbe4.spc
```

Edit `~/.steampipe/config/knowbe4.spc`:

```hcl
connection "knowbe4" {
  plugin = "local/knowbe4"

  # KnowBe4 API token. Generate one at Account Settings → Account Info → API.
  api_token = "YOUR_API_TOKEN"

  # KnowBe4 API region (optional, defaults to "us").
  # Valid values: "us", "eu", "ca", "uk", "de"
  # api_region = "us"
}
```

## Tables

| Table | Description |
|-------|-------------|
| [knowbe4_account](tables/knowbe4_account.md) | KnowBe4 account details including subscription level and current risk score. |
| [knowbe4_account_risk_score_history](tables/knowbe4_account_risk_score_history.md) | Historical risk score records for the KnowBe4 account. |
| [knowbe4_group](tables/knowbe4_group.md) | User groups defined in the KnowBe4 account. |
| [knowbe4_group_risk_score_history](tables/knowbe4_group_risk_score_history.md) | Historical risk score records for each group in the KnowBe4 account. |
| [knowbe4_phishing_campaign](tables/knowbe4_phishing_campaign.md) | Phishing campaigns configured in the KnowBe4 account. |
| [knowbe4_phishing_recipient](tables/knowbe4_phishing_recipient.md) | Phishing security test recipient results in the KnowBe4 account. |
| [knowbe4_phishing_security_test](tables/knowbe4_phishing_security_test.md) | Phishing security tests (PSTs) run within KnowBe4 phishing campaigns. |
| [knowbe4_training_campaign](tables/knowbe4_training_campaign.md) | Training campaigns configured in the KnowBe4 account. |
| [knowbe4_training_enrollment](tables/knowbe4_training_enrollment.md) | Training enrollment records for users in the KnowBe4 account. |
| [knowbe4_training_policy](tables/knowbe4_training_policy.md) | Uploaded compliance policies in the KnowBe4 account. |
| [knowbe4_training_store_purchase](tables/knowbe4_training_store_purchase.md) | Training modules purchased from the KnowBe4 content store. |
| [knowbe4_user](tables/knowbe4_user.md) | Users registered in the KnowBe4 account. |
| [knowbe4_user_risk_score_history](tables/knowbe4_user_risk_score_history.md) | Historical risk score records for each user in the KnowBe4 account. |
