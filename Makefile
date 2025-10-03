
run_local:
	 ENV_FILE=./.env.local go run cmd/main.go


deploy-server:
	./deploy/deploy.sh

server-logs:
	./deploy/logs.sh

swag:
	swag init -g cmd/main.go
