package runner_supervisor

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	app "gravity-dsa/app/interface"
	pb "gravity-dsa/pb"
)

type Service struct {
	app app.AppImpl
}

func CreateService(a app.AppImpl) *Service {

	// Preparing service
	service := &Service{
		app: a,
	}

	return service
}

func (service *Service) Publish(ctx context.Context, in *pb.PublishRequest) (*pb.PublishReply, error) {

	// TODO: Getting URL for data handler to process event

	log.WithFields(log.Fields{
		"event": in.EventName,
	}).Info("Received event")

	// Set up a connection to the server.
	address := viper.GetString("data_handler.host")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error("did not connect: %v", err)
		return &pb.PublishReply{
			Success: false,
			Reason:  "Cannot connect to data handler",
		}, nil
	}

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
