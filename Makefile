.PHONY: gen-proto
gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./internal/proto/v1/goph_keeper_v1.proto