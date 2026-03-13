# Table: knowbe4_training_campaign

Training campaigns configured in the KnowBe4 account.

## Examples

### Basic info

```sql
select
  campaign_id,
  name,
  status,
  start_date,
  end_date,
  completion_percentage
from
  knowbe4_training_campaign;
```

### List all active training campaigns

```sql
select
  campaign_id,
  name,
  status,
  start_date,
  end_date,
  completion_percentage
from
  knowbe4_training_campaign
where
  status = 'Active'
order by
  completion_percentage asc;
```

### Campaigns with low completion rates

```sql
select
  campaign_id,
  name,
  status,
  completion_percentage,
  end_date
from
  knowbe4_training_campaign
where
  status = 'Active'
  and completion_percentage < 50
order by
  end_date asc;
```
