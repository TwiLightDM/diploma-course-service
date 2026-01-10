package entities

type Lesson struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Position    int64  `json:"position"`
	ModuleId    string `json:"module_id"`
}
