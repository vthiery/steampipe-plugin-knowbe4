# Table: knowbe4_account

KnowBe4 account details including subscription level, number of seats, and current risk score.

## Examples

### Basic info

```sql
select
  name,
  type,
  subscription_level,
  subscription_end_date,
  number_of_seats,
  current_risk_score
from
  knowbe4_account;
```

### List account domains

```sql
select
  name,
  jsonb_array_elements_text(domains::jsonb) as domain
from
  knowbe4_account;
```
