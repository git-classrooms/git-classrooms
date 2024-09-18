package api

import (
	"context"
	"fmt"
	"testing"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils/factory"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPatchMemberRole(t *testing.T) {
	restoreDatabase(t)

	creator := factory.User()
	classroom := factory.Classroom(creator.ID)
	factory.UserClassroom(creator.ID, classroom.ID, database.Owner)

	user := factory.User()
	userClassroom := factory.UserClassroom(user.ID, classroom.ID, database.Student)

	t.Run("Owner > Student by creator, viewOtherProjects false", func(t *testing.T) {
		updateUserClassroom(t, userClassroom, database.Owner)
		updateClassroom(t, classroom, false)

		r := RunTest{
			loginUser:      creator,
			user:           user,
			changeTo:       database.Student,
			accessLevel:    model.GuestPermissions,
			expectedStatus: fiber.StatusAccepted,
			classroom:      classroom,
		}
		r.run(t)
	})

	t.Run("Owner > Student by creator, viewOtherProjects true", func(t *testing.T) {
		updateUserClassroom(t, userClassroom, database.Owner)
		updateClassroom(t, classroom, true)

		r := RunTest{
			loginUser:      creator,
			user:           user,
			changeTo:       database.Student,
			accessLevel:    model.ReporterPermissions,
			expectedStatus: fiber.StatusAccepted,
			classroom:      classroom,
		}
		r.run(t)
	})

	t.Run("Moderator > Student by creator, viewOtherProjects false", func(t *testing.T) {
		updateUserClassroom(t, userClassroom, database.Moderator)
		updateClassroom(t, classroom, false)

		r := RunTest{
			loginUser:      creator,
			user:           user,
			changeTo:       database.Student,
			accessLevel:    model.GuestPermissions,
			expectedStatus: fiber.StatusAccepted,
			classroom:      classroom,
		}
		r.run(t)
	})

	 // TODO: Adjust case, the permissions should not change
	 //t.Run("Moderator > Student by creator, viewOtherProjects true", func(t *testing.T) {
	 //	updateUserClassroom(t, userClassroom, database.Moderator)
	 //	updateClassroom(t, classroom, true)

	 //	r := RunTest{
	 //		loginUser:      creator,
	 //		user:           user,
	 //		changeTo:       database.Student,
	 //		accessLevel:    model.ReporterPermissions,
	 //		expectedStatus: fiber.StatusAccepted,
	 //		classroom:      classroom,
	 //	}
	 //	r.run(t)
	 //})

	 t.Run("Moderator > Owner by creator, viewOtherProjects false", func(t *testing.T) {
	 	updateUserClassroom(t, userClassroom, database.Moderator)
	 	updateClassroom(t, classroom, false)

	 	r := RunTest{
	 		loginUser:      creator,
	 		user:           user,
	 		changeTo:       database.Owner,
	 		accessLevel:    model.OwnerPermissions,
	 		expectedStatus: fiber.StatusAccepted,
	 		classroom:      classroom,
	 	}
	 	r.run(t)
	 })

	 t.Run("Moderator > Owner by creator, viewOtherProjects true", func(t *testing.T) {
	 	updateUserClassroom(t, userClassroom, database.Moderator)
	 	updateClassroom(t, classroom, true)

	 	r := RunTest{
	 		loginUser:      creator,
	 		user:           user,
	 		changeTo:       database.Owner,
	 		accessLevel:    model.OwnerPermissions,
	 		expectedStatus: fiber.StatusAccepted,
	 		classroom:      classroom,
	 	}
	 	r.run(t)
	 })

	 t.Run("Student > Moderator by creator, viewOtherProjects false", func(t *testing.T) {
	 	updateUserClassroom(t, userClassroom, database.Student)
	 	updateClassroom(t, classroom, false)

	 	r := RunTest{
	 		loginUser:      creator,
	 		user:           user,
	 		changeTo:       database.Moderator,
	 		accessLevel:    model.ReporterPermissions,
	 		expectedStatus: fiber.StatusAccepted,
	 		classroom:      classroom,
	 	}
	 	r.run(t)
	 })

	 // TODO: Adjust case, the permissions should not change
	 //t.Run("Student > Moderator by creator, viewOtherProjects true", func(t *testing.T) {
	 //	updateUserClassroom(t, userClassroom, database.Student)
	 //	updateClassroom(t, classroom, true)

	 //	r := RunTest{
	 //		loginUser:      creator,
	 //		user:           user,
	 //		changeTo:       database.Moderator,
	 //		accessLevel:    model.ReporterPermissions,
	 //		expectedStatus: fiber.StatusAccepted,
	 //		classroom:      classroom,
	 //	}
	 //	r.run(t)
	 //})

	 t.Run("Student > Owner by creator, viewOtherProjects false", func(t *testing.T) {
	 	updateUserClassroom(t, userClassroom, database.Student)
	 	updateClassroom(t, classroom, false)

	 	r := RunTest{
	 		loginUser:      creator,
	 		user:           user,
	 		changeTo:       database.Owner,
	 		accessLevel:    model.OwnerPermissions,
	 		expectedStatus: fiber.StatusAccepted,
	 		classroom:      classroom,
	 	}
	 	r.run(t)
	 })

	 t.Run("Student > Owner by creator, viewOtherProjects true", func(t *testing.T) {
	 	updateUserClassroom(t, userClassroom, database.Student)
	 	updateClassroom(t, classroom, true)

	 	r := RunTest{
	 		loginUser:      creator,
	 		user:           user,
	 		changeTo:       database.Owner,
	 		accessLevel:    model.OwnerPermissions,
	 		expectedStatus: fiber.StatusAccepted,
	 		classroom:      classroom,
	 	}
	 	r.run(t)
	 })
}

func updateClassroom(t *testing.T, classroom *database.Classroom, viewOtherProjects bool) {
	classroom.StudentsViewAllProjects = viewOtherProjects
	err := query.Classroom.WithContext(context.Background()).Save(classroom)

	if err != nil {
		t.Fatal(err)
	}
}

func updateUserClassroom(t *testing.T, userClassrooms *database.UserClassrooms, role database.Role) {
	userClassrooms.Role = role
	err := query.UserClassrooms.WithContext(context.Background()).Save(userClassrooms)

	if err != nil {
		t.Fatal(err)
	}
}

type RunTest struct {
	loginUser      *database.User
	user           *database.User
	changeTo       database.Role
	accessLevel    model.AccessLevelValue
	classroom      *database.Classroom
	expectedStatus int
}

func (r RunTest) run(t *testing.T) {
	// ------------ END OF SEEDING DATA -----------------
	app, gitlabRepo, _ := setupApp(t, r.loginUser)
	route := fmt.Sprintf("/api/v2/classrooms/%s/members/%d/role", r.classroom.ID.String(), r.user.ID)

	requestBody := &updateMemberRoleRequest{
		Role: utils.NewPtr(r.changeTo),
	}

	gitlabRepo.
		EXPECT().
		GroupAccessLogin(r.classroom.GroupAccessToken).
		Return(nil)

	gitlabRepo.
		EXPECT().
		ChangeUserAccessLevelInGroup(r.classroom.GroupID, r.user.ID, r.accessLevel).
		Return(nil)

	req := newJsonRequest(route, requestBody, "PATCH")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, r.expectedStatus, resp.StatusCode)

	updatedUserClassroom, err :=
		query.
			UserClassrooms.
			WithContext(context.Background()).
			Where(query.UserClassrooms.UserID.Eq(r.user.ID)).
			Where(query.UserClassrooms.ClassroomID.Eq(r.classroom.ID)).
			First()

	assert.Equal(t, r.changeTo, updatedUserClassroom.Role)
}
