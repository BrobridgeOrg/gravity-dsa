package data_source_adapter

import (
	"errors"
	"sync"

	data_handler "github.com/BrobridgeOrg/gravity-api/service/data_handler"
	log "github.com/sirupsen/logrus"
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

	// Getting stream from pool
	err := service.grpcPool.GetStream("push", func(s interface{}) error {

		// Send request
		return s.(data_handler.DataHandler_PushStreamClient).Send(req)
	})
	if err != nil {
		log.Error("Failed to get connection: %v", err)
		return errors.New("Cannot connect to data handler")
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
