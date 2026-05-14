package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/TwiLightDM/diploma-course-service/internal/config"
	"github.com/TwiLightDM/diploma-course-service/internal/services/course-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/group-course-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/lesson-file-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/lesson-progress-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/lesson-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/module-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/task-attempt-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/task-service"
	"github.com/TwiLightDM/diploma-course-service/package/databases"
	"github.com/TwiLightDM/diploma-course-service/proto/courseservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/groupcourseservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonfileservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonprogressservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/moduleservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/taskattemptservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/taskservicepb"

	"google.golang.org/grpc"
)

func Run(cfg *config.Config) error {
	postgres, err := databases.InitDB(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
	)
	if err != nil {
		return err
	}

	storage, err := databases.NewStorage(
		cfg.Storage.Endpoint,
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.BucketName,
	)
	if err != nil {
		return err
	}

	mongo, err := databases.InitMongoDB(
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
		cfg.Mongo.Name,
	)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return err
	}

	log.Printf("Starting course-service on %s", listener.Addr().String())

	grpcServer := grpc.NewServer()

	courseRepo := course_service.NewCourseRepository(postgres)
	courseService := course_service.NewCourseService(courseRepo)
	courseHandler := course_service.NewCourseHandler(courseService)

	moduleRepo := module_service.NewModuleRepository(postgres)
	moduleService := module_service.NewModuleService(moduleRepo)
	moduleHandler := module_service.NewModuleHandler(moduleService)

	lessonRepo := lesson_service.NewLessonRepository(postgres)
	lessonStore := lesson_service.NewLessonStorage(storage)
	lessonService := lesson_service.NewLessonService(lessonRepo, lessonStore)
	lessonHandler := lesson_service.NewLessonHandler(lessonService)

	lessonFileRepo := lesson_file_service.NewLessonFileRepository(postgres)
	lessonFileStore := lesson_file_service.NewLessonFileStorage(storage)
	lessonFileService := lesson_file_service.NewLessonFileService(lessonFileRepo, lessonFileStore)
	lessonFileHandler := lesson_file_service.NewLessonFileHandler(lessonFileService)

	taskRepo := task_service.NewTaskRepository(mongo)
	taskService := task_service.NewTaskService(taskRepo)
	taskHandler := task_service.NewTaskHandler(taskService)

	taskAttemptRepo := task_attempt_service.NewTaskAttemptRepository(mongo)
	taskAttemptService := task_attempt_service.NewTaskAttemptService(taskAttemptRepo, taskService)
	taskAttemptHandler := task_attempt_service.NewTaskAttemptHandler(taskAttemptService)

	groupCourseRepo := group_course_service.NewGroupCourseRepository(postgres)
	groupCourseService := group_course_service.NewGroupCourseService(groupCourseRepo)
	groupCourseHandler := group_course_service.NewGroupCourseHandler(groupCourseService)

	lessonProgressRepo := lesson_progress_service.NewLessonProgressRepository(postgres)
	lessonProgressService := lesson_progress_service.NewLessonProgressService(lessonProgressRepo)
	lessonProgressHandler := lesson_progress_service.NewLessonProgressHandler(lessonProgressService)

	courseservicepb.RegisterCourseServiceServer(grpcServer, courseHandler)
	moduleservicepb.RegisterModuleServiceServer(grpcServer, moduleHandler)
	lessonservicepb.RegisterLessonServiceServer(grpcServer, lessonHandler)
	lessonfileservicepb.RegisterLessonFileServiceServer(grpcServer, lessonFileHandler)
	taskservicepb.RegisterTaskServiceServer(grpcServer, taskHandler)
	taskattemptservicepb.RegisterTaskAttemptServiceServer(grpcServer, taskAttemptHandler)
	groupcourseservicepb.RegisterGroupCourseServiceServer(grpcServer, groupCourseHandler)
	lessonprogressservicepb.RegisterLessonProgressServiceServer(grpcServer, lessonProgressHandler)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down gRPC server...")

	grpcServer.GracefulStop()

	sqlDB, err := postgres.DB()
	if err == nil {
		log.Println("Closing database connection...")
		_ = sqlDB.Close()
	}

	if err = mongo.Client().Disconnect(ctx); err != nil {
		log.Println("Error while disconnecting mongo:", err)
	}
	log.Println("Closing mongo connection...")

	log.Println("Course-service stopped gracefully")
	return nil
}
