-- +goose Up
ALTER TABLE "public"."assignment_projects"
ADD COLUMN "http_url_to_repo" TEXT,
ADD COLUMN "ssh_url_to_repo" TEXT;


-- +goose Down
ALTER TABLE "public"."assignment_projects"
DROP COLUMN "http_url_to_repo",
DROP COLUMN "ssh_url_to_repo";
