-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "public"."users" (
    "id" BIGINT PRIMARY KEY,
    "gitlab_username" TEXT NOT NULL UNIQUE,
    "gitlab_email" TEXT NOT NULL UNIQUE,
    "avatar_url" TEXT,
    "fallback_avatar_url" TEXT,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE "public"."classrooms" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "owner_id" BIGINT NOT NULL,
    "create_teams" BOOLEAN NOT NULL,
    "max_team_size" BIGINT NOT NULL DEFAULT 1,
    "max_teams" BIGINT NOT NULL DEFAULT 0,
    "group_id" BIGINT NOT NULL,
    "group_access_token_id" BIGINT NOT NULL,
    "group_access_token" TEXT NOT NULL,
    "group_access_token_created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "students_view_all_projects" BOOLEAN NOT NULL,
    "archived" BOOLEAN NOT NULL DEFAULT FALSE,
    "potentially_deleted" BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT "fk_users_owned_classrooms" FOREIGN KEY ("owner_id") REFERENCES "public"."users"("id") ON DELETE CASCADE
);

CREATE TABLE "public"."teams" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "name" TEXT NOT NULL,
    "group_id" BIGINT NOT NULL,
    "classroom_id" UUID NOT NULL,
    CONSTRAINT "fk_classrooms_teams" FOREIGN KEY ("classroom_id") REFERENCES "public"."classrooms"("id") ON DELETE CASCADE
);

CREATE TABLE "public"."user_classrooms" (
    "user_id" BIGINT NOT NULL,
    "classroom_id" UUID NOT NULL,
    "team_id" UUID,
    "role" SMALLINT NOT NULL,
    PRIMARY KEY ("user_id", "classroom_id"),
    CONSTRAINT "fk_users_classrooms" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_classrooms_member" FOREIGN KEY ("classroom_id") REFERENCES "public"."classrooms"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_teams_member" FOREIGN KEY ("team_id") REFERENCES "public"."teams"("id") ON DELETE SET NULL
);
CREATE INDEX "idx_user_classrooms_team_id" ON "public"."user_classrooms" USING btree ("team_id");

CREATE TABLE "public"."classroom_invitations" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "status" SMALLINT NOT NULL,
    "classroom_id" UUID NOT NULL,
    "email" TEXT NOT NULL,
    "expiry_date" TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT "fk_classrooms_invitations" FOREIGN KEY ("classroom_id") REFERENCES "public"."classrooms"("id") ON DELETE CASCADE
);

CREATE TABLE "public"."assignments" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "classroom_id" UUID NOT NULL,
    "template_project_id" BIGINT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "due_date" TIMESTAMP WITH TIME ZONE,
    "closed" BOOLEAN DEFAULT FALSE,
    "grading_j_unit_auto_grading_active" BOOLEAN,
    CONSTRAINT "fk_classrooms_assignments" FOREIGN KEY ("classroom_id") REFERENCES "public"."classrooms"("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX "idx_unique_classroom_assignmentName" ON "public"."assignments" USING btree ("classroom_id", "name");

CREATE TABLE "public"."assignment_projects" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "team_id" UUID NOT NULL,
    "assignment_id" UUID NOT NULL,
    "project_status" TEXT NOT NULL DEFAULT 'pending'::TEXT,
    "project_id" BIGINT,
    "grading_j_unit_test_result" JSONB,
    CONSTRAINT "fk_teams_assignment_projects" FOREIGN KEY ("team_id") REFERENCES "public"."teams"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_assignments_projects" FOREIGN KEY ("assignment_id") REFERENCES "public"."assignments"("id") ON DELETE CASCADE
);

CREATE TABLE "public"."manual_grading_rubrics" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "classroom_id" UUID NOT NULL,
    "max_score" BIGINT NOT NULL,
    CONSTRAINT "fk_classrooms_manual_grading_rubrics" FOREIGN KEY ("classroom_id") REFERENCES "public"."classrooms"("id") ON DELETE CASCADE
);

CREATE TABLE "public"."assignment_manual_grading_rubrics" (
    "manual_grading_rubric_id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "assignment_id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    PRIMARY KEY ("manual_grading_rubric_id", "assignment_id"),
    CONSTRAINT "fk_assignment_manual_grading_rubrics_manual_grading_rubric" FOREIGN KEY ("manual_grading_rubric_id") REFERENCES "public"."manual_grading_rubrics"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_assignment_manual_grading_rubrics_assignment" FOREIGN KEY ("assignment_id") REFERENCES "public"."assignments"("id") ON DELETE CASCADE
);

CREATE TABLE "public"."assignment_junit_tests" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "name" TEXT NOT NULL,
    "assignment_id" UUID NOT NULL,
    "score" BIGINT NOT NULL,
    CONSTRAINT "fk_assignments_j_unit_tests" FOREIGN KEY ("assignment_id") REFERENCES "public"."assignments"("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX "idx_unique_assignment_assignmentjunittestName" ON "public"."assignment_junit_tests" USING btree ("name", "assignment_id");

CREATE TABLE "public"."manual_grading_results" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "rubric_id" UUID NOT NULL,
    "assignment_project_id" UUID NOT NULL,
    "score" BIGINT NOT NULL,
    "feedback" TEXT,
    CONSTRAINT "fk_manual_grading_rubrics_results" FOREIGN KEY ("rubric_id") REFERENCES "public"."manual_grading_rubrics"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_assignment_projects_grading_manual_results" FOREIGN KEY ("assignment_project_id") REFERENCES "public"."assignment_projects"("id") ON DELETE CASCADE
);

-- +goose Down
DROP TABLE "public"."manual_grading_results";
DROP TABLE "public"."assignment_junit_tests";
DROP TABLE "public"."assignment_manual_grading_rubrics";
DROP TABLE "public"."manual_grading_rubrics";
DROP TABLE "public"."assignment_projects";
DROP TABLE "public"."assignments";
DROP TABLE "public"."classroom_invitations";
DROP TABLE "public"."user_classrooms";
DROP TABLE "public"."teams";
DROP TABLE "public"."classrooms";
DROP TABLE "public"."users";
