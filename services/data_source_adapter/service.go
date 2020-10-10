package runner_supervisor

import (
	"sync"
	"time"

	grpc_connection_pool "github.com/cfsghost/grpc-connection-pool"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	app "gravity-dsa/app/interface"

	data_handler "github.com/BrobridgeOrg/gravity-api/service/data_handler"
	pb "github.com/BrobridgeOrg/gravity-api/service/dsa"
)

//var counter uint64

var PublishSuccess = pb.PublishReply{
	Success: true,
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &data_handler.PushRequest{}
	},
}

type Service struct {
	app      app.AppImpl
	grpcPool *grpc_connection_pool.GRPCPool
}

func CreateService(a app.AppImpl) *Service {

	address := viper.GetString("data_handler.host")

	options := &grpc_connection_pool.Options{
		InitCap:     8,
		MaxCap:      16,
		DialTimeout: time.Second * 20,
	}

	// Initialize connection pool
	p, err := grpc_connection_pool.NewGRPCPool(address, options, grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return nil
	}

	if p == nil {
		log.Error(err)
		return nil
	}

	// Preparing service
	service := &Service{
		app:      a,
		grpcPool: p,
	}

	return service
}

func (service *Service) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishReply, error) {
	/*
		id := atomic.AddUint64((*uint64)(&counter), 1)

		if id%1000 == 0 {
			log.Info(id)
		}
	*/
	/*
		log.WithFields(log.Fields{
			"event": in.EventName,
		}).Info("Received event")
	*/
	// Getting connection from pool
	conn, err := service.grpcPool.Get()
	if err != nil {
		log.Error("Failed to get connection: %v", err)
		return &pb.PublishReply{
			Success: false,
			Reason:  "Cannot connect to data handler",
		}, nil
	}

	// Preparing context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Prepare request
	req := requestPool.Get().(*data_handler.PushRequest)
	req.EventName = in.EventName
	req.Payload = in.Payload

	// Push message to data handler
	res, err := data_handler.NewDataHandlerClient(conn).Push(ctx, req)
	if err != nil {
		log.Error(err)
		requestPool.Put(req)
		return &pb.PublishReply{
			Success: false,
			Reason:  "Data handler cannot handle this event",
		}, nil
	}

	requestPool.Put(req)

	if res.Success == false {
		return &pb.PublishReply{
			Success: false,
			Reason:  res.Reason,
		}, nil
	}

	return &PublishSuccess, nil
}
