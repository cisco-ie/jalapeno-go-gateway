module github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient

go 1.13

require (
	github.com/cisco-ie/jalapeno-go-gateway/pkg/apis v0.0.0-00010101000000-000000000000
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	google.golang.org/grpc v1.26.0
)

replace github.com/cisco-ie/jalapeno-go-gateway/pkg/apis => ../apis