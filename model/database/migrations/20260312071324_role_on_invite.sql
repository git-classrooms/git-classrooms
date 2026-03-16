-- +goose Up
ALTER TABLE "public"."classroom_invitations"
ADD COLUMN "role" SMALLINT;

UPDATE classroom_invitations as ci
SET role = roles.role
FROM
  (SELECT ci.id as id, uc.role as role
  FROM classroom_invitations AS ci
  INNER JOIN user_classrooms as uc ON ci.classroom_id = uc.classroom_id
  INNER JOIN users as u ON u.id = uc.user_id
  WHERE ci.email = u.gitlab_email
  ORDER BY ci.id) as roles
WHERE roles.id = ci.id;

UPDATE "public"."classroom_invitations"
SET "role" = 2
WHERE "role" IS NULL;

ALTER TABLE "public"."classroom_invitations"
ALTER COLUMN "role" SET NOT NULL;

-- +goose Down
ALTER TABLE "public"."classroom_invitations"
DROP COLUMN "role";
