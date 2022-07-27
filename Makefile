all: build-win build-lin build-darwin

build-win:
	GOOS=windows GOARCH=amd64 go build -o ./bin/win/amd64/cloudsync.exe ./cmd/cli/main.go
	GOOS=windows GOARCH=arm64 go build -o ./bin/win/arm64/cloudsync.exe ./cmd/cli/main.go

build-lin:
	GOOS=linux GOARCH=amd64 go build -o ./bin/lin/amd64/cloudsync ./cmd/cli/main.go
	GOOS=linux GOARCH=arm64 go build -o ./bin/lin/arm64/cloudsync ./cmd/cli/main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/amd64/cloudsync ./cmd/cli/main.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin/arm64/cloudsync ./cmd/cli/main.go

test:
	go test ./... -v --cover

test-cov:
	go test ./... -coverprofile coverage.out . && go tool cover -html=coverage.out

cloud-plan:
	cd ./deployments/terraform/workspaces/$(stage) && terraform plan

cloud-deploy:
	cd ./deployments/terraform/workspaces/$(stage) && terraform apply

cloud-destroy:
	cd ./deployments/terraform/workspaces/$(stage) && terraform destroy
