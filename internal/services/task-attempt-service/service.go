package task_attempt_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-course-service/internal/entities"
	"github.com/google/uuid"
)

type TaskAttemptRepository interface {
	Create(ctx context.Context, taskAttempt *entities.TaskAttempt) error
	ReadById(ctx context.Context, id string) (*entities.TaskAttempt, error)
	ReadAllByUserIdAndModuleId(ctx context.Context, userId, moduleId string) ([]entities.TaskAttempt, error)
	ReadAllByUserIdAndCourseId(ctx context.Context, userId, courseId string) ([]entities.TaskAttempt, error)
	GetLastAttemptNumber(ctx context.Context, userId string, courseId string, moduleId string) (int, error)
}

type TaskService interface {
	ReadTaskById(ctx context.Context, id string) (*entities.Task, error)
}

type taskAttemptService struct {
	repo        TaskAttemptRepository
	taskService TaskService
}

func NewTaskAttemptService(repo TaskAttemptRepository, taskService TaskService) TaskAttemptService {
	return &taskAttemptService{repo: repo, taskService: taskService}
}

func (s *taskAttemptService) SubmitTaskAttempt(ctx context.Context, taskAttempt *entities.TaskAttempt) (*entities.TaskAttempt, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	lastAttemptNumber, err := s.repo.GetLastAttemptNumber(ctx, taskAttempt.UserId, taskAttempt.CourseId, taskAttempt.ModuleId)
	if err != nil {
		return nil, err
	}

	taskAttempt.Id = uuid.NewString()
	taskAttempt.AttemptNumber = lastAttemptNumber + 1

	correctAnswers := 0

	for i, answer := range taskAttempt.Answers {
		task, err := s.taskService.ReadTaskById(ctx, answer.TaskId)
		if err != nil {
			return nil, err
		}

		switch task.Type {
		case entities.TaskTypeTextInput:
			if answer.TextAnswer != "" {
				isCorrect := answer.TextAnswer == task.CorrectTextAnswer
				taskAttempt.Answers[i].IsCorrect = isCorrect

				if isCorrect {
					correctAnswers++
				}
			}

		case entities.TaskTypeChoice:
			correctOptionIds := make([]string, 0)
			for _, option := range task.Options {
				if option.IsCorrect {
					correctOptionIds = append(correctOptionIds, option.Id)
				}
			}

			isCorrect := compareStringSlices(answer.SelectedOptionIds, correctOptionIds)
			taskAttempt.Answers[i].IsCorrect = isCorrect

			if isCorrect {
				correctAnswers++
			}
		}
	}

	taskAttempt.TotalQuestions = len(taskAttempt.Answers)

	taskAttempt.CorrectAnswers = correctAnswers

	if taskAttempt.TotalQuestions > 0 {
		taskAttempt.Score = (float64(correctAnswers) / float64(taskAttempt.TotalQuestions)) * 100
	}

	err = s.repo.Create(ctx, taskAttempt)
	if err != nil {
		return nil, err
	}

	return taskAttempt, nil
}

func (s *taskAttemptService) ReadTaskAttemptById(ctx context.Context, id string) (*entities.TaskAttempt, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.ReadById(ctx, id)
}

func (s *taskAttemptService) ReadAllTaskAttemptsByUserIdAndModuleId(ctx context.Context, userId, moduleId string) ([]entities.TaskAttempt, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.ReadAllByUserIdAndModuleId(ctx, userId, moduleId)
}

func (s *taskAttemptService) ReadAllTaskAttemptsByUserIdAndCourseId(ctx context.Context, userId, courseId string) ([]entities.TaskAttempt, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	return s.repo.ReadAllByUserIdAndCourseId(ctx, userId, courseId)
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int)

	for _, value := range a {
		m[value]++
	}

	for _, value := range b {
		if _, ok := m[value]; !ok {
			return false
		}

		m[value]--
		if m[value] < 0 {
			return false
		}
	}

	return true
}
