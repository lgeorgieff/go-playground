.PHONY: generate-proto
generate-proto:
	protoc \
		--go_out=. \
		--go_opt=module=github.com/lgeorgieff/go-playground \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/lgeorgieff/go-playground \
		proto/dummy/v1/dummy.proto

.PHONY: generate
generate: generate-proto
