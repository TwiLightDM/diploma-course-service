package entities

type LessonFile struct {
	Id          string
	FileName    string
	ObjectName  string
	ContentType string
	LessonId    string
	Size        int64
	Url         string `gorm:"->"`
}
