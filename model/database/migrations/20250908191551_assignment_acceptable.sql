-- +goose Up
ALTER TABLE "public"."assignments" ADD COLUMN "acceptable_since" TIMESTAMP WITH TIME ZONE;

-- +goose Down
ALTER TABLE "public"."assignments" DROP COLUMN "acceptable_since";
