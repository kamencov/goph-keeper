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
	grep -Ev "./cli/|mock.go|./proto/" cover.out > cover.filtered.out
	go tool cover -func=cover.filtered.out
	rm cover.out cover.filtered.out

.PHONY: git
git:
	git add .
	git commit -m "update"
	git push origin iter3