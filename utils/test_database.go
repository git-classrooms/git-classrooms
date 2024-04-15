package utils

import (
	"context"
	"testing"

	databaseConfig "gitlab.hs-flensburg.de/gitlab-classroom/config/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/tests"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestDB struct {
	db *gorm.DB
	t  *testing.T
}

func NewTestDB(t *testing.T) *TestDB {
	db := TestDB{t: t}

	db.Setup()
	query.SetDefault(db.db)

	return &db
}

func (testDb *TestDB) Setup() {
	testDb.t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "false")
	pq, err := tests.StartPostgres()
	if err != nil {
		testDb.t.Fatalf("could not start database container: %s", err.Error())
	}
	testDb.t.Cleanup(func() {
		pq.Terminate(context.Background())
	})
	port, err := pq.MappedPort(context.Background(), "5432")
	if err != nil {
		testDb.t.Fatalf("could not get database container port: %s", err.Error())
	}
	dbConfig := databaseConfig.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	}
	testDb.db, err = gorm.Open(postgres.Open(dbConfig.Dsn()), &gorm.Config{})
	if err != nil {
		testDb.t.Fatalf("could not connect to database: %s", err.Error())
	}
	err = MigrateDatabase(testDb.db)
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
