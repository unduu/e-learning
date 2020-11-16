package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/unduu/e-learning/evaluation/model"
	"strings"
)

type EvaluationRepository struct {
	conn *sqlx.DB
}

func NewEvaluationRepository(db *sqlx.DB) *EvaluationRepository {
	return &EvaluationRepository{
		conn: db,
	}
}

// GetQuestions return question list
func (a *EvaluationRepository) GetQuestions(module string, page int, limit int) ([]*model.Question, int, error) {
	offset := (page - 1) * limit
	questions := make([]*model.Question, 0)

	var count struct {
		Total int `db:"total"`
	}
	queryParams := map[string]interface{}{
		"module": module,
		"offset": offset,
		"limit":  limit,
	}
	query, err := a.conn.PrepareNamed(`SELECT id, type, attachment_type, attachment, question, choices FROM questions WHERE module = :module ORDER BY RAND() LIMIT :offset, :limit `)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	queryTotal, err := a.conn.PrepareNamed(`SELECT COUNT(*) AS total FROM questions WHERE module = :module`)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	err = query.Select(&questions, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	err = queryTotal.Get(&count, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}
	count.Total = 5
	return questions, count.Total, err
}

// GetQuestions return question list
func (a *EvaluationRepository) GetPrePostQuestions() ([]*model.Question, int, error) {
	questions := make([]*model.Question, 0)

	var count struct {
		Total int `db:"total"`
	}
	queryParams := map[string]interface{}{}
	query, err := a.conn.PrepareNamed(`
		select id, type, attachment_type, attachment, question, choices
		from (
				 select
					 id,
					 type,
					 attachment_type,
					 attachment,
					 choices,
					 question,
					 module,
					 @group_rank := IF(@current_group=module, @group_rank + 1, 1) as  group_rank,
					 @current_group := module
				 from (
						  select
							  id,
							  type,
							  attachment_type,
							  attachment,
							  choices,
							  question,
							  module,
							  CONCAT(module, '-', round(rand() * 100)) as rand_rank
						  from questions order by rand_rank
					  ) tmp
			 ) ranked
		where group_rank <= 2 AND module !="prepost" order by module
	`)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	queryTotal, err := a.conn.PrepareNamed(`SELECT COUNT(*) AS total FROM questions WHERE module != "prepost GROUP BY module"`)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	err = query.Select(&questions, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	err = queryTotal.Get(&count, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}
	count.Total = 10
	return questions, count.Total, err
}

// GetQuestionByIds return question by specified id
func (a *EvaluationRepository) GetQuestionByIds(ids []string, page int, limit int) ([]*model.Question, int, error) {
	offset := (page - 1) * limit
	questions := make([]*model.Question, 0)

	var count struct {
		Total int `db:"total"`
	}

	queryParams := map[string]interface{}{
		"ids":    strings.Join(ids, ","),
		"offset": offset,
		"limit":  limit,
	}
	query, err := a.conn.PrepareNamed(`SELECT id, type, attachment_type, attachment, question, choices FROM questions 
												WHERE id IN (` + strings.Join(ids, ",") + `)
												ORDER BY RAND() LIMIT :offset, :limit `)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	queryTotal, err := a.conn.PrepareNamed(`SELECT COUNT(*) AS total FROM questions WHERE id IN (` + strings.Join(ids, ",") + `)`)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	err = query.Select(&questions, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	err = queryTotal.Get(&count, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	return questions, count.Total, err
}

// GetCorrectAnswerByQuestionID return correct answer for specific question
func (a *EvaluationRepository) GetCorrectAnswerByQuestionID(id int) *model.CorrectAnswerDB {
	// DB Response struct
	answers := make([]*model.CorrectAnswerDB, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"id": id,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT id,answer FROM questions WHERE id = :id`)
	if err != nil {
		fmt.Println("Error db GetCorrectAnswerByQuestionID->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&answers, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	return answers[0]
}

// GetUserAnswers return last user answer in specific pre,post, or quiz test
func (a *EvaluationRepository) GetUserAnswers(username string, module string) *model.UserAnswerDB {
	// DB Response struct
	answers := make([]*model.UserAnswerDB, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"username": username,
		"module":   module,
		"status":   "latest",
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT username,answer,grade FROM answers WHERE username = :username AND type = :module AND status = :status`)
	if err != nil {
		fmt.Println("Error db GetCorrectAnswerByQuestionID->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&answers, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	if len(answers) <= 0 {
		return nil
	}
	return answers[0]
}

// InsertAnswer persis to database
func (a *EvaluationRepository) InsertAnswer(username string, testType string, answer string, grade float64) (affected int64) {
	gradeStr := fmt.Sprintf("%.1f", grade)
	// Data for query
	queryParams := map[string]interface{}{
		"username": username,
		"type":     testType,
		"answer":   answer,
		"grade":    gradeStr,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`INSERT INTO answers SET username = :username, type = :type, answer = :answer, grade = :grade, updated = CURRENT_TIMESTAMP ON DUPLICATE KEY UPDATE updated = CURRENT_TIMESTAMP`)
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

// UpdateUserAnswerStatus change latest status to archived
func (a *EvaluationRepository) UpdateUserAnswerStatus(username string, module string, newStatus string) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username": username,
		"type":     module,
		"status":   newStatus,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`UPDATE answers SET status = :status WHERE username = :username AND type = :type AND status = "latest"`)
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

// InsertQuestion Add a new question
func (a *EvaluationRepository) InsertQuestion(question model.Question) {
	_, err := a.conn.NamedQuery(`INSERT INTO questions (module, type,question,choices,answer) 
											VALUES (:module, :type, :question, :choices, :answer)`,
		question)
	if err != nil {
		fmt.Println("ERROR InsertQuestion ", err)
	}
}

// UpdateQuestion update question
func (a *EvaluationRepository) UpdateQuestion(question model.Question) {
	_, err := a.conn.NamedQuery(`UPDATE questions SET question = :question, choices = :choices, answer = :answer WHERE id = :id`,
		question)
	if err != nil {
		fmt.Println("ERROR UpdateQuestion ", err)
	}
}

// DeleteQuestion remove question from database
func (a *EvaluationRepository) DeleteQuestion(question model.Question) {
	_, err := a.conn.NamedQuery(`DELETE FROM questions WHERE id = :id`,
		question)
	if err != nil {
		fmt.Println("ERROR DeleteQuestion ", err)
	}
}

func (a *EvaluationRepository) GetAllQuestions(page int, limit int) ([]*model.Question, int, error) {
	offset := (page - 1) * limit
	questions := make([]*model.Question, 0)

	var count struct {
		Total int `db:"total"`
	}
	queryParams := map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	}
	query, err := a.conn.PrepareNamed(`SELECT id, type, attachment_type, attachment, question, choices FROM questions ORDER BY id DESC LIMIT :offset, :limit `)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	queryTotal, err := a.conn.PrepareNamed(`SELECT COUNT(*) AS total FROM questions`)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	err = query.Select(&questions, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	err = queryTotal.Get(&count, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	return questions, count.Total, err
}

func (a *EvaluationRepository) FetchQuestionGroups() []*model.QuestionGroup {
	// DB Response struct
	questionGroups := make([]*model.QuestionGroup, 0)

	// Data for query
	queryParams := map[string]interface{}{}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT module,count(*) as total_questions FROM questions GROUP BY module;`)
	if err != nil {
		fmt.Println("Error db GetCorrectAnswerByQuestionID->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&questionGroups, queryParams)
	if err != nil {
		fmt.Println("Error db FetchQuestionGroups->query.Select : ", err)
	}

	return questionGroups
}

func (a *EvaluationRepository) FetchAvailableQuestionGroups() []*model.QuestionGroup {
	// DB Response struct
	questionGroups := make([]*model.QuestionGroup, 0)

	// Data for query
	queryParams := map[string]interface{}{}

	// Compose query
	query, err := a.conn.PrepareNamed(`
		SELECT questions.module
		FROM questions
		LEFT JOIN course_contents cc on questions.module = cc.title
		WHERE questions.module != "prepost" AND course_content_id is null
	`)
	if err != nil {
		fmt.Println("Error db FetchAvailableQuestionGroups->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&questionGroups, queryParams)
	if err != nil {
		fmt.Println("Error db FetchAvailableQuestionGroups->query.Select : ", err)
	}

	return questionGroups
}

func (a *EvaluationRepository) FetchQuestionsByModule(moduleName string, page int, limit int) ([]*model.Question, int, error) {
	offset := (page - 1) * limit
	questions := make([]*model.Question, 0)

	var count struct {
		Total int `db:"total"`
	}
	queryParams := map[string]interface{}{
		"module": moduleName,
		"offset": offset,
		"limit":  limit,
	}
	query, err := a.conn.PrepareNamed(`
			SELECT id, type, attachment_type, attachment, question, choices 
			FROM questions 
			WHERE module = :module 
			ORDER BY id DESC 
			LIMIT :offset, :limit 
	`)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	queryTotal, err := a.conn.PrepareNamed(`SELECT COUNT(*) AS total FROM questions WHERE module = :module `)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	err = query.Select(&questions, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	err = queryTotal.Get(&count, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	return questions, count.Total, err
}

// DeleteGroupByName remove question from database
func (a *EvaluationRepository) DeleteGroupByName(name string) {
	q := model.Question{
		Module: name,
	}
	_, err := a.conn.NamedQuery(`DELETE FROM questions WHERE module = :module`, q)
	if err != nil {
		fmt.Println("ERROR DeleteGroupByName ", err)
	}
}
