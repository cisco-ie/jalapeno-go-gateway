module gateway

go 1.13

require (
	cloud.google.com/go v0.52.0 // indirect
	github.com/cisco-ie/jalapeno-go-gateway/pkg/gateway v0.0.0-00010101000000-000000000000
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	google.golang.org/grpc v1.26.0
)

replace (
	github.com/cisco-ie/jalapeno-go-gateway/pkg/apis => ./pkg/apis
	github.com/cisco-ie/jalapeno-go-gateway/pkg/gateway => ./pkg/gateway
)
