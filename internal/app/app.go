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
	"github.com/TwiLightDM/diploma-course-service/internal/services/lesson-service"
	"github.com/TwiLightDM/diploma-course-service/internal/services/module-service"
	"github.com/TwiLightDM/diploma-course-service/package/databases"
	"github.com/TwiLightDM/diploma-course-service/proto/courseservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/groupcourseservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/lessonservicepb"
	"github.com/TwiLightDM/diploma-course-service/proto/moduleservicepb"

	"google.golang.org/grpc"
)

func Run(cfg *config.Config) error {
	db, err := databases.InitDB(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
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

	courseRepo := course_service.NewCourseRepository(db)
	courseService := course_service.NewCourseService(courseRepo)
	courseHandler := course_service.NewCourseHandler(courseService)

	moduleRepo := module_service.NewModuleRepository(db)
	moduleService := module_service.NewModuleService(moduleRepo)
	moduleHandler := module_service.NewModuleHandler(moduleService)

	lessonRepo := lesson_service.NewLessonRepository(db)
	lessonService := lesson_service.NewLessonService(lessonRepo)
	lessonHandler := lesson_service.NewLessonHandler(lessonService)

	groupCourseRepo := group_course_service.NewGroupCourseRepository(db)
	groupCourseService := group_course_service.NewGroupCourseService(groupCourseRepo)
	groupCourseHandler := group_course_service.NewGroupCourseHandler(groupCourseService)

	courseservicepb.RegisterCourseServiceServer(grpcServer, courseHandler)
	moduleservicepb.RegisterModuleServiceServer(grpcServer, moduleHandler)
	lessonservicepb.RegisterLessonServiceServer(grpcServer, lessonHandler)
	groupcourseservicepb.RegisterGroupCourseServiceServer(grpcServer, groupCourseHandler)

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

	sqlDB, err := db.DB()
	if err == nil {
		log.Println("Closing database connection...")
		_ = sqlDB.Close()
	}

	log.Println("Course-service stopped gracefully")
	return nil
}
