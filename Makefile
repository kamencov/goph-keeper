.PHONY: gen-proto
gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./internal/proto/v1/goph_keeper_v1.proto

.PHONY: cover
cover:
	go test -short -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: test
test:
	go test ./... -coverprofile=cover.out
	grep -Ev "./cli/|mock.go|./proto/|./internal/migrations/|run.go|./docs/" cover.out > cover.filtered.out
	go tool cover -func=cover.filtered.out
	rm cover.out cover.filtered.out

.PHONY: git
git:
	git add .
	git commit -m "update"
	git push origin iter3

.PHONY: swag
swag:
	swag init -g cmd/service/main.go --parseInternal

.PHONY: build_mac
build_mac:
	GOARCH=amd64 go build -o myapp_amd64 ./cmd/client
	GOARCH=arm64 go build -o myapp_arm64 ./cmd/client
	lipo -create -output myapp myapp_amd64 myapp_arm64

.PHONY: build_win
build_win:
	GOARCH=amd64 go build -o myapp.exe ./cmd/client

.PHONY: build_linux
build_linux:
	GOARCH=amd64 go build -o myapp ./cmd/client