-- +goose Up
ALTER TABLE "public"."classroom_invitations"
ADD COLUMN "role" SMALLINT;

-- +goose Down
ALTER TABLE "public"."classroom_invitations"
DROP COLUMN "role";
