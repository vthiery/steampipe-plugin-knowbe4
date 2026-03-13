# Table: knowbe4_training_store_purchase

Training modules purchased from the KnowBe4 content store.

## Examples

### Basic info

```sql
select
  store_purchased_id,
  name,
  type,
  duration,
  publisher,
  retired
from
  knowbe4_training_store_purchase
order by
  purchase_date desc;
```

### Active (non-retired) modules

```sql
select
  store_purchased_id,
  name,
  type,
  duration
from
  knowbe4_training_store_purchase
where
  not retired
order by
  name;
```
