package data_source_adapter

import (
	"io"
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

var counter uint64

var PublishSuccess = pb.PublishReply{
	Success: true,
}

type Service struct {
	app      app.AppImpl
	grpcPool *grpc_connection_pool.GRPCPool
	incoming chan *data_handler.PushRequest
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

	// Register initializer for stream
	p.SetStreamInitializer("push", func(conn *grpc.ClientConn) (interface{}, error) {
		client := data_handler.NewDataHandlerClient(conn)
		return client.PushStream(context.Background())
	})

	// Preparing service
	service := &Service{
		app:      a,
		grpcPool: p,
		incoming: make(chan *data_handler.PushRequest, 4096),
	}

	go service.startWorker()
	go service.startWorker()
	go service.startWorker()
	go service.startWorker()
	go service.startWorker()
	go service.startWorker()
	go service.startWorker()
	go service.startWorker()

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

	err := service.publish(in.EventName, in.Payload)
	if err != nil {
		log.Error(err)
		return &pb.PublishReply{
			Success: false,
			Reason:  err.Error(),
		}, nil
	}

	return &PublishSuccess, nil
}

func (service *Service) PublishEvents(stream pb.DataSourceAdapter_PublishEventsServer) error {

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		/*
			id := atomic.AddUint64((*uint64)(&counter), 1)

			if id%1000 == 0 {
				log.Info(id)
			}
		*/
		service.publishAsync(in.EventName, in.Payload)
	}
}
