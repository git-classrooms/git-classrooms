-- +goose Up
create table "public"."assignment_dates" (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "assignment_id" UUID NOT NULL,
    "description" TEXT NOT NULL,
    "due_date" TIMESTAMP WITH TIME ZONE,
    "closed" BOOLEAN DEFAULT FALSE,
    CONSTRAINT "fk_assignments_assignment_dates" FOREIGN KEY ("assignment_id") REFERENCES "public"."assignments"("id") ON DELETE CASCADE
);

INSERT INTO "assignment_dates" (assignment_id, due_date, closed, created_at, updated_at, description)
  SELECT a.id, a.due_date, a.closed, NOW(), NOW(), 'final assignment' FROM assignments AS a;

UPDATE "assignment_dates" SET due_date = NOW() + '14 days' WHERE due_date IS NULL;

alter table "public"."assignment_dates" alter column "due_date" set not null;

alter table "public"."assignments" drop column "closed";

alter table "public"."assignments" drop column "due_date";

create table "public"."assignment_dates_manual_grading_rubrics" (
    "manual_grading_rubric_id" uuid not null default uuid_generate_v4(),
    "assignment_date_id" uuid not null default uuid_generate_v4(),
    PRIMARY KEY ("manual_grading_rubric_id", "assignment_date_id"),
    constraint "fk_assignment_dates_manual_grading_rubrics_assignment_date" FOREIGN KEY (assignment_date_id) REFERENCES assignment_dates(id) ON DELETE CASCADE,
    constraint "fk_assignment_dates_manual_grading_rubrics_manual_grading_rubri" FOREIGN KEY (manual_grading_rubric_id) REFERENCES manual_grading_rubrics(id) ON DELETE CASCADE
);

INSERT INTO "public"."assignment_dates_manual_grading_rubrics"
  SELECT amg.manual_grading_rubric_id, ad.id
  FROM assignments AS a
  INNER JOIN assignment_dates AS ad ON a.id = ad.assignment_id
  INNER JOIN assignment_manual_grading_rubrics as amg ON amg.assignment_id = a.id;

drop table "public"."assignment_manual_grading_rubrics";

create table "public"."assignment_project_grading_dates" (
    "id" uuid PRIMARY KEY not null default uuid_generate_v4(),
    "assignment_project_id" uuid not null,
    "assignment_date_id" uuid not null,
    "branchname" text,
    "grading_j_unit_test_result" jsonb,
    constraint "fk_assignment_dates_assignment_project_grading_date" FOREIGN KEY (assignment_date_id) REFERENCES assignment_dates(id) ON DELETE CASCADE,
    constraint "fk_assignment_projects_gradings" FOREIGN KEY (assignment_project_id) REFERENCES assignment_projects(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_project_date ON public.assignment_project_grading_dates USING btree (assignment_project_id, assignment_date_id);

INSERT INTO "public"."assignment_project_grading_dates" (assignment_project_id, assignment_date_id, grading_j_unit_test_result)
  SELECT ap.id, ad.id, ap.grading_j_unit_test_result
  FROM assignment_projects as ap
  INNER JOIN assignments as a ON ap.assignment_id = a.id
  INNER JOIN assignment_dates as ad ON ad.assignment_id = a.id;

alter table "public"."assignment_projects" drop column "grading_j_unit_test_result";

alter table "public"."assignment_junit_tests" add column "assignment_date_id" uuid;
alter table "public"."assignment_junit_tests" add constraint "fk_assignment_dates_j_unit_tests" FOREIGN KEY (assignment_date_id) REFERENCES assignment_dates(id) ON DELETE CASCADE;

UPDATE assignment_junit_tests AS ajt
SET assignment_date_id = ad.id
FROM assignment_dates as ad
WHERE ajt.assignment_id = ad.assignment_id;

alter table "public"."assignment_junit_tests" ALTER column "assignment_date_id" SET NOT NULL;

alter table "public"."assignment_junit_tests" drop column "assignment_id";

alter table "public"."manual_grading_results" add column "assignment_project_grading_date_id" uuid;
CREATE UNIQUE INDEX "idx_unique_assignment_date_assignmentjunittestName" ON public.assignment_junit_tests USING btree (name, assignment_date_id);
alter table "public"."manual_grading_results" add constraint "fk_assignment_project_grading_dates_grading_manual_results" FOREIGN KEY (assignment_project_grading_date_id) REFERENCES assignment_project_grading_dates(id) ON DELETE CASCADE;

UPDATE manual_grading_results AS mgr
SET assignment_project_grading_date_id = apgd.id
FROM assignment_project_grading_dates as apgd
WHERE apgd.assignment_project_id = mgr.assignment_project_id;

alter table "public"."manual_grading_results" ALTER column "assignment_project_grading_date_id" SET NOT NULL;

alter table "public"."manual_grading_results" drop column "assignment_project_id";

-- +goose Down

-- 1. manual_grading_results: restore assignment_project_id
alter table "public"."manual_grading_results" add column "assignment_project_id" uuid;

UPDATE manual_grading_results AS mgr
SET assignment_project_id = apgd.assignment_project_id
FROM assignment_project_grading_dates AS apgd
WHERE mgr.assignment_project_grading_date_id = apgd.id;

alter table "public"."manual_grading_results" ALTER column "assignment_project_id" SET NOT NULL;
alter table "public"."manual_grading_results" add constraint "fk_assignment_projects_grading_manual_results" FOREIGN KEY (assignment_project_id) REFERENCES assignment_projects(id) ON DELETE CASCADE;
alter table "public"."manual_grading_results" drop constraint "fk_assignment_project_grading_dates_grading_manual_results";
alter table "public"."manual_grading_results" drop column "assignment_project_grading_date_id";

-- 2. assignment_junit_tests: restore assignment_id
DROP INDEX "public"."idx_unique_assignment_date_assignmentjunittestName";

alter table "public"."assignment_junit_tests" add column "assignment_id" uuid;

UPDATE assignment_junit_tests AS ajt
SET assignment_id = ad.assignment_id
FROM assignment_dates AS ad
WHERE ajt.assignment_date_id = ad.id;

alter table "public"."assignment_junit_tests" ALTER column "assignment_id" SET NOT NULL;
alter table "public"."assignment_junit_tests" add constraint "fk_assignments_j_unit_tests" FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX "idx_unique_assignment_assignmentjunittestName" ON "public"."assignment_junit_tests" USING btree ("name", "assignment_id");

alter table "public"."assignment_junit_tests" drop constraint "fk_assignment_dates_j_unit_tests";
alter table "public"."assignment_junit_tests" drop column "assignment_date_id";

-- 3. assignment_projects: restore grading_j_unit_test_result
alter table "public"."assignment_projects" add column "grading_j_unit_test_result" jsonb;

UPDATE assignment_projects AS ap
SET grading_j_unit_test_result = apgd.grading_j_unit_test_result
FROM assignment_project_grading_dates AS apgd
WHERE apgd.assignment_project_id = ap.id;

-- 4. drop assignment_project_grading_dates
drop table "public"."assignment_project_grading_dates";

-- 5. restore assignment_manual_grading_rubrics
create table "public"."assignment_manual_grading_rubrics" (
    "manual_grading_rubric_id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    "assignment_id" UUID NOT NULL DEFAULT uuid_generate_v4(),
    PRIMARY KEY ("manual_grading_rubric_id", "assignment_id"),
    CONSTRAINT "fk_assignment_manual_grading_rubrics_manual_grading_rubric" FOREIGN KEY ("manual_grading_rubric_id") REFERENCES "public"."manual_grading_rubrics"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_assignment_manual_grading_rubrics_assignment" FOREIGN KEY ("assignment_id") REFERENCES "public"."assignments"("id") ON DELETE CASCADE
);

INSERT INTO "public"."assignment_manual_grading_rubrics" (manual_grading_rubric_id, assignment_id)
  SELECT DISTINCT admgr.manual_grading_rubric_id, ad.assignment_id
  FROM assignment_dates_manual_grading_rubrics AS admgr
  INNER JOIN assignment_dates AS ad ON admgr.assignment_date_id = ad.id;

-- 6. drop assignment_dates_manual_grading_rubrics
drop table "public"."assignment_dates_manual_grading_rubrics";

-- 7. assignments: restore due_date and closed
alter table "public"."assignments" add column "due_date" TIMESTAMP WITH TIME ZONE;
alter table "public"."assignments" add column "closed" BOOLEAN DEFAULT FALSE;

UPDATE assignments AS a
SET due_date = ad.due_date, closed = ad.closed
FROM (
  SELECT DISTINCT ON (assignment_id) assignment_id, due_date, closed
  FROM assignment_dates
  ORDER BY assignment_id, created_at ASC
) AS ad
WHERE a.id = ad.assignment_id;

-- 8. drop assignment_dates
drop table "public"."assignment_dates";
