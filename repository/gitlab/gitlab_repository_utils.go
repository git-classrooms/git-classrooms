package gitlab

import (
	"time"

	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"

	goGitlab "github.com/xanzy/go-gitlab"
)

func ProjectFromGoGitlab(gitlabProject goGitlab.Project) *model.Project {
	var owner *model.User = nil
	if gitlabProject.Owner != nil {
		owner = UserFromGoGitlab(*gitlabProject.Owner)
	}

	return &model.Project{
		Name:          gitlabProject.Name,
		ID:            gitlabProject.ID,
		Visibility:    VisibilityFromGoGitlab(gitlabProject.Visibility),
		WebUrl:        gitlabProject.WebURL,
		Description:   gitlabProject.Description,
		Owner:         owner,
		DefaultBranch: gitlabProject.DefaultBranch,
	}
}

func ProjectFromGoGitlabWithProjectMembers(gitlabProject goGitlab.Project, gitlabMembers []*goGitlab.ProjectMember) *model.Project {
	var owner *model.User = nil
	if gitlabProject.Owner != nil {
		owner = UserFromGoGitlab(*gitlabProject.Owner)
	}

	members := make([]model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		members[i] = *UserFromGoGitlabProjectMember(*gitlabMember)
	}

	return &model.Project{
		Name:        gitlabProject.Name,
		ID:          gitlabProject.ID,
		Visibility:  VisibilityFromGoGitlab(gitlabProject.Visibility),
		WebUrl:      gitlabProject.WebURL,
		Description: gitlabProject.Description,
		Owner:       owner,
		Members:     members,
	}
}

func ProjectFromGoGitlabWithGroupMembers(gitlabProject goGitlab.Project, gitlabMembers []*goGitlab.GroupMember) *model.Project {
	var owner *model.User = nil
	if gitlabProject.Owner != nil {
		owner = UserFromGoGitlab(*gitlabProject.Owner)
	}

	members := make([]model.User, len(gitlabMembers))
	for i, gitlabMember := range gitlabMembers {
		members[i] = *UserFromGoGitlabGroupMember(*gitlabMember)
	}

	return &model.Project{
		Name:        gitlabProject.Name,
		ID:          gitlabProject.ID,
		Visibility:  VisibilityFromGoGitlab(gitlabProject.Visibility),
		WebUrl:      gitlabProject.WebURL,
		Description: gitlabProject.Description,
		Owner:       owner,
		Members:     members,
	}
}

func BranchFromGoGitlab(input *goGitlab.Branch) *model.Branch {
	return &model.Branch{
		Name:      input.Name,
		Protected: input.Protected,
		Default:   input.Default,
		WebURL:    input.WebURL,
	}
}

func VisibilityFromGoGitlab(input goGitlab.VisibilityValue) model.Visibility {
	if input == "public" {
		return model.Public
	} else if input == "internal" {
		return model.Internal
	} else if input == "private" {
		return model.Private
	}
	return 0
}

func VisibilityFromModel(input model.Visibility) goGitlab.VisibilityValue {
	switch input {
	case model.Public:
		return goGitlab.PublicVisibility
	case model.Internal:
		return goGitlab.InternalVisibility
	case model.Private:
		return goGitlab.PrivateVisibility
	default:
		return goGitlab.PrivateVisibility
	}
}

func AccessLevelFromGoGitlab(input goGitlab.AccessLevelValue) model.AccessLevelValue {
	switch input {
	case goGitlab.NoPermissions:
		return model.NoPermissions
	case goGitlab.MinimalAccessPermissions:
		return model.MinimalAccessPermissions
	case goGitlab.GuestPermissions:
		return model.GuestPermissions
	case goGitlab.ReporterPermissions:
		return model.ReporterPermissions
	case goGitlab.DeveloperPermissions:
		return model.DeveloperPermissions
	case goGitlab.MaintainerPermissions:
		return model.MaintainerPermissions
	case goGitlab.OwnerPermissions:
		return model.OwnerPermissions
	case goGitlab.AdminPermissions:
		return model.AdminPermissions
	default:
		return model.NoPermissions // Default case
	}
}

func AccessLevelFromModel(input model.AccessLevelValue) goGitlab.AccessLevelValue {
	switch input {
	case model.NoPermissions:
		return goGitlab.NoPermissions
	case model.MinimalAccessPermissions:
		return goGitlab.MinimalAccessPermissions
	case model.GuestPermissions:
		return goGitlab.GuestPermissions
	case model.ReporterPermissions:
		return goGitlab.ReporterPermissions
	case model.DeveloperPermissions:
		return goGitlab.DeveloperPermissions
	case model.MaintainerPermissions:
		return goGitlab.MaintainerPermissions
	case model.OwnerPermissions:
		return goGitlab.OwnerPermissions
	case model.AdminPermissions:
		return goGitlab.AdminPermissions
	default:
		return goGitlab.NoPermissions // Default case
	}
}

func UserFromGoGitlab(input goGitlab.User) *model.User {
	return &model.User{
		ID:       input.ID,
		Username: input.Username,
		Name:     input.Name,
		WebUrl:   input.WebURL,
		Email:    input.Email,
		Avatar:   model.UserAvatar{AvatarURL: &input.AvatarURL},
	}
}

func UserFromGoGitlabProjectMember(input goGitlab.ProjectMember) *model.User {
	return &model.User{
		ID:       input.ID,
		Username: input.Username,
		Name:     input.Name,
		WebUrl:   input.WebURL,
		Email:    input.Email,
	}
}

func UserFromGoGitlabGroupMember(input goGitlab.GroupMember) *model.User {
	return &model.User{
		ID:       input.ID,
		Username: input.Username,
		Name:     input.Name,
		WebUrl:   input.WebURL,
		Email:    input.Email,
	}
}

func GroupFromGoGitlab(input goGitlab.Group) *model.Group {
	return &model.Group{
		Name:        input.Name,
		ID:          input.ID,
		Description: input.Description,
		WebUrl:      input.WebURL,
		Visibility:  VisibilityFromGoGitlab(input.Visibility),
	}
}

func GroupFromGoGitlabWithMembersAndProjects(group goGitlab.Group, members []model.User, projects []model.Project) *model.Group {

	return &model.Group{
		Name:        group.Name,
		ID:          group.ID,
		Description: group.Description,
		WebUrl:      group.WebURL,
		Visibility:  VisibilityFromGoGitlab(group.Visibility),
		Projects:    projects,
		Member:      members,
	}
}

func TestReportFromGoGitlabTestReport(testReport *goGitlab.PipelineTestReport) *model.TestReport {
	var report model.TestReport
	report.TotalTime = testReport.TotalTime
	report.TotalCount = testReport.TotalCount
	report.SuccessCount = testReport.SuccessCount
	report.FailedCount = testReport.FailedCount
	report.SkippedCount = testReport.SkippedCount
	report.ErrorCount = testReport.ErrorCount

	report.TestSuites = make([]model.TestReportTestSuite, len(testReport.TestSuites))
	for i, ts := range testReport.TestSuites {
		var suite model.TestReportTestSuite
		suite.Name = ts.Name
		suite.TotalTime = ts.TotalTime
		suite.TotalCount = ts.TotalCount
		suite.SuccessCount = ts.SuccessCount
		suite.FailedCount = ts.FailedCount
		suite.SkippedCount = ts.SkippedCount
		suite.ErrorCount = ts.ErrorCount

		suite.TestCases = make([]model.TestReportTestCase, len(ts.TestCases))
		for j, tc := range ts.TestCases {
			var case_ model.TestReportTestCase
			case_.Status = tc.Status
			case_.Name = tc.Name
			case_.Classname = tc.Classname
			case_.ExecutionTime = tc.ExecutionTime
			case_.SystemOutput = tc.SystemOutput
			case_.StackTrace = tc.StackTrace

			suite.TestCases[j] = case_
		}

		report.TestSuites[i] = suite
	}

	return &report
}

func PendingInviteFromGoGitlab(input goGitlab.PendingInvite) *model.PendingInvite {
	return &model.PendingInvite{
		ID:            input.ID,
		InviteEmail:   input.InviteEmail,
		CreatedAt:     input.CreatedAt,
		AccessLevel:   AccessLevelFromGoGitlab(input.AccessLevel),
		ExpiresAt:     input.ExpiresAt,
		UserName:      input.UserName,
		CreatedByName: input.CreatedByName,
	}
}

func ConvertUserPointerSlice(input []*model.User) []model.User {
	output := make([]model.User, len(input))
	for i, ptr := range input {
		output[i] = *ptr
	}
	return output
}

func ConvertProjectPointerSlice(input []*model.Project) []model.Project {
	output := make([]model.Project, len(input))
	for i, ptr := range input {
		output[i] = *ptr
	}
	return output
}

func ConvertPendingInvitePointerSlice(input []*model.PendingInvite) []model.PendingInvite {
	output := make([]model.PendingInvite, len(input))
	for i, ptr := range input {
		output[i] = *ptr
	}
	return output
}

func GroupAccessTokenFromGoGitlabGroupAccessToken(input goGitlab.GroupAccessToken) *model.GroupAccessToken {
	return &model.GroupAccessToken{
		ID:          input.ID,
		UserID:      input.UserID,
		Name:        input.Name,
		Scopes:      input.Scopes,
		ExpiresAt:   time.Time(*input.ExpiresAt),
		Token:       input.Token,
		AccessLevel: AccessLevelFromGoGitlab(input.AccessLevel),
	}
}
