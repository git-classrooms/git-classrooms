package worker

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gitlabConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/gitlab"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	gitlabRepoMock "gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/_mock"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
)

// TestSyncClassroomsWork tests the main syncing functionalities of SyncGitlabDbWork.
func TestSyncClassroomsWork(t *testing.T) {
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	repo := gitlabRepoMock.NewMockRepository(t)

	pg, err := tests.StartPostgres()
	if err != nil {
		t.Fatalf("could not start database container: %s", err.Error())
	}

	dbUrl, err := pg.ConnectionString(context.Background())
	if err != nil {
		t.Fatalf("could not get database connection string: %s", err.Error())
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not connect to database: %s", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("could not get database connection: %s", err.Error())
	}

	// 1. Migrate database
	err = database.MigrateDatabase(sqlDB)
	if err != nil {
		t.Fatalf("could not migrate database: %s", err.Error())
	}

	query.SetDefault(db)

	// Setup test data: users, classroom, teams, assignments
	owner := factory.User()
	member1 := factory.User()
	member2 := factory.User()
	classroom1 := factory.Classroom(owner.ID)

	members := []*database.UserClassrooms{
		factory.UserClassroom(owner.ID, classroom1.ID, database.Owner),
		factory.UserClassroom(member1.ID, classroom1.ID, database.Student),
		factory.UserClassroom(member2.ID, classroom1.ID, database.Student),
	}

	team1 := factory.Team(classroom1.ID, members)
	dueDate := time.Now().Add(1 * time.Hour)
	assignment1 := factory.Assignment(classroom1.ID, &dueDate, false)
	assignment1Project := factory.AssignmentProject(assignment1.ID, team1.ID)

	publicUrl := &url.URL{Scheme: "http", Host: "localhost"}

	w := NewSyncGitlabDbWork(&gitlabConfig.GitlabConfig{}, publicUrl)

	// Test the getUnarchivedClassrooms method.
	t.Run("getUnarchivedClassrooms", func(t *testing.T) {
		classroom2 := factory.Classroom(owner.ID)
		factory.UserClassroom(owner.ID, classroom2.ID, database.Owner)

		classroom2.Archived = true

		query.Classroom.WithContext(context.Background()).Save(classroom2)

		dbClassrooms := w.getUnarchivedClassrooms(context.Background())

		if len(dbClassrooms) != 1 {
			t.Errorf("Expected 1 classroom, got %d", len(dbClassrooms))
		}
		assert.Equal(t, classroom1.ID, dbClassrooms[0].ID)
		assert.Len(t, dbClassrooms[0].Member, 3)
		assert.Len(t, dbClassrooms[0].Teams, 1)
	})

	// Test the syncClassroom method.
	t.Run("syncClassroom", func(t *testing.T) {
		newName := "new name"
		newDescription := "new description"

		repo.EXPECT().
			GetGroupById(classroom1.GroupID).
			Return(&model.Group{
				Name:        newName,
				Description: newDescription,
			}, nil).
			Times(1)

		repo.EXPECT().
			ChangeGroupName(classroom1.GroupID, classroom1.Name).
			Return(nil, nil).
			Times(1)

		repo.EXPECT().
			ChangeGroupDescription(classroom1.GroupID, utils.CreateClassroomGitlabDescription(classroom1, publicUrl)).
			Return(nil, nil).
			Times(1)

		w.syncClassroom(context.Background(), *classroom1, repo)

		repo.AssertExpectations(t)

		dbClassroom1, err := query.Classroom.WithContext(context.Background()).
			Where(query.Classroom.ID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, classroom1.Name, dbClassroom1.Name)
		assert.Equal(t, classroom1.Description, dbClassroom1.Description)

		// Revert changes of db object for the next tests
		dbClassroom1.Name = classroom1.Name
		dbClassroom1.Description = classroom1.Description
		query.UserClassrooms.WithContext(context.Background()).Updates(dbClassroom1)
	})

	// Test syncClassroomMember method: handle case when members have left via GitLab.
	t.Run("syncClassroomMember - left via gitlab", func(t *testing.T) {
		repo.EXPECT().
			GetAllUsersOfGroup(classroom1.GroupID).
			Return([]*model.User{
				{
					ID:       owner.ID,
					Username: owner.GitlabUsername,
					Email:    owner.GitlabEmail,
				},
				{
					ID:       member1.ID,
					Username: member1.GitlabUsername,
					Email:    member1.GitlabEmail,
				},
			}, nil).
			Times(1)

		classroom, err := query.Classroom.WithContext(context.Background()).
			Preload(query.Classroom.Member).
			Where(query.Classroom.ID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)

		w.syncClassroomMember(context.Background(), classroom1.GroupID, classroom.Member, repo)

		repo.AssertExpectations(t)

		leftMember, err := query.UserClassrooms.WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(member2.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(classroom1.ID)).
			First()
		assert.Error(t, err)
		assert.Nil(t, leftMember)

		// Revert changes for the next tests
		factory.UserClassroom(member2.ID, classroom1.ID, database.Student)
	})

	// Test syncClassroomMember method: handle case when members are added via GitLab.
	t.Run("syncClassroomMember - added via gitlab", func(t *testing.T) {
		repo.EXPECT().
			GetAllUsersOfGroup(classroom1.GroupID).
			Return([]*model.User{
				{
					ID:       owner.ID,
					Username: owner.GitlabUsername,
					Email:    owner.GitlabEmail,
				},
				{
					ID:       member1.ID,
					Username: member1.GitlabUsername,
					Email:    member1.GitlabEmail,
				},
				{
					ID:       member2.ID,
					Username: member2.GitlabUsername,
					Email:    member2.GitlabEmail,
				},
				{
					ID:       4,
					Username: "new",
					Email:    "new",
				},
			}, nil).
			Times(1)

		repo.EXPECT().
			RemoveUserFromGroup(classroom1.GroupID, 4).
			Return(nil).
			Times(1)

		classroom, err := query.Classroom.WithContext(context.Background()).
			Preload(query.Classroom.Member).
			Where(query.Classroom.ID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)

		w.syncClassroomMember(context.Background(), classroom1.GroupID, classroom.Member, repo)

		repo.AssertExpectations(t)
	})

	// Test syncTeam method.
	t.Run("syncTeam", func(t *testing.T) {
		newName := "new name"
		newDescription := "new description"
		repo.EXPECT().
			GetGroupById(team1.GroupID).
			Return(&model.Group{
				Name:        newName,
				Description: newDescription,
			}, nil).
			Times(1)

		repo.EXPECT().
			ChangeGroupName(team1.GroupID, team1.Name).
			Return(nil, nil).
			Times(1)

		repo.EXPECT().
			ChangeGroupDescription(team1.GroupID, utils.CreateTeamGitlabDescription(classroom1, team1, publicUrl)).
			Return(nil, nil).
			Times(1)

		w.syncTeam(context.Background(), classroom1, *team1, repo)

		repo.AssertExpectations(t)

		dbTeam1, err := query.Team.WithContext(context.Background()).
			Where(query.Team.ID.Eq(team1.ID)).
			First()
		assert.NoError(t, err)
		assert.Equal(t, team1.Name, dbTeam1.Name)
	})

	// Test syncTeamMember method: handle case when members have left the team via GitLab.
	t.Run("syncTeamMember - left via gitlab", func(t *testing.T) {
		repo.EXPECT().
			GetAllUsersOfGroup(team1.GroupID).
			Return([]*model.User{
				{
					ID:       owner.ID,
					Username: owner.GitlabUsername,
					Email:    owner.GitlabEmail,
				},
				{
					ID:       member1.ID,
					Username: member1.GitlabUsername,
					Email:    member1.GitlabEmail,
				},
			}, nil).
			Times(1)

		team, err := query.Team.WithContext(context.Background()).
			Preload(query.Team.Member).
			Where(query.Team.ID.Eq(team1.ID)).
			First()
		assert.NoError(t, err)

		w.syncTeamMember(context.Background(), team1.GroupID, team.Member, repo)

		repo.AssertExpectations(t)

		leftMember, err := query.UserClassrooms.WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(member2.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		assert.Nil(t, leftMember.TeamID)
		assert.Nil(t, leftMember.Team)

		// Revert changes to the db object for the next test
		leftMember, err = query.UserClassrooms.WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(member2.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(classroom1.ID)).
			First()
		assert.NoError(t, err)
		leftMember.TeamID = &team1.ID
		leftMember.Team = team1
		query.UserClassrooms.WithContext(context.Background()).Updates(leftMember)
	})

	// Test syncTeamMember method: handle case when members are added to the team via GitLab.
	t.Run("syncTeamMember - added via gitlab", func(t *testing.T) {
		repo.EXPECT().
			GetAllUsersOfGroup(team1.GroupID).
			Return([]*model.User{
				{
					ID:       owner.ID,
					Username: owner.GitlabUsername,
					Email:    owner.GitlabEmail,
				},
				{
					ID:       member1.ID,
					Username: member1.GitlabUsername,
					Email:    member1.GitlabEmail,
				},
				{
					ID:       member2.ID,
					Username: member2.GitlabUsername,
					Email:    member2.GitlabEmail,
				},
				{
					ID:       4,
					Username: "new",
					Email:    "new",
				},
			}, nil).
			Times(1)

		repo.EXPECT().
			RemoveUserFromGroup(team1.GroupID, 4).
			Return(nil).
			Times(1)

		team, err := query.Team.WithContext(context.Background()).
			Preload(query.Team.Member).
			Where(query.Team.ID.Eq(team1.ID)).
			First()
		assert.NoError(t, err)

		w.syncTeamMember(context.Background(), team1.GroupID, team.Member, repo)

		repo.AssertExpectations(t)
	})

	// Test getAssignmentProjects method.
	t.Run("getAssignmentProjects", func(t *testing.T) {
		dbAssignments := w.getAssignmentProjects(context.Background(), assignment1.ID)

		if len(dbAssignments) != 1 {
			t.Errorf("Expected 1 assignment, got %d", len(dbAssignments))
		}
		assert.Equal(t, assignment1Project.ProjectID, dbAssignments[0].ProjectID)
	})

	// Test syncProject method.
	t.Run("syncProject", func(t *testing.T) {
		repo.EXPECT().
			GetProjectById(assignment1Project.ProjectID).
			Return(nil, fiber.NewError(404, "404 {message: 404 Project Not Found}")).
			Times(1)

		w.syncProject(context.Background(), *assignment1Project, repo)

		repo.AssertExpectations(t)

		deletedProject, err := query.AssignmentProjects.WithContext(context.Background()).
			Where(query.AssignmentProjects.ID.Eq(assignment1Project.ID)).
			First()
		assert.Error(t, err)
		assert.Nil(t, deletedProject)
	})
}
