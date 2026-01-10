package entities

type Module struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int64  `json:"position"`
	CourseId    string `json:"course_id"`
}
