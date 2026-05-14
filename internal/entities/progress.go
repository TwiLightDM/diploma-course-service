package entities

type LessonProgress struct {
	UserId   string
	LessonId string
}

type ModuleProgress struct {
	ModuleId           string
	TotalLessons       int64
	CompletedLessons   int64
	ProgressPercent    float64
	CompletedLessonIds []string
}

type UserModuleProgress struct {
	UserId           string
	CompletedLessons int64
	TotalLessons     int64
	ProgressPercent  float64
	Completed        bool
}

type CourseProgress struct {
	CourseId           string
	TotalLessons       int64
	CompletedLessons   int64
	ProgressPercent    float64
	CompletedLessonIds []string
}

type UserCourseProgress struct {
	UserId           string
	CompletedLessons int64
	TotalLessons     int64
	ProgressPercent  float64
	Completed        bool
}
