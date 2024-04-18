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
	session.InitSessionStore(db.dbUrl)

	return &db
}

func (testDb *TestDB) Setup() {
	testDb.t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "false")
	pq, err := StartPostgres()
	if err != nil {
		testDb.t.Fatalf("could not start database container: %s", err.Error())
	}
	testDb.t.Cleanup(func() {
		pq.Terminate(context.Background())
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
