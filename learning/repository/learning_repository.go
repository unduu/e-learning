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
	query, err := a.conn.PrepareNamed(`SELECT id,alias,title,subtitle,thumbnail FROM courses`)
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

func (a *LearningRepository) GetCourseByAlias(alias string) *model.Course {
	// DB Response struct
	courses := make([]*model.Course, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"alias": alias,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT id,alias,title,subtitle,thumbnail FROM courses WHERE alias = :alias LIMIT 1`)
	if err != nil {
		fmt.Println("Error db GetCourseByAlias->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&courses, queryParams)
	if err != nil {
		fmt.Println("Error db GetCourseByAlias->query.Get : ", err)
	}
	if len(courses) <= 0 {
		return &model.Course{}
	}
	return courses[0]
}

// GetCourseByQuiz return course by quiz name
func (a *LearningRepository) GetCourseByQuiz(quiz string) *model.Course {
	// DB Response struct
	courses := make([]*model.Course, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"title": quiz,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT courses.id,courses.alias,courses.title,subtitle,thumbnail FROM courses JOIN course_contents ON courses.id = course_id WHERE course_contents.title = :title`)
	if err != nil {
		fmt.Println("Error db GetLessonsByCourseId->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&courses, queryParams)

	if err != nil {
		fmt.Println("Error db GetLessonsByCourseId->query.Get : ", err)
	}
	if len(courses) <= 0 {
		return &model.Course{}
	}
	return courses[0]
}

func (a *LearningRepository) GetLessonsByCourseId(id int, username string) []*model.SectionLessons {
	// DB Response struct
	sectionLessons := make([]*model.SectionLessons, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"course_id": id,
		"username":  username,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`
			SELECT section_name,section_desc,course_contents.course_content_id,type,title,permalink,duration,content,IFNULL(progress_time, 0) as progress_time
			FROM course_contents 
			LEFT JOIN course_contents_progress ON course_contents.course_content_id = course_contents_progress.course_content_id AND username = :username
			WHERE course_id = :course_id`)
	if err != nil {
		fmt.Println("Error db GetLessonsByCourseId->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&sectionLessons, queryParams)
	if err != nil {
		fmt.Println("Error db GetLessonsByCourseId->query.Get : ", err)
	}

	return sectionLessons
}

// GetLessonsByPermalink return lesson by permalink
func (a *LearningRepository) GetLessonByPermalink(course string, permalink string) *model.SectionLessons {
	// DB Response struct
	sectionLessons := make([]*model.SectionLessons, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"alias":     course,
		"permalink": permalink,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`
			SELECT section_name,section_desc,course_contents.course_content_id,type,course_contents.title,permalink,duration,content 
			FROM courses
			JOIN course_contents
			WHERE alias = :alias AND permalink = :permalink`)
	if err != nil {
		fmt.Println("Error db GetLessonsByPermalink->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&sectionLessons, queryParams)
	if err != nil {
		fmt.Println("Error db GetLessonsByPermalink->query.Get : ", err)
	}

	if len(sectionLessons) <= 0 {
		return &model.SectionLessons{}
	}
	return sectionLessons[0]
}

// AddCourseParticipant assign user to course
func (a *LearningRepository) AddCourseParticipant(username string, courseId int, status int) (affected int64) {

	// Data for query
	queryParams := map[string]interface{}{
		"username":  username,
		"course_id": courseId,
		"status":    status,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`INSERT INTO course_participants SET username = :username, course_id = :course_id, status = :status`)
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->RowsAffected : ", err)
	}

	return affected
}

func (a *LearningRepository) AddLessonProgress(username string, learningId int) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username":          username,
		"course_content_id": learningId,
		"progress_time":     0,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`INSERT INTO course_contents_progress SET username = :username, course_content_id = :course_content_id, progress_time = :progress_time`)
	if err != nil {
		fmt.Println("Error db AddCourseContentsProgress->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db AddCourseContentsProgress->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db AddCourseContentsProgress->RowsAffected : ", err)
	}

	return affected
}

// DeleteUserFromAllCourse remove user access to all course
func (a *LearningRepository) DeleteUserFromAllCourse(username string) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username": username,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`DELETE FROM course_participants WHERE username = :username`)
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->RowsAffected : ", err)
	}

	return affected
}

// DeleteUserAllLessonProgress remove user progress from all lessons
func (a *LearningRepository) DeleteUserAllLessonProgress(username string) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username": username,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`DELETE FROM course_contents_progress WHERE username = :username`)
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db AddCourseParticipant->RowsAffected : ", err)
	}

	return affected
}

// UpdateParticipantStatus change user access status in a course
func (a *LearningRepository) UpdateParticipantStatus(username string, id int, newStatus int) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username":  username,
		"course_id": id,
		"status":    newStatus,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`UPDATE course_participants SET status = :status WHERE username = :username AND course_id = :course_id`)
	if err != nil {
		fmt.Println("Error db InsertAnswer->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db InsertAnswer->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db InsertAnswer->RowsAffected : ", err)
	}

	return affected
}

// UpdateLearningVideoTimebar change user vide last progress time
func (a *LearningRepository) UpdateLearningVideoTimebar(username string, learningId int, time int) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username":          username,
		"course_content_id": learningId,
		"progress_time":     time,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`UPDATE course_contents_progress SET progress_time = :progress_time WHERE username = :username AND course_content_id = :course_content_id`)
	if err != nil {
		fmt.Println("Error db InsertAnswer->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db InsertAnswer->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db InsertAnswer->RowsAffected : ", err)
	}

	return affected
}

func (a *LearningRepository) InsertCourse(course *model.Course) (affected int64) {
	_, err := a.conn.NamedQuery(`INSERT INTO courses (alias, title, subtitle, thumbnail) 
											VALUES (:alias, :title, :subtitle, :thumbnail)`,
		course)
	if err != nil {
		fmt.Println("ERROR InsertQuestion ", err)
		return 0
	}
	return 1
}

func (a *LearningRepository) UpdateCourse(alias string, course *model.Course) (affected int64) {
	param := map[string]interface{}{
		"alias":     alias,
		"title":     course.Title,
		"subtitle":  course.Subtitle,
		"thumbnail": course.Thumbnail,
		"aliasNew":  course.Alias,
	}
	_, err := a.conn.NamedQuery(`
			UPDATE courses 
			SET alias = :aliasNew, title = :title, subtitle = :subtitle, thumbnail = :thumbnail 
			WHERE alias = :alias`,
		param)
	if err != nil {
		fmt.Println("ERROR UpdateQuestion ", err)
		return 0
	}
	return 1
}

func (a *LearningRepository) DeleteCourse(course *model.Course) (affected int64) {
	_, err := a.conn.NamedQuery(`DELETE FROM courses WHERE alias = :alias`,
		course)
	if err != nil {
		fmt.Println("ERROR DeleteQuestion ", err)
		return 0
	}
	return 1
}

func (a *LearningRepository) SaveCourseContent(courseId int, sectionName string, sectionDesc string, content *model.Lesson) (affected int64) {

	// Data for query
	queryParams := map[string]interface{}{
		"courseID":    courseId,
		"type":        content.Type,
		"title":       content.Title,
		"permalink":   content.Permalink,
		"content":     content.Video,
		"sectionName": sectionName,
		"sectionDesc": sectionDesc,
		"duration":    content.Duration,
	}
	fmt.Println("queryParams ", queryParams)
	// Compose query
	query, err := a.conn.PrepareNamed(`
		INSERT INTO course_contents 
		(course_id, type, title, permalink, content, section_name, section_desc, duration) 
		VALUES (:courseID, :type, :title, :permalink, :content, :sectionName, :sectionDesc, :duration)
	`)
	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db SaveCourseContent->query.Exec : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db InsertAnswer->RowsAffected : ", err)
	}

	return affected
}

func (a *LearningRepository) FetchSectionContentByCourseAndSection(courseID int, sectionName string) *model.SectionLessons {
	// DB Response struct
	sections := make([]*model.SectionLessons, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"courseID":    courseID,
		"sectionName": sectionName,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`
		SELECT section_name, section_desc
		FROM course_contents 
		WHERE course_id = :courseID AND section_name = :sectionName
		LIMIT 1
	`)
	if err != nil {
		fmt.Println("Error db GetCourseByAlias->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&sections, queryParams)
	if err != nil {
		fmt.Println("Error db FetchSectionContentByCourseAndSection->query.Select : ", err)
	}
	if len(sections) <= 0 {
		return &model.SectionLessons{}
	}
	return sections[0]
}
