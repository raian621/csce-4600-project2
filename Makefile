.PHONY: cover
cover: test
	-go tool cover -html=coverage.out -o coverage.html
	-go tool cover -func=coverage.out

.PHONY: test
test:
	-go test -v -coverprofile=coverage.out ./...