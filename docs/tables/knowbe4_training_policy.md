# Table: knowbe4_training_policy

Uploaded compliance policies in the KnowBe4 account.

## Examples

### Basic info

```sql
select
  policy_id,
  name,
  content_type,
  minimum_time,
  default_language,
  status
from
  knowbe4_training_policy
order by
  name;
```
