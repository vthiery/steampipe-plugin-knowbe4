# Table: knowbe4_group_risk_score_history

Historical risk score records for groups in the KnowBe4 account.

Note: `group_id` must be specified in the `where` clause.

## Examples

### Basic info

```sql
select
  group_id,
  date,
  risk_score
from
  knowbe4_group_risk_score_history
where
  group_id = 185273
order by
  date desc;
```
