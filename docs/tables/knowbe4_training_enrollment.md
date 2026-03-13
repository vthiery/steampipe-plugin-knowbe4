# Table: knowbe4_training_enrollment

Training enrollment records for users in the KnowBe4 account.

## Examples

### Basic info

```sql
select
  enrollment_id,
  module_name,
  "user" ->> 'email' as email,
  campaign_name,
  status,
  enrollment_date
from
  knowbe4_training_enrollment;
```

### List all incomplete enrollments

```sql
select
  enrollment_id,
  module_name,
  "user" ->> 'email' as email,
  campaign_name,
  status,
  enrollment_date
from
  knowbe4_training_enrollment
where
  status != 'Passed'
  and status != 'Exempted'
order by
  enrollment_date asc;
```

### Enrollments for a specific campaign

```sql
select
  enrollment_id,
  module_name,
  "user" ->> 'email' as email,
  status,
  completion_date
from
  knowbe4_training_enrollment
where
  campaign_id = 45472
order by
  status;
```

### Enrollments for a specific user

```sql
select
  enrollment_id,
  module_name,
  campaign_name,
  status,
  completion_date,
  time_spent
from
  knowbe4_training_enrollment
where
  user_id = 807836
order by
  enrollment_date desc;
```

### Completion rate per module

```sql
select
  module_name,
  count(*) as total,
  count(*) filter (where status = 'Passed') as passed,
  round(count(*) filter (where status = 'Passed')::numeric / count(*) * 100, 1) as completion_pct
from
  knowbe4_training_enrollment
group by
  module_name
order by
  completion_pct asc;
```
