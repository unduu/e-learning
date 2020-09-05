package usecase

import (
	"github.com/unduu/e-learning/evaluation"
	"github.com/unduu/e-learning/learning"
	"github.com/unduu/e-learning/learning/model"
	"strings"
)

type LearningUsecase struct {
	repository   learning.Repository
	evaluationUC evaluation.Usecase
}

func NewLearningUsecase(repository learning.Repository, evaluationUC evaluation.Usecase) *LearningUsecase {
	return &LearningUsecase{
		repository:   repository,
		evaluationUC: evaluationUC,
	}
}

// GetCourseList return course / modules list
func (a *LearningUsecase) GetCourseList() (results []*model.Course) {
	courseArr := a.repository.GetCourses()
	for _, course := range courseArr {
		participants := a.repository.GetParticipantByCourse(course.Id)
		// Get detail total lessons & hours
		course := a.GetCourseLessons(course.Alias, "")
		course.AddParticipant(participants)
		results = append(results, course)
	}

	return results
}

// GetCourseLessons return lessons in a course
func (a *LearningUsecase) GetCourseLessons(alias string, username string) (course *model.Course) {

	course = a.repository.GetCourseByAlias(alias)
	participants := a.repository.GetParticipantByCourse(course.Id)
	course.AddParticipant(participants)
	data := a.repository.GetLessonsByCourseId(course.Id, username)
	for _, row := range data {
		lessonSplit := a.repository.FetchCourseContentSplit(row.LessonID)
		lesson := &model.Lesson{
			Type:      row.Type,
			Title:     row.Title,
			Permalink: row.Permalink,
			Duration:  row.Duration,
			Video:     row.Video,
			Timebar:   row.Timebar,
			Progress:  0,
			Split:     lessonSplit,
		}

		currSection := course.GetSection(row.Name)
		section := currSection
		if currSection == nil {
			section = &model.Section{
				Name: row.Name,
				Desc: row.Desc,
			}
		}
		section.AddLesson(lesson)

		if currSection == nil {
			course.AddSection(section)
		}
	}

	return course
}

// SetDefaultCourse set default course for new registered user
func (a *LearningUsecase) SetDefaultCourse(username string) {
	courseArr := []int{1, 2, 3, 4, 5}
	a.repository.DeleteUserFromAllCourse(username)
	a.repository.DeleteUserAllLessonProgress(username)
	for i, courseId := range courseArr {
		// Set first course status to open
		status := 0
		if i == 0 {
			status = 1
		}
		// Set default course for user
		a.repository.AddCourseParticipant(username, courseId, status)
		// Set default learning progress
		lessons := a.repository.GetLessonsByCourseId(courseId, "")
		for _, lesson := range lessons {
			a.repository.AddLessonProgress(username, lesson.LessonID)
		}
	}
}

// UpdateUserCourseProgress set user course last progress based on quiz result
func (a *LearningUsecase) UpdateUserCourseProgress(username string, quiz string) {
	pass := true
	// Get course by quiz name
	course := a.repository.GetCourseByQuiz(quiz)
	// Get lesson from quiz
	lessons := a.repository.GetLessonsByCourseId(course.Id, username)
	for _, row := range lessons {
		lesson := &model.Lesson{
			Type:     row.Type,
			Title:    row.Title,
			Duration: row.Duration,
			Video:    row.Video,
		}
		lesson = a.SetLessonProgress(username, lesson)
		if lesson.Progress != 2 && lesson.IsQuiz() {
			pass = false
		}
	}

	if pass {
		// Mark current course as finish
		a.repository.UpdateParticipantStatus(username, course.Id, course.GetStatusCode("completed"))
		// Open next course
		a.repository.UpdateParticipantStatus(username, course.GetNextCourseId(), course.GetStatusCode("open"))
	}
}

// SetLessonProgress set user lesson progress
func (a *LearningUsecase) SetLessonProgress(username string, lesson *model.Lesson) *model.Lesson {
	if lesson.IsQuiz() {
		lesson.Progress = 0
		exist, answer := a.evaluationUC.IsAnswerExists(username, lesson.Title)
		if !exist {
			return lesson
		}

		lesson.Progress = 1
		if answer.Grade == 100 {
			lesson.Progress = 2
		}

	}
	return lesson
}

func (a *LearningUsecase) UpdateVideoProgress(username string, course string, lesson string, time int) {
	les := a.repository.GetLessonByPermalink(course, lesson)
	a.repository.UpdateLearningVideoTimebar(username, les.LessonID, time)
}

func (a *LearningUsecase) AddCourse(title string, subtitle string, thumbnail string) {
	lower := strings.ToLower(title)
	permalink := strings.Replace(lower, " ", "-", -1)
	course := &model.Course{
		Title:     title,
		Subtitle:  subtitle,
		Alias:     permalink,
		Thumbnail: thumbnail,
	}
	a.repository.InsertCourse(course)
}

func (a *LearningUsecase) EditCourse(alias string, title string, subtitle string, thumbnail string) {
	lower := strings.ToLower(title)
	permalink := strings.Replace(lower, " ", "-", -1)
	course := &model.Course{
		Alias:     permalink,
		Title:     title,
		Subtitle:  subtitle,
		Thumbnail: thumbnail,
	}
	a.repository.UpdateCourse(alias, course)
}

func (a *LearningUsecase) DeleteCourse(alias string) {
	course := &model.Course{
		Alias: alias,
	}
	a.repository.DeleteCourse(course)
}

func (a *LearningUsecase) AddCourseContent(courseAlias string, sectionName string, sectionDesc string, lessonType string, title string, video string) {

	content := &model.Lesson{
		Type:  lessonType,
		Title: title,
		Video: video,
	}
	content.GeneratePermalink()

	course := a.repository.GetCourseByAlias(courseAlias)
	a.repository.SaveCourseContent(course.Id, sectionName, sectionDesc, content)
}
