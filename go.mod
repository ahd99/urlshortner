module github.com/ahd99/urlshortner

go 1.16

require (
	//mysql driver
	github.com/go-sql-driver/mysql v1.6.0 // indirect

	// prometheus
	github.com/prometheus/client_golang v1.11.0
	// ------

	// mongo
	go.mongodb.org/mongo-driver v1.7.1
	// -----

	// zap
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0
	// ------

	// grpc, protobuf
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
// ------

)
