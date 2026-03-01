ifneq (,$(wildcard .env))
    include .env
    export
endif

SQLC_CONFIG_PATH=./server/sqlc/config.yaml
DATABASE_URL=postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE_NAME)?sslmode=disable

gen-server-proto:
	@echo "Generating server protobuf ..."
	protoc \
    	--go_out=. \
    	--go_opt=paths=source_relative \
    	shared/packets.proto
	mv shared/packets.pb.go server/pkg/network/protocol/packets.pb.go

sql-generate:
	@echo "Generating SQL source code ..."
	@sqlc generate -f $(SQLC_CONFIG_PATH)
	@echo "Done"
