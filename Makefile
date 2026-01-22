proto-gen:
	protoc --go_out=. --go-grpc_out=. proto/course_service.proto
	protoc --go_out=. --go-grpc_out=. proto/module_service.proto
	protoc --go_out=. --go-grpc_out=. proto/lesson_service.proto
	protoc --go_out=. --go-grpc_out=. proto/group_course_service.proto