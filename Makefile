test: 
	go test ./... -v -race

run: 
	go run cmd/flowd/main.go

staticcheck: 
	staticcheck ./...


.PHONY: staticcheck test run