.PHONY: generate-dummy-proto
generate-dummy-proto:
	protoc \
		--go_out=. \
		--go_opt=module=github.com/lgeorgieff/go-playground \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/lgeorgieff/go-playground \
		proto/dummy/v1/dummy.proto

.PHONY: generate-todo-proto
generate-todo-proto:
	protoc \
		--go_out=. \
		--go_opt=module=github.com/lgeorgieff/go-playground \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/lgeorgieff/go-playground \
		proto/todo/v1/todo.proto

.PHONY: generate-proto
generate-proto: generate-dummy-proto generate-todo-proto

.PHONY: generate
generate: generate-proto
