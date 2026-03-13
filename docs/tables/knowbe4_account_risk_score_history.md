# Table: knowbe4_account_risk_score_history

Historical risk score records for the KnowBe4 account.

## Examples

### Basic info

```sql
select
  date,
  risk_score
from
  knowbe4_account_risk_score_history;
```

### Recent risk score trend

```sql
select
  date,
  risk_score
from
  knowbe4_account_risk_score_history
order by
  date desc
limit 30;
```

### Highest risk score recorded

```sql
select
  date,
  risk_score
from
  knowbe4_account_risk_score_history
order by
  risk_score desc
limit 1;
```
