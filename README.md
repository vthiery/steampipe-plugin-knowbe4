# KnowBe4 Plugin for Steampipe

Use SQL to query users, phishing campaigns, training enrollments and more from [KnowBe4](https://www.knowbe4.com).

- **[Get started →](https://hub.steampipe.io/plugins/vthiery/knowbe4/tables)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/vthiery/knowbe4/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-knowbe4/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install ghcr.io/vthiery/knowbe4
```

Configure your credentials in `~/.steampipe/config/knowbe4.spc`:

```hcl
connection "knowbe4" {
  plugin = "ghcr.io/vthiery/knowbe4"

  # KnowBe4 API token (required).
  # You can generate one at https://training.knowbe4.com/account/info under "API".
  # api_token = "your-api-token"

  # KnowBe4 API region (optional, defaults to "us").
  # Valid values: "us", "eu", "ca", "uk", "de"
  # api_region = "us"
}
```

Run a query:

```sql
select
  id,
  email,
  phish_prone_percentage,
  current_risk_score
from
  knowbe4_user
where
  status = 'active'
order by
  current_risk_score desc
limit 10;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```shell
git clone https://github.com/vthiery/steampipe-plugin-knowbe4.git
cd steampipe-plugin-knowbe4
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make install
```

Configure the plugin:

```shell
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/knowbe4.spc
```

Try it!

```shell
steampipe query
> .inspect knowbe4
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [KnowBe4 Reporting API](https://developer.knowbe4.com/rest/reporting)
