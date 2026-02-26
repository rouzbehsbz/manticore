gen-server-proto:
	@echo "Generating server protobuf ..."
	protoc \
    	--go_out=. \
    	--go_opt=paths=source_relative \
    	shared/packets.proto
	mv shared/packets.pb.go server/pkg/network/protocol/packets.pb.go
