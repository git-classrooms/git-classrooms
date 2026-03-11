-- +goose Up
ALTER TABLE "public"."classrooms"
ADD COLUMN "teaching_group_id" BIGINT;

-- +goose Down
ALTER TABLE "public"."classrooms"
DROP COLUMN "teaching_group_id";
