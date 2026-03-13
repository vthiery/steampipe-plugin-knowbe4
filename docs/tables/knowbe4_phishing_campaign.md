# Table: knowbe4_phishing_campaign

Phishing campaigns configured in the KnowBe4 account.

## Examples

### Basic info

```sql
select
  campaign_id,
  name,
  status,
  frequency,
  psts_count
from
  knowbe4_phishing_campaign;
```

### List active campaigns

```sql
select
  campaign_id,
  name,
  last_phish_prone_percentage,
  psts_count
from
  knowbe4_phishing_campaign
where
  status = 'Active'
order by
  last_phish_prone_percentage desc;
```

### Get a specific phishing campaign

```sql
select
  campaign_id,
  name,
  status,
  frequency,
  last_phish_prone_percentage,
  psts_count
from
  knowbe4_phishing_campaign
where
  campaign_id = 53333;
```
