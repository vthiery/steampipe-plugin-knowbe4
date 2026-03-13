# Table: knowbe4_group

User groups defined in the KnowBe4 account.

## Examples

### Basic info

```sql
select
  id,
  name,
  group_type,
  member_count,
  current_risk_score
from
  knowbe4_group;
```

### List all active groups

```sql
select
  id,
  name,
  group_type,
  member_count,
  current_risk_score
from
  knowbe4_group
where
  status = 'active'
order by
  member_count desc;
```

### Groups with highest risk score

```sql
select
  id,
  name,
  member_count,
  current_risk_score
from
  knowbe4_group
where
  status = 'active'
order by
  current_risk_score desc
limit 10;
```

### Get a specific group by ID

```sql
select
  id,
  name,
  group_type,
  member_count
from
  knowbe4_group
where
  id = 185273;
```
