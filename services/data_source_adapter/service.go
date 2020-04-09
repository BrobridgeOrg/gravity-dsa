package runner_supervisor

import (
	"time"

	"github.com/flyaways/pool"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	app "gravity-dsa/app/interface"
	pb "gravity-dsa/pb"
)

type Service struct {
	app      app.AppImpl
	grpcPool *pool.GRPCPool
}

func CreateService(a app.AppImpl) *Service {

	address := viper.GetString("data_handler.host")

	options := &pool.Options{
		InitTargets:  []string{address},
		InitCap:      5,
		MaxCap:       30,
		DialTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	// Initialize connection pool
	p, err := pool.NewGRPCPool(options, grpc.WithInsecure())

	if err != nil {
		log.Printf("%#v\n", err)
		return nil
	}

	if p == nil {
		log.Printf("p= %#v\n", p)
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

	// TODO: Getting URL for data handler to process event

	log.WithFields(log.Fields{
		"event": in.EventName,
	}).Info("Received event")

	// Getting connection from pool
	conn, err := service.grpcPool.Get()
	if err != nil {
		log.Error("Failed to get connection: %v", err)
		return &pb.PublishReply{
			Success: false,
			Reason:  "Cannot connect to data handler",
		}, nil
	}
	defer service.grpcPool.Put(conn)

	// Preparing context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.PushRequest{
		EventName: in.EventName,
		Payload:   in.Payload,
	}

	// Push message to data handler
	res, err := pb.NewDataHandlerClient(conn).Push(ctx, req)
	if err != nil {
		log.Error(err)
		return &pb.PublishReply{
			Success: false,
			Reason:  "Data handler cannot handle this event",
		}, nil
	}

	if res.Success == false {
		return &pb.PublishReply{
			Success: false,
			Reason:  res.Reason,
		}, nil
	}

	return &pb.PublishReply{
		Success: true,
	}, nil
}
