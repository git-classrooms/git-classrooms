package main

import (
	"context"

	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gorm.io/gen/field"
)

func GetAllUsers(ctx context.Context) ([]*database.User, error) {
	queryUser := query.User

	users, err := queryUser.WithContext(ctx).
		Find()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetAllClassrooms(ctx context.Context) ([]*database.Classroom, error) {
	queryClassroom := query.Classroom

	classrooms, err := queryClassroom.WithContext(ctx).
		Preload(queryClassroom.Member).
		Preload(queryClassroom.Member.User).
		Preload(queryClassroom.Invitations).
		Preload(queryClassroom.Assignments).
		Preload(queryClassroom.Teams).
		Preload(field.NewRelation("Teams.Member", "")).
		Preload(field.NewRelation("Teams.Member.User", "")).
		Find()
	if err != nil {
		return nil, err
	}

	return classrooms, nil
}

func SaveClassroom(ctx context.Context, classroom *database.Classroom) error {
	return query.Classroom.WithContext(ctx).Save(classroom)
}
