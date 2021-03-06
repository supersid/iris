package broker

import (
	"fmt"
	"github.com/supersid/iris/service"
	"github.com/supersid/iris/client"
	"github.com/supersid/iris/message"
	"github.com/satori/go.uuid"
	"github.com/supersid/iris/worker"
)

/*
Data obtained when listening on a ZMQ socket
It's just an slice of strings.
*/

const WORKER_READY string = "WORKER_READY"

func (broker *Broker) ParseMessage(msg []string) message.Message {
	var m message.Message
	for index, message := range msg {
		fmt.Println(fmt.Sprintf("%d. %s", index, message))
	}

	if len(msg) < 4 {
		return message.Message{}
	}

	command := msg[2]
	sender  := msg[0]

	if command == WORKER_READY || command == client.CLIENT_REQUEST{
		m = message.Message{
			Sender:       sender,
			Command:      command,
			Identity:     msg[4],
			ServiceName:  msg[3],
			Data:         msg[5],
			MessageId:    uuid.NewV4().String(),
		}
	} else if command == worker.WORKER_RESPONSE {
		fmt.Println("WORKER RESPONSE arrived. Yaay!!")
		m = message.Message{
			Sender:       msg[1],
			Command:      command,
			Data:  	      msg[3],
			ResponseData: msg[4],
			Identity:     msg[5],
			MessageId:    uuid.NewV4().String(),
		}
	} else {
		m = message.Message{}
	}

	return m
}

func (broker *Broker) ProcessMessage(msg message.Message) {
	if msg.Command == WORKER_READY {
		broker.WorkerReadyHandler(msg)
	} else if msg.Command == client.CLIENT_REQUEST {
		fmt.Println("Client Request arrived.")
		broker.ClientRequestHandler(msg)
	} else if msg.Command == worker.WORKER_RESPONSE {
		fmt.Println("Worker response arrived.")
		broker.WorkerResponseHandler(msg)
	}

}

func (broker *Broker) FindOrCreateService(serviceName string) (*service.Service, bool) {
	srvc, present := broker.services[serviceName]

	if !present {
		srvc = service.NewService(serviceName)
	}

	return srvc, present
}
