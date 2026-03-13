# Table: knowbe4_phishing_security_test

Phishing security tests (PSTs) run within KnowBe4 phishing campaigns.

## Examples

### Basic info

```sql
select
  pst_id,
  campaign_id,
  name,
  status,
  phish_prone_percentage,
  started_at
from
  knowbe4_phishing_security_test
order by
  started_at desc;
```

### PSTs for a specific campaign

```sql
select
  pst_id,
  name,
  status,
  scheduled_count,
  clicked_count,
  phish_prone_percentage
from
  knowbe4_phishing_security_test
where
  campaign_id = 53333;
```

### Most recent PST click rates

```sql
select
  pst_id,
  name,
  started_at,
  scheduled_count,
  clicked_count,
  round((clicked_count::numeric / nullif(delivered_count, 0)) * 100, 1) as click_rate_pct
from
  knowbe4_phishing_security_test
where
  status = 'Closed'
order by
  started_at desc
limit 10;
```
