package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/unduu/e-learning/learning/model"
)

type LearningRepository struct {
	conn *sqlx.DB
}

func NewLearningRepository(db *sqlx.DB) *LearningRepository {
	return &LearningRepository{
		conn: db,
	}
}

// GetCourses return all course list
func (a *LearningRepository) GetCourses() []*model.Course {
	// DB Response struct
	courses := make([]*model.Course, 0)

	// Data for query
	queryParams := map[string]interface{}{}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT id,alias,title,subtitle FROM courses`)
	if err != nil {
		fmt.Println("Error db GetCourses->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&courses, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	return courses
}

// GetParticipantByCourse return user who join the course
func (a *LearningRepository) GetParticipantByCourse(id int) []*model.Participant {
	// DB Response struct
	participants := make([]*model.Participant, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"course_id": id,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT username,status FROM course_participants WHERE course_id = :course_id`)
	if err != nil {
		fmt.Println("Error db GetParticipantByCourse->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&participants, queryParams)
	if err != nil {
		fmt.Println("Error db GetParticipantByCourse->query.Get : ", err)
	}

	return participants
}
