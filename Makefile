.PHONY: help run build

build:
	go build -o cloudsync$(extension) ./cmd/cli/main.go

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
