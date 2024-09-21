create extension if not exists "uuid-ossp" with schema "public" version '1.1';

create sequence "public"."user_avatars_user_id_seq";

create table "public"."assignment_junit_tests" (
    "id" uuid not null default uuid_generate_v4(),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "name" text not null,
    "assignment_id" uuid not null,
    "score" bigint not null
);


create table "public"."assignment_manual_grading_rubrics" (
    "manual_grading_rubric_id" uuid not null default uuid_generate_v4(),
    "assignment_id" uuid not null default uuid_generate_v4()
);


create table "public"."assignment_projects" (
    "id" uuid not null default uuid_generate_v4(),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "team_id" uuid not null,
    "assignment_id" uuid not null,
    "project_status" text not null default 'pending'::text,
    "project_id" bigint,
    "grading_j_unit_test_result" jsonb
);


create table "public"."assignments" (
    "id" uuid not null default uuid_generate_v4(),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "classroom_id" uuid not null,
    "template_project_id" bigint not null,
    "name" text not null,
    "description" text,
    "due_date" timestamp with time zone,
    "closed" boolean default false,
    "grading_j_unit_auto_grading_active" boolean
);


create table "public"."classroom_invitations" (
    "id" uuid not null default uuid_generate_v4(),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "status" smallint not null,
    "classroom_id" uuid not null,
    "email" text not null,
    "expiry_date" timestamp with time zone not null
);


create table "public"."classrooms" (
    "id" uuid not null default uuid_generate_v4(),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "name" text not null,
    "description" text,
    "owner_id" bigint not null,
    "create_teams" boolean not null,
    "max_team_size" bigint not null default 1,
    "max_teams" bigint not null default 0,
    "group_id" bigint not null,
    "group_access_token_id" bigint not null,
    "group_access_token" text not null,
    "students_view_all_projects" boolean not null,
    "archived" boolean not null default false,
    "potentially_deleted" boolean not null default false
);


create table "public"."fiber_storage" (
    "k" character varying(64) not null default ''::character varying,
    "v" bytea not null,
    "e" bigint not null default '0'::bigint
);


create table "public"."manual_grading_results" (
    "id" uuid not null default uuid_generate_v4(),
    "rubric_id" uuid not null,
    "assignment_project_id" uuid not null,
    "score" bigint not null,
    "feedback" text
);


create table "public"."manual_grading_rubrics" (
    "id" uuid not null default uuid_generate_v4(),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "name" text not null,
    "description" text not null,
    "classroom_id" uuid not null,
    "max_score" bigint not null
);


create table "public"."teams" (
    "id" uuid not null default uuid_generate_v4(),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "name" text not null,
    "group_id" bigint not null,
    "classroom_id" uuid not null
);


create table "public"."user_avatars" (
    "user_id" bigint not null default nextval('user_avatars_user_id_seq'::regclass),
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone,
    "avatar_url" text,
    "fallback_avatar_url" text
);


create table "public"."user_classrooms" (
    "user_id" bigint not null,
    "classroom_id" uuid not null,
    "team_id" uuid,
    "role" smallint not null
);


create table "public"."users" (
    "id" bigint not null,
    "gitlab_username" text not null,
    "gitlab_email" text not null,
    "name" text not null,
    "created_at" timestamp with time zone,
    "updated_at" timestamp with time zone
);


alter sequence "public"."user_avatars_user_id_seq" owned by "public"."user_avatars"."user_id";

CREATE UNIQUE INDEX assignment_junit_tests_pkey ON public.assignment_junit_tests USING btree (id);

CREATE UNIQUE INDEX assignment_manual_grading_rubrics_pkey ON public.assignment_manual_grading_rubrics USING btree (manual_grading_rubric_id, assignment_id);

CREATE UNIQUE INDEX assignment_projects_pkey ON public.assignment_projects USING btree (id);

CREATE UNIQUE INDEX assignments_pkey ON public.assignments USING btree (id);

CREATE UNIQUE INDEX classroom_invitations_pkey ON public.classroom_invitations USING btree (id);

CREATE UNIQUE INDEX classrooms_pkey ON public.classrooms USING btree (id);

CREATE INDEX e ON public.fiber_storage USING btree (e);

CREATE UNIQUE INDEX fiber_storage_pkey ON public.fiber_storage USING btree (k);

CREATE UNIQUE INDEX "idx_unique_assignment_assignmentjunittestName" ON public.assignment_junit_tests USING btree (name, assignment_id);

CREATE UNIQUE INDEX "idx_unique_classroom_assignmentName" ON public.assignments USING btree (classroom_id, name);

CREATE INDEX idx_user_classrooms_team_id ON public.user_classrooms USING btree (team_id);

CREATE UNIQUE INDEX manual_grading_results_pkey ON public.manual_grading_results USING btree (id);

CREATE UNIQUE INDEX manual_grading_rubrics_pkey ON public.manual_grading_rubrics USING btree (id);

CREATE UNIQUE INDEX teams_pkey ON public.teams USING btree (id);

CREATE UNIQUE INDEX user_avatars_pkey ON public.user_avatars USING btree (user_id);

CREATE UNIQUE INDEX user_classrooms_pkey ON public.user_classrooms USING btree (user_id, classroom_id);

CREATE UNIQUE INDEX users_gitlab_email_key ON public.users USING btree (gitlab_email);

CREATE UNIQUE INDEX users_gitlab_username_key ON public.users USING btree (gitlab_username);

CREATE UNIQUE INDEX users_pkey ON public.users USING btree (id);

alter table "public"."assignment_junit_tests" add constraint "assignment_junit_tests_pkey" PRIMARY KEY using index "assignment_junit_tests_pkey";

alter table "public"."assignment_manual_grading_rubrics" add constraint "assignment_manual_grading_rubrics_pkey" PRIMARY KEY using index "assignment_manual_grading_rubrics_pkey";

alter table "public"."assignment_projects" add constraint "assignment_projects_pkey" PRIMARY KEY using index "assignment_projects_pkey";

alter table "public"."assignments" add constraint "assignments_pkey" PRIMARY KEY using index "assignments_pkey";

alter table "public"."classroom_invitations" add constraint "classroom_invitations_pkey" PRIMARY KEY using index "classroom_invitations_pkey";

alter table "public"."classrooms" add constraint "classrooms_pkey" PRIMARY KEY using index "classrooms_pkey";

alter table "public"."fiber_storage" add constraint "fiber_storage_pkey" PRIMARY KEY using index "fiber_storage_pkey";

alter table "public"."manual_grading_results" add constraint "manual_grading_results_pkey" PRIMARY KEY using index "manual_grading_results_pkey";

alter table "public"."manual_grading_rubrics" add constraint "manual_grading_rubrics_pkey" PRIMARY KEY using index "manual_grading_rubrics_pkey";

alter table "public"."teams" add constraint "teams_pkey" PRIMARY KEY using index "teams_pkey";

alter table "public"."user_avatars" add constraint "user_avatars_pkey" PRIMARY KEY using index "user_avatars_pkey";

alter table "public"."user_classrooms" add constraint "user_classrooms_pkey" PRIMARY KEY using index "user_classrooms_pkey";

alter table "public"."users" add constraint "users_pkey" PRIMARY KEY using index "users_pkey";

alter table "public"."assignment_junit_tests" add constraint "fk_assignments_j_unit_tests" FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE not valid;

alter table "public"."assignment_junit_tests" validate constraint "fk_assignments_j_unit_tests";

alter table "public"."assignment_manual_grading_rubrics" add constraint "fk_assignment_manual_grading_rubrics_assignment" FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE not valid;

alter table "public"."assignment_manual_grading_rubrics" validate constraint "fk_assignment_manual_grading_rubrics_assignment";

alter table "public"."assignment_manual_grading_rubrics" add constraint "fk_assignment_manual_grading_rubrics_manual_grading_rubric" FOREIGN KEY (manual_grading_rubric_id) REFERENCES manual_grading_rubrics(id) ON DELETE CASCADE not valid;

alter table "public"."assignment_manual_grading_rubrics" validate constraint "fk_assignment_manual_grading_rubrics_manual_grading_rubric";

alter table "public"."assignment_projects" add constraint "fk_assignments_projects" FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE not valid;

alter table "public"."assignment_projects" validate constraint "fk_assignments_projects";

alter table "public"."assignment_projects" add constraint "fk_teams_assignment_projects" FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE not valid;

alter table "public"."assignment_projects" validate constraint "fk_teams_assignment_projects";

alter table "public"."assignments" add constraint "fk_classrooms_assignments" FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE not valid;

alter table "public"."assignments" validate constraint "fk_classrooms_assignments";

alter table "public"."classroom_invitations" add constraint "fk_classrooms_invitations" FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE not valid;

alter table "public"."classroom_invitations" validate constraint "fk_classrooms_invitations";

alter table "public"."classrooms" add constraint "fk_users_owned_classrooms" FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE not valid;

alter table "public"."classrooms" validate constraint "fk_users_owned_classrooms";

alter table "public"."manual_grading_results" add constraint "fk_assignment_projects_grading_manual_results" FOREIGN KEY (assignment_project_id) REFERENCES assignment_projects(id) ON DELETE CASCADE not valid;

alter table "public"."manual_grading_results" validate constraint "fk_assignment_projects_grading_manual_results";

alter table "public"."manual_grading_results" add constraint "fk_manual_grading_rubrics_results" FOREIGN KEY (rubric_id) REFERENCES manual_grading_rubrics(id) ON DELETE CASCADE not valid;

alter table "public"."manual_grading_results" validate constraint "fk_manual_grading_rubrics_results";

alter table "public"."manual_grading_rubrics" add constraint "fk_classrooms_manual_grading_rubrics" FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE not valid;

alter table "public"."manual_grading_rubrics" validate constraint "fk_classrooms_manual_grading_rubrics";

alter table "public"."teams" add constraint "fk_classrooms_teams" FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE not valid;

alter table "public"."teams" validate constraint "fk_classrooms_teams";

alter table "public"."user_avatars" add constraint "fk_users_git_lab_avatar" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE not valid;

alter table "public"."user_avatars" validate constraint "fk_users_git_lab_avatar";

alter table "public"."user_classrooms" add constraint "fk_classrooms_member" FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE not valid;

alter table "public"."user_classrooms" validate constraint "fk_classrooms_member";

alter table "public"."user_classrooms" add constraint "fk_teams_member" FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE not valid;

alter table "public"."user_classrooms" validate constraint "fk_teams_member";

alter table "public"."user_classrooms" add constraint "fk_users_classrooms" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE not valid;

alter table "public"."user_classrooms" validate constraint "fk_users_classrooms";

alter table "public"."users" add constraint "users_gitlab_email_key" UNIQUE using index "users_gitlab_email_key";

alter table "public"."users" add constraint "users_gitlab_username_key" UNIQUE using index "users_gitlab_username_key";


