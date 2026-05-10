package entities

type TaskType string

const (
	TaskTypeTextInput TaskType = "TEXT_INPUT"
	TaskTypeChoice    TaskType = "CHOICE"
)

type TaskOption struct {
	Id        string `bson:"_id"`
	Text      string `bson:"text"`
	IsCorrect bool   `bson:"is_correct"`
}

type Task struct {
	Id string `bson:"_id"`

	Title string `bson:"title"`

	CourseId string `bson:"course_id"`
	ModuleId string `bson:"module_id"`

	Type TaskType `bson:"type"`

	CorrectTextAnswer string `bson:"correct_text_answer"`

	Options []TaskOption `bson:"options"`
}
