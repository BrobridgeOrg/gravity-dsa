package data_source_adapter

import (
	"errors"
	"sync"
	"time"

	data_handler "github.com/BrobridgeOrg/gravity-api/service/data_handler"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var requestPool = sync.Pool{
	New: func() interface{} {
		return &data_handler.PushRequest{}
	},
}

func (service *Service) publishAsync(eventName string, payload []byte) error {

	// Prepare request
	req := requestPool.Get().(*data_handler.PushRequest)
	req.EventName = eventName
	req.Payload = payload

	service.incoming <- req

	return nil
}

func (service *Service) publish(eventName string, payload []byte) error {

	// Prepare request
	req := requestPool.Get().(*data_handler.PushRequest)
	req.EventName = eventName
	req.Payload = payload

	err := service.sendData(req)
	requestPool.Put(req)

	return err
}

func (service *Service) sendData(req *data_handler.PushRequest) error {

	// Getting connection from pool
	conn, err := service.grpcPool.Get()
	if err != nil {
		log.Error("Failed to get connection: %v", err)
		return errors.New("Cannot connect to data handler")
	}

	// Preparing context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Push message to data handler
	res, err := data_handler.NewDataHandlerClient(conn).Push(ctx, req)
	if err != nil {
		log.Error(err)
		return errors.New("Data handler cannot handle this event")
	}

	if res.Success == false {
		return errors.New(res.Reason)
	}

	return nil
}

func (service *Service) startWorker() {

	for {
		select {
		case req := <-service.incoming:
			err := service.sendData(req)
			requestPool.Put(req)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
