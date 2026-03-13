# Table: knowbe4_user

Users registered in the KnowBe4 account.

## Examples

### Basic info

```sql
select
  id,
  email,
  first_name,
  last_name,
  department,
  phish_prone_percentage
from
  knowbe4_user;
```

### List all active users

```sql
select
  id,
  email,
  first_name,
  last_name,
  department,
  phish_prone_percentage
from
  knowbe4_user
where
  status = 'active';
```

### Users with high phish-prone percentage

```sql
select
  email,
  first_name,
  last_name,
  department,
  phish_prone_percentage
from
  knowbe4_user
where
  phish_prone_percentage > 30
  and status = 'active'
order by
  phish_prone_percentage desc;
```

### Get a specific user by ID

```sql
select
  id,
  email,
  current_risk_score
from
  knowbe4_user
where
  id = 807836;
```

### Users who have never signed in

```sql
select
  email,
  first_name,
  last_name,
  joined_on
from
  knowbe4_user
where
  last_sign_in is null
  and status = 'active';
```
