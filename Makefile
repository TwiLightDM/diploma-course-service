proto-gen:
	protoc --go_out=. --go-grpc_out=. proto/course_service.proto
	protoc --go_out=. --go-grpc_out=. proto/module_service.proto
	protoc --go_out=. --go-grpc_out=. proto/lesson_service.proto
	protoc --go_out=. --go-grpc_out=. proto/lesson_file_service.proto
	protoc --go_out=. --go-grpc_out=. proto/group_course_service.proto
	protoc --go_out=. --go-grpc_out=. proto/task_service.proto
	protoc --go_out=. --go-grpc_out=. proto/task_attempt_service.proto
	protoc --go_out=. --go-grpc_out=. proto/lesson_progress_service.proto
	protoc --go_out=. --go-grpc_out=. proto/completed_module_service.proto
	protoc --go_out=. --go-grpc_out=. proto/completed_course_service.proto
	protoc --go_out=. --go-grpc_out=. proto/completed_theory_course_service.proto