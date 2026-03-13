# Table: knowbe4_user_risk_score_history

Historical risk score records for users in the KnowBe4 account.

Note: `user_id` must be specified in the `where` clause.

## Examples

### Basic info

```sql
select
  user_id,
  date,
  risk_score
from
  knowbe4_user_risk_score_history
where
  user_id = 807836
order by
  date desc;
```
