-- +goose Up
ALTER TABLE "public"."user_classrooms"
ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE,
ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE;

UPDATE "public"."user_classrooms"
SET "created_at" = NOW(), "updated_at" = NOW();

-- +goose Down
ALTER TABLE "public"."assignment_projects"
ADD COLUMN "created_at" TEXT,
ADD COLUMN "updated_at" TEXT;
