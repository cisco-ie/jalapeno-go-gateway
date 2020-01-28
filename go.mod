module github.com/cisco-ie/jalapeno-go-gateway/jalapeno-gateway

go 1.13

require (
	github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient v0.0.0-00010101000000-000000000000
	github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient v0.0.0-00010101000000-000000000000
	github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient/mock v0.0.0-00010101000000-000000000000
	github.com/cisco-ie/jalapeno-go-gateway/pkg/gateway v0.0.0-00010101000000-000000000000
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/osrg/gobgp v2.0.0+incompatible // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200113162924-86b910548bc1 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200115191322-ca5a22157cba // indirect
)

replace (
	github.com/cisco-ie/jalapeno-go-gateway/pkg/apis => ./pkg/apis
	github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient => ./pkg/bgpclient
	github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient => ./pkg/dbclient
	github.com/cisco-ie/jalapeno-go-gateway/pkg/dbclient/mock => ./pkg/dbclient/mock
	github.com/cisco-ie/jalapeno-go-gateway/pkg/gateway => ./pkg/gateway
)
