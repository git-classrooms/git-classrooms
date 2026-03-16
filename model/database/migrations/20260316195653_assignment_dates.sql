-- +goose Up
create table "public"."assignment_dates" (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "assignment_id" UUID NOT NULL,
    "due_date" TIMESTAMP WITH TIME ZONE,
    "closed" BOOLEAN DEFAULT FALSE,
    CONSTRAINT "fk_assignments_assignment_dates" FOREIGN KEY ("assignment_id") REFERENCES "public"."assignments"("id") ON DELETE CASCADE
);

INSERT INTO "assignment_dates" (assignment_id, due_date, closed, created_at, updated_at)
  SELECT a.id, a.due_date, a.closed, NOW(), NOW() FROM assignments AS a;

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
    "grading_j_unit_test_result" jsonb,
    constraint "fk_assignment_project_grading_dates_assignment_date" FOREIGN KEY (assignment_date_id) REFERENCES assignment_dates(id),
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
