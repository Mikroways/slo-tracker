schema "public" {
  comment = "standard public schema"
}

table "incidents" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "sli_name" {
    null = true
    type = text
  }
  column "slo_id" {
    null = true
    type = bigint
  }
  column "alertsource" {
    null = true
    type = text
  }
  column "state" {
    null = true
    type = text
  }
  column "error_budget_spent" {
    null = true
    type = numeric
  }
  column "real_error_budget" {
    null = true
    type = numeric
  }
  column "mark_false_positive" {
    null = true
    type = boolean
  }
  primary_key {
    columns = [column.id]
  }
}

table "slos" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "slo_name" {
    null = false
    type = text
  }
  column "target_slo" {
    null = true
    type = numeric
  }
  column "holidays_enabled" {
    null = true
    type = boolean
  }
  primary_key {
    columns = [column.id]
  }
  unique "uni_slos_slo_name" {
    columns = [column.slo_name]
  }
}

table "store_working_schedules" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "slo_id" {
    null = true
    type = bigint
  }
  column "weekday" {
    null = true
    type = bigint
  }
  column "open_hour" {
    null = true
    type = time
  }
  column "close_hour" {
    null = true
    type = time
  }
  primary_key {
    columns = [column.id]
  }
}
