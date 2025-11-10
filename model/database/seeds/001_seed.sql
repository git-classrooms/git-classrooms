-- +goose Up
INSERT INTO "public"."users" ("id", "gitlab_username", "gitlab_email", "name", "created_at", "updated_at", "avatar_url", "fallback_avatar_url")
VALUES (
  2, 'tonitester', 'toni@tester.de', 'Toni Tester', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/e79b309c6c0ab9f6620ffa06b2ffcef89cbcc36f174a80a3dc62f2c195308ae3?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/e79b309c6c0ab9f6620ffa06b2ffcef89cbcc36f174a80a3dc62f2c195308ae3?s=80&d=identicon'
), (
  3, 'lauratester', 'laura@tester.de', 'Laura Tester', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/f882fbeb99363bd61790edd1845040152a8198ebcda9192e06a4ed9ed2496f48?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/f882fbeb99363bd61790edd1845040152a8198ebcda9192e06a4ed9ed2496f48?s=80&d=identicon'
), (
  4, 'harri', 'harald@hueftschmerz.de', 'Harald Hüftschmerz', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/b8565777dac214ba1a4e9623c16f3baf60ff1cdb22229e1d35c068841be294ef?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/b8565777dac214ba1a4e9623c16f3baf60ff1cdb22229e1d35c068841be294ef?s=80&d=identicon'
), (
  5, 'mustermanfred', 'manfred@muster.de', 'Manfred Muster', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/9b9613b9a2632819b9f5e5a85420fc80f6ad8dd235f312e64069610c503a62d6?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/9b9613b9a2632819b9f5e5a85420fc80f6ad8dd235f312e64069610c503a62d6?s=80&d=identicon'
), (
  6, 'hanshotfix', 'hans@hotfix.de', 'Hans Hotfix', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/c3d427b59dc3e5c60b14b53742e75fe39e6373981b628ed8e619ca94cfe687e3?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/c3d427b59dc3e5c60b14b53742e75fe39e6373981b628ed8e619ca94cfe687e3?s=80&d=identicon'
), (
  7, 'peterpush', 'peter@push.de', 'Peter Push', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/a13361b3397245a4fcd3792baf9b8080d15ec044d134046305a3f900c98ac229?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/a13361b3397245a4fcd3792baf9b8080d15ec044d134046305a3f900c98ac229?s=80&d=identicon'
), (
  8, 'rebaserandy', 'rebase@randy.de', 'Rebase Randy', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/300c3be6ed362cac33850d5dd5e1c815419144f2c4798e32fa40fd0141c229bc?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/300c3be6ed362cac33850d5dd5e1c815419144f2c4798e32fa40fd0141c229bc?s=80&d=identicon'
), (
  9, 'stashalot', 'sirstash@alot.de', 'Sir Stash Alot', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/30f8e8c1a908ec86cabbfe34409b3a844e95e991eb2b0865c3af7be7d780c9fb?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/30f8e8c1a908ec86cabbfe34409b3a844e95e991eb2b0865c3af7be7d780c9fb?s=80&d=identicon'
), (
  10, 'ingeissue', 'inge@issue.de', 'Inge Issue', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/5702878b262d980c133639bede83866fed57dfd84e27af0265ec031807802c1e?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/5702878b262d980c133639bede83866fed57dfd84e27af0265ec031807802c1e?s=80&d=identicon'
), (
  11, 'franzifeature', 'franzi@feature.de', 'Franzi Feature', NOW(), NOW(),
  'https://secure.gravatar.com/avatar/9457a2bb061c4677314535ccba8a39a0343da4d851cb21aae335d9911a20fb58?s=80&d=identicon',
  'https://secure.gravatar.com/avatar/9457a2bb061c4677314535ccba8a39a0343da4d851cb21aae335d9911a20fb58?s=80&d=identicon'
);

INSERT INTO "public"."classrooms" (
  "id",
  "name",
  "description",
  "owner_id",
  "create_teams",
  "max_team_size",
  "max_teams",
  "group_id",
  "group_access_token_id",
  "group_access_token",
  "group_access_token_created_at",
  "students_view_all_projects",
  "invite_code",
  "created_at",
  "updated_at"
) VALUES
( 'a8566e13-f696-494a-87f9-98d40b178291', 'Classroom with Teams', 'Classroom with Teams!', 2, true, 2, 0, -1, -1, '', '1970-01-01 00:00:00.000000+00', false, '253cdd84-9414-42a6-ada2-2acb74f75cdd', NOW(), NOW() ),
( '827469b8-0dd1-431c-8b00-031b00547dc7', 'Classroom without Teams', 'Classroom without Teams!', 2, false, 1, 0, -1, -1, '', '1970-01-01 00:00:00.000000+00', false, '90913851-cd98-4c63-a767-bba7eb63d47e', NOW(), NOW() );


INSERT INTO "public"."teams" (
  "id", "name", "group_id", "classroom_id", "created_at", "updated_at"
) VALUES
('ddb55138-4dc8-4c6e-9fcf-76f162920841', 'Push Pandas', -1, 'a8566e13-f696-494a-87f9-98d40b178291', NOW(), NOW()),
('164eb7ad-b631-4871-b2fc-15085285bcc4', 'Deploy Dackel', -1, 'a8566e13-f696-494a-87f9-98d40b178291', NOW(), NOW()),
('a741b295-fd74-4622-b763-19cc0eddc847', 'Hotfix Hamster', -1, 'a8566e13-f696-494a-87f9-98d40b178291', NOW(), NOW()),
('45a910fa-478c-4187-8012-cc3567682592', 'Feature Füchse', -1, 'a8566e13-f696-494a-87f9-98d40b178291', NOW(), NOW()),

('9e50f780-d74b-4de9-b516-2fae56c37794', 'Rebase Randy', -1, '827469b8-0dd1-431c-8b00-031b00547dc7', NOW(), NOW()),
('2a6506f0-e43a-4d42-a195-309f7adc7f82', 'Sir Stash Alot', -1, '827469b8-0dd1-431c-8b00-031b00547dc7', NOW(), NOW()),
('f7abf4fb-9937-48ed-85be-18dc19f8c2e7', 'Inge Issue', -1, '827469b8-0dd1-431c-8b00-031b00547dc7', NOW(), NOW());

INSERT INTO "public"."user_classrooms" (
  "user_id", "classroom_id", "team_id", "role", "created_at", "updated_at"
) VALUES
( 2, 'a8566e13-f696-494a-87f9-98d40b178291', null, 0, NOW(), NOW() ),
( 3, 'a8566e13-f696-494a-87f9-98d40b178291', null, 1, NOW(), NOW() ),
( 4, 'a8566e13-f696-494a-87f9-98d40b178291', null, 2, NOW(), NOW() ),
( 5, 'a8566e13-f696-494a-87f9-98d40b178291', 'a741b295-fd74-4622-b763-19cc0eddc847', 2, NOW(), NOW() ),
( 6, 'a8566e13-f696-494a-87f9-98d40b178291', 'a741b295-fd74-4622-b763-19cc0eddc847', 2, NOW(), NOW() ),
( 7, 'a8566e13-f696-494a-87f9-98d40b178291', 'ddb55138-4dc8-4c6e-9fcf-76f162920841', 2, NOW(), NOW() ),
( 8, 'a8566e13-f696-494a-87f9-98d40b178291', null, 2, NOW(), NOW() ),
( 9, 'a8566e13-f696-494a-87f9-98d40b178291', 'ddb55138-4dc8-4c6e-9fcf-76f162920841', 2, NOW(), NOW() ),
( 10, 'a8566e13-f696-494a-87f9-98d40b178291', null, 2, NOW(), NOW() ),
( 11, 'a8566e13-f696-494a-87f9-98d40b178291', '45a910fa-478c-4187-8012-cc3567682592', 2, NOW(), NOW() ),

( 2, '827469b8-0dd1-431c-8b00-031b00547dc7', null, 0, NOW(), NOW() ),
( 3, '827469b8-0dd1-431c-8b00-031b00547dc7', null, 1, NOW(), NOW() ),
( 8, '827469b8-0dd1-431c-8b00-031b00547dc7', '9e50f780-d74b-4de9-b516-2fae56c37794', 2, NOW(), NOW() ),
( 9, '827469b8-0dd1-431c-8b00-031b00547dc7', '2a6506f0-e43a-4d42-a195-309f7adc7f82', 2, NOW(), NOW() ),
( 10, '827469b8-0dd1-431c-8b00-031b00547dc7', 'f7abf4fb-9937-48ed-85be-18dc19f8c2e7', 2, NOW(), NOW() );

  -- 2, "tonitester", "toni@tester.de", "Toni Tester",
  -- 3, "lauratester", "laura@tester.de", "Laura Tester",
  -- 4, "harri", "harald@hueftschmerz.de", "Harald Hüftschmerz",
  -- 5, "mustermanfred", "manfred@muster.de", "Manfred Muster",
  -- 6, "hanshotfix", "hans@hotfix.de", "Hans Hotfix",
  -- 7, "peterpush", "peter@push.de", "Peter Push",
  -- 8, "rebaserandy", "rebase@randy.de", "Rebase Randy",
  -- 9, "stashalot", "sirstash@alot.de", "Sir Stash Alot",
  -- 10, "ingeissue", "inge@issue.de", "Inge Issue",
  -- 11, "franzifeature", "franzi@feature.de", "Franzi Feature",

INSERT INTO "public"."classroom_invitations" (
    "id",
    "status",
    "classroom_id",
    "email",
    "expiry_date",
    "created_at",
    "updated_at"
) VALUES
( 'bb87e0cc-b958-4ada-90b7-0748a82588e9', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'toni@tester.de', NOW() + '14 days', NOW(), NOW() ),
( 'f97c36b7-6e42-45a3-9b0e-96bdd3c2c9fa', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'laura@tester.de', NOW() + '14 days', NOW(), NOW() ),
( 'd97b471d-f343-44a3-8838-b5de33d3896b', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'harald@hueftschmerz.de', NOW() + '14 days', NOW(), NOW() ),
( '4d69f23b-cee0-4cff-b8d9-a1674a15192b', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'manfred@muster.de', NOW() + '14 days', NOW(), NOW() ),
( '1674cd23-002d-4f79-b4a9-b93dc5830648', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'hans@hotfix.de', NOW() + '14 days', NOW(), NOW() ),
( '1bbca05e-c3ce-4854-85cb-18ac5c20465f', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'peter@push.de', NOW() + '14 days', NOW(), NOW() ),
( 'dbcf9ae8-eb0b-4786-82c9-93b48fcf3702', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'rebase@randy.de', NOW() + '14 days', NOW(), NOW() ),
( 'e0a41fc6-f197-4437-b6b0-0c09aedd37f1', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'sirstash@alot.de', NOW() + '14 days', NOW(), NOW() ),
( '8d711f63-301e-47b7-b300-f0e0439a2427', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'inge@issue.de', NOW() + '14 days', NOW(), NOW() ),
( 'd0d5f9cd-ec19-40cc-9ee0-43f812735f7b', 1 , 'a8566e13-f696-494a-87f9-98d40b178291', 'franzi@feature.de', NOW() + '14 days', NOW(), NOW() ),

( 'd55bb0a7-6309-477b-b433-862fd1b9ac17', 1, '827469b8-0dd1-431c-8b00-031b00547dc7', 'toni@tester.de', NOW() + '14 days', NOW(), NOW() ),
( 'ff5fff6e-5a15-45cc-b6e0-b8dccee4107d', 1, '827469b8-0dd1-431c-8b00-031b00547dc7', 'laura@tester.de', NOW() + '14 days', NOW(), NOW() ),
( '3d81f8bc-1949-46e1-8d6e-a766f192fbd3', 1, '827469b8-0dd1-431c-8b00-031b00547dc7', 'rebase@randy.de', NOW() + '14 days', NOW(), NOW() ),
( '889a6fff-adb0-40e9-8b7b-9f1d67abc591', 1, '827469b8-0dd1-431c-8b00-031b00547dc7', 'sirstash@alot.de', NOW() + '14 days', NOW(), NOW() ),
( '109574a5-f7d1-46ad-a9ad-4953d518639b', 1, '827469b8-0dd1-431c-8b00-031b00547dc7', 'inge@issue.de', NOW() + '14 days', NOW(), NOW() ),
( '9054d699-ef36-44b5-8677-152f20b47b5f', 1, '827469b8-0dd1-431c-8b00-031b00547dc7', 'inge@issue.de', NOW() + '14 days', NOW(), NOW() );

INSERT INTO "public"."assignments" (
  "id",
  "classroom_id",
  "template_project_id",
  "name",
  "description",
  "due_date",
  "closed",
  "grading_j_unit_auto_grading_active",
  "acceptable_since",
  "created_at",
  "updated_at"
) VALUES
( '1ea03cc3-645f-4591-8c6c-a4742a1d91fb', 'a8566e13-f696-494a-87f9-98d40b178291', -1, 'First Go Assignment', 'Test description', NOW() + '14 days', false, false, NOW(), NOW(), NOW() ),
( 'b78ac40a-ba52-49c9-b596-a363bf0e14ea', 'a8566e13-f696-494a-87f9-98d40b178291', -2, 'First Simple Assignment', 'Test description', NOW() + '14 days', false, false, null, NOW(), NOW() ),

( '979efba7-fe05-42d1-997e-f00250b11a44', '827469b8-0dd1-431c-8b00-031b00547dc7', -1, 'First Simple Assignment', 'Test description', NOW() + '14 days', false, false, NOW(), NOW(), NOW() ),
( '22715641-c93f-492d-8101-77bb3ee1d94a', '827469b8-0dd1-431c-8b00-031b00547dc7', -2, 'Second Go Assignment', 'Test description', NOW() + '14 days', false, false, null, NOW(), NOW() );

INSERT INTO "public"."assignment_projects"(
  "id",
  "team_id",
  "assignment_id",
  "created_at",
  "updated_at"
) VALUES
('1cac3ab2-b11e-4230-85bd-4b32aec9888b', 'ddb55138-4dc8-4c6e-9fcf-76f162920841', '1ea03cc3-645f-4591-8c6c-a4742a1d91fb', NOW(), NOW()),
('e5421a39-da2c-4fe4-b4da-2577345dc0d5', '164eb7ad-b631-4871-b2fc-15085285bcc4', '1ea03cc3-645f-4591-8c6c-a4742a1d91fb', NOW(), NOW()),
('351d4d4f-ba83-4671-ae85-65755aaa1f5e', 'a741b295-fd74-4622-b763-19cc0eddc847', '1ea03cc3-645f-4591-8c6c-a4742a1d91fb', NOW(), NOW()),
('97616c03-d11c-4e7e-9fd4-677aed0d02d4', '45a910fa-478c-4187-8012-cc3567682592', '1ea03cc3-645f-4591-8c6c-a4742a1d91fb', NOW(), NOW()),

('1befbc6d-f210-4070-b8b3-cbce9e3365bf', '9e50f780-d74b-4de9-b516-2fae56c37794', '979efba7-fe05-42d1-997e-f00250b11a44', NOW(), NOW()),
('1ac6a69b-080b-479c-a37a-5c9670967f38', '2a6506f0-e43a-4d42-a195-309f7adc7f82', '979efba7-fe05-42d1-997e-f00250b11a44', NOW(), NOW()),
('6e493847-c7e6-4376-8bdf-53ec9501dc2f', 'f7abf4fb-9937-48ed-85be-18dc19f8c2e7', '979efba7-fe05-42d1-997e-f00250b11a44', NOW(), NOW());

-- +goose Down
DELETE FROM "public"."assignment_projects" WHERE true;
DELETE FROM "public"."assignments" WHERE true;
DELETE FROM "public"."classroom_invitations" WHERE true;
DELETE FROM "public"."user_classrooms" WHERE true;
DELETE FROM "public"."teams" WHERE true;
DELETE FROM "public"."classrooms" WHERE true;
DELETE FROM "public"."users" WHERE true;
