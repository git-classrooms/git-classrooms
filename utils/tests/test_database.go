package tests

import (
	"context"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestDB struct {
	db    *gorm.DB
	t     *testing.T
	dbUrl string
}

func NewTestDB(t *testing.T) *TestDB {
	db := TestDB{t: t}

	db.Setup()

	query.SetDefault(db.db)
	session.InitSessionStore(nil)

	return &db
}

func (testDb *TestDB) Setup() {
	testDb.t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "false")
	pq, err := StartPostgres()
	if err != nil {
		testDb.t.Fatalf("could not start database container: %s", err.Error())
	}
	testDb.t.Cleanup(func() {
		err = pq.Restore(context.Background())
		if err != nil {
			testDb.t.Fatal(err)
		}
	})

	testDb.dbUrl, err = pq.ConnectionString(context.Background())
	if err != nil {
		testDb.t.Fatalf("could not get database connection string: %s", err.Error())
	}

	testDb.db, err = gorm.Open(postgres.Open(testDb.dbUrl), &gorm.Config{})
	if err != nil {
		testDb.t.Fatalf("could not connect to database: %s", err.Error())
	}

	err = utils.MigrateDatabase(testDb.db)
	if err != nil {
		testDb.t.Fatalf("could not migrate database: %s", err.Error())
	}
}

func (testDb *TestDB) InsertUser(user *database.User) {
	err := query.User.WithContext(context.Background()).Create(user)
	if err != nil {
		testDb.t.Fatalf("could not insert user: %s", err.Error())
	}
}

func (testDb *TestDB) InsertClassroom(classroom *database.Classroom) {
	err := query.Classroom.WithContext(context.Background()).Create(classroom)
	if err != nil {
		testDb.t.Fatalf("could not insert classroom: %s", err.Error())
	}
}

func (testDb *TestDB) InsertAssignment(assignment *database.Assignment) {
	err := query.Assignment.WithContext(context.Background()).Create(assignment)
	if err != nil {
		testDb.t.Fatalf("could not insert assignment: %s", err.Error())
	}
}

func (testDb *TestDB) InsertTeam(team *database.Team) {
	err := query.Team.WithContext(context.Background()).Create(team)
	if err != nil {
		testDb.t.Fatalf("could not insert team: %s", err.Error())
	}
}

func (testDb *TestDB) InsertAssignmentProject(assignmentProject *database.AssignmentProjects) {
	err := query.AssignmentProjects.WithContext(context.Background()).Create(assignmentProject)
	if err != nil {
		testDb.t.Fatalf("could not insert assignment project: %s", err.Error())
	}
}

func (db *TestDB) SaveAssignmentProject(project *database.AssignmentProjects) {
	err := query.AssignmentProjects.WithContext(context.Background()).Save(project)
	if err != nil {
		db.t.Fatalf("could not update assignment project: %s", err.Error())
	}
}

func (testDb *TestDB) InsertInvitation(invitation *database.ClassroomInvitation) {
	err := query.ClassroomInvitation.WithContext(context.Background()).Create(invitation)
	if err != nil {
		testDb.t.Fatalf("could not insert invitation: %s", err.Error())
	}
}

func (testDb *TestDB) SaveInvitation(invitation *database.ClassroomInvitation) {
	err := query.ClassroomInvitation.WithContext(context.Background()).Save(invitation)
	if err != nil {
		testDb.t.Fatalf("could not update invitation: %s", err.Error())
	}
}
