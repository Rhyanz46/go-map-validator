test-for-publish:
	go clean -testcache
	@go test ./... -coverprofile=cover.out && go tool cover -html=cover.out

tests-silent:
	go clean -testcache
	@go test ./...

tests:
	go clean -testcache
	go test -v ./...
