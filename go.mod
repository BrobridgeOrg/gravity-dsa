module gravity-dsa

go 1.13

require (
	github.com/BrobridgeOrg/gravity-api v0.2.0
	github.com/cfsghost/grpc-connection-pool v0.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.2
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	google.golang.org/grpc v1.31.1
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api
//replace github.com/cfsghost/grpc-connection-pool => /Users/fred/works/opensource/grpc-connection-pool
