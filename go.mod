module gravity-dsa

go 1.13

require (
	github.com/BrobridgeOrg/gravity-adapter-nats v0.0.0-20201013172819-41aa2d00d5ab // indirect
	github.com/BrobridgeOrg/gravity-api v0.2.1
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/cfsghost/grpc-connection-pool v0.3.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.1.4 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	//	google.golang.org/grpc v1.31.1
	google.golang.org/grpc v1.32.0
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api

//replace github.com/cfsghost/grpc-connection-pool => /Users/fred/works/opensource/grpc-connection-pool
