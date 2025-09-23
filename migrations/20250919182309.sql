-- Create "incident_reqs" table
CREATE TABLE "public"."incident_reqs" (
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
-- Create index "idx_incident_reqs_deleted_at" to table: "incident_reqs"
CREATE INDEX "idx_incident_reqs_deleted_at" ON "public"."incident_reqs" ("deleted_at");
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
-- Create index "idx_incidents_deleted_at" to table: "incidents"
CREATE INDEX "idx_incidents_deleted_at" ON "public"."incidents" ("deleted_at");
-- Create "slos" table
CREATE TABLE "public"."slos" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "slo_name" text NOT NULL,
  "open_hour" time NULL,
  "close_hour" time NULL,
  "target_slo" numeric NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_slos_slo_name" UNIQUE ("slo_name")
);
-- Create index "idx_slos_deleted_at" to table: "slos"
CREATE INDEX "idx_slos_deleted_at" ON "public"."slos" ("deleted_at");
