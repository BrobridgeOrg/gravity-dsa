module gravity-dsa

go 1.13

require (
	github.com/BrobridgeOrg/gravity-api v0.0.0-00010101000000-000000000000
	github.com/flyaways/pool v1.0.1
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	google.golang.org/grpc v1.31.0
)

//replace github.com/BrobridgeOrg/gravity-api => ../gravity-api
