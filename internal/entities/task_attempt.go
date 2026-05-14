package entities

type TaskAttemptAnswer struct {
	TaskId            string   `bson:"task_id"`
	TextAnswer        string   `bson:"text_answer"`
	SelectedOptionIds []string `bson:"selected_option_ids"`
	IsCorrect         bool     `bson:"is_correct"`
}

type TaskAttempt struct {
	Id       string `bson:"_id,omitempty"`
	UserId   string `bson:"user_id"`
	CourseId string `bson:"course_id"`
	ModuleId string `bson:"module_id"`

	AttemptNumber  int                 `bson:"attempt_number"`
	Answers        []TaskAttemptAnswer `bson:"answers"`
	CorrectAnswers int                 `bson:"correct_answers"`
	TotalQuestions int                 `bson:"total_questions"`
	Score          float64             `bson:"score"`
}
