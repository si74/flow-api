test: 
	go test ./... -v -race

run: 
	go run cmd/flowd/main.go

staticcheck: 
	staticcheck ./...

lint: 
	golint -set_exit_status $(PWD)/cmd/... $(PWD)/internal/...


.PHONY: staticcheck test run lint