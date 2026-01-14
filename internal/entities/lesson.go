package entities

type Lesson struct {
	Id          string
	Title       string
	Description string
	Content     string
	Position    int64
	ModuleId    string
}
