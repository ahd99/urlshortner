## proto compile command

protoc --go_out=internal/monitoring/proto --go_opt=paths=source_relative  --go-grpc_out=internal/monitoring/proto  --go-grpc_opt=paths=source_relative internal/monitoring/proto/monitoring.proto 

current directory is urlshortner project path (for example: ~/workspace/urlshortner)