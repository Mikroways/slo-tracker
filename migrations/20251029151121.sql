-- Create "incidents" table
CREATE TABLE "public"."incidents" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "sli_name" text NULL,
  "slo_id" bigint NULL,
  "alertsource" text NULL,
  "state" text NULL,
  "error_budget_spent" numeric NULL,
  "real_error_budget" numeric NULL,
  "mark_false_positive" boolean NULL,
  PRIMARY KEY ("id")
);
-- Create "slos" table
CREATE TABLE "public"."slos" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "slo_name" text NOT NULL,
  "target_slo" numeric NULL,
  "holidays_enabled" boolean NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_slos_slo_name" UNIQUE ("slo_name")
);
-- Create "store_working_schedules" table
CREATE TABLE "public"."store_working_schedules" (
  "id" bigserial NOT NULL,
  "slo_id" bigint NULL,
  "weekday" bigint NULL,
  "open_hour" time NULL,
  "close_hour" time NULL,
  PRIMARY KEY ("id")
);
