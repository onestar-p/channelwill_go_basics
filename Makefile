#编译可执行文件
gatewayBuildFileName="./bin/gateway"
serverBuildFileName="./bin/server"
build:
	go build -x -o ${gatewayBuildFileName} cmd/client/main.go
	go build -x -o ${serverBuildFileName} cmd/server/main.go

	@echo ""
	@echo "########## build successfully! ##########"
	@echo "Gateway file path:" ${gatewayBuildFileName}
	@echo "Server file path:" ${serverBuildFileName}
	@echo ""

rc:
	go run cmd/client/main.go
	
rs:
	go run cmd/server/main.go

gen-proto:
	cd proto; ./genProto.sh