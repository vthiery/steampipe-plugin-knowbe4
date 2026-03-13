# Table: knowbe4_phishing_recipient

Phishing security test recipient results.

Note: `pst_id` must be specified in the `where` clause.

## Examples

### Basic info

```sql
select
  recipient_id,
  "user" ->> 'email' as email,
  clicked_at,
  reported_at,
  ip_location,
  browser,
  os
from
  knowbe4_phishing_recipient
where
  pst_id = 93104;
```

### Users who clicked but didn't report

```sql
select
  recipient_id,
  "user" ->> 'email' as email,
  clicked_at,
  ip_location
from
  knowbe4_phishing_recipient
where
  pst_id = 93104
  and clicked_at is not null
  and reported_at is null;
```

### Get a specific recipient result

```sql
select
  *
from
  knowbe4_phishing_recipient
where
  pst_id = 93104
  and recipient_id = 5781808;
```
