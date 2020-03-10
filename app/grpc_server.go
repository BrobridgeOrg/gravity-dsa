package app

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	app "gravity-dsa/app/interface"
	pb "gravity-dsa/pb"
	data_source_adapter "gravity-dsa/services/data_source_adapter"
)

func (a *App) InitGRPCServer(host string) error {

	// Start to listen on port
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.WithFields(log.Fields{
		"host": host,
	}).Info("Starting gRPC server on " + host)

	// Create gRPC server
	s := grpc.NewServer()

	// Register data source adapter service
	dsaService := data_source_adapter.CreateService(app.AppImpl(a))
	pb.RegisterDataSourceAdapterServer(s, dsaService)
	reflection.Register(s)

	log.WithFields(log.Fields{
		"service": "DataSourceAdapter",
	}).Info("Registered service")

	// Starting server
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
