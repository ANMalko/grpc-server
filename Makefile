.PHONY: build
build:
	go build -o grpc ./cmd

.PHONY: install-tools
install-tools:
	go get -v
	go install -v \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: generate
generate:
	protoc -I ./proto \
		--go_out ./proto --go_opt paths=source_relative \
		--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt logtostderr=true \
		./proto/users/users.proto

.PHONY: generate_docs
generate_docs:
	protoc -I ./proto --openapiv2_out ./api \
    --openapiv2_opt logtostderr=true,allow_merge=true,merge_file_name=api \
	./proto/users/users.proto \
