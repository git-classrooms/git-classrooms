-- +goose Up
ALTER TABLE "public"."classrooms" ADD COLUMN "invite_code" UUID NOT NULL DEFAULT uuid_generate_v4();

-- +goose Down
ALTER TABLE "public"."classrooms" DROP COLUMN "invite_code";
