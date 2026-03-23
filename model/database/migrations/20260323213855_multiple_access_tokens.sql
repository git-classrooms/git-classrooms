-- +goose Up
CREATE TABLE "public"."classroom_tokens" (
    "id" UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    "group_access_token_id" BIGINT NOT NULL,
    "group_access_token" TEXT NOT NULL,
    "group_access_token_created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "classroom_id" UUID NOT NULL,
    CONSTRAINT "fk_classrooms_tokens" FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE
);

INSERT INTO classroom_tokens ("group_access_token_id", "group_access_token", "group_access_token_created_at", "classroom_id")
SELECT "group_access_token_id", "group_access_token", "group_access_token_created_at", "id" FROM classrooms;

ALTER TABLE "public"."classrooms"
DROP COLUMN "group_access_token",
DROP COLUMN "group_access_token_created_at",
DROP COLUMN "group_access_token_id";

-- +goose Down
ALTER TABLE "public"."classrooms"
ADD COLUMN "group_access_token_id" BIGINT,
ADD COLUMN "group_access_token" TEXT,
ADD COLUMN "group_access_token_created_at" TIMESTAMP WITH TIME ZONE:

UPDATE classrooms c
SET "group_access_token_id" = ct."group_access_token_id",
    "group_access_token" = ct."group_access_token",
    "group_access_token_created_at" = ct."group_access_token_created_at"
FROM classroom_tokens ct
WHERE ct."classroom_id" = c."id";

DROP TABLE "public"."classroom_tokens";

ALTER TABLE classrooms
ALTER COLUMN "group_access_token_id" SET NOT NULL,
ALTER COLUMN "group_access_token" SET NOT NULL,
ALTER COLUMN "group_access_token_created_at" SET NOT NULL;
