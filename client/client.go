package client

import (
	uuid "github.com/satori/go.uuid"
	zmq "github.com/pebbe/zmq4"
	"errors"
	"fmt"
	"time"
)

const (
	CLIENT_REQUEST = "CLIENT_REQUEST"
)

type Client struct {
	socket *zmq.Socket
	brokerUrl  string
	identity   string
	poller     *zmq.Poller
}

func (client *Client) setIdentity() {
	client.identity = uuid.NewV4().String()
}

func newClient(brokerUrl string) (*Client, error) {
	client := &Client{
		brokerUrl: brokerUrl,
	}

	socket, err := zmq.NewSocket(zmq.DEALER)

	if err != nil {
		fmt.Println("[ERROR]: Could not create a new client socket.")
		return nil, err
	}

	client.socket = socket
	client.setIdentity()

	client.poller = zmq.NewPoller()
	client.poller.Add(client.socket, zmq.POLLIN)

	return client, err
}

func (client *Client) Close() {
	client.socket.Close()
}

func (client *Client) ReceiveMessage() (error, []string){
	sockets, err := client.poller.Poll(300 * time.Millisecond)

	if len(sockets) == 0 {
		return err, make([]string, 0)
	}
	if err != nil {
		return err, make([]string,0)
	}

	fmt.Println("Waiting on a response.")
	msg, err2 := client.socket.RecvMessage(0)
	if err2 != nil {
		return err2, make([]string, 0)
	}

	return nil, msg
}

/*
   CLIENT_REQUEST
   0: Blank Frame
   1: Command
   2: Blank Frame
   2: Service Name
   3: Data
*/
func (client *Client) createMessage(serviceName, message string) []string {
	msg := make([]string, 5)
	msg[0] = ""
	msg[1] = CLIENT_REQUEST
	msg[2] = serviceName
	msg[3] = client.identity
	msg[4] = message

	return msg
 }

func (client *Client) SendMessage(serviceName, message string) error{
	if serviceName == "" {
		return errors.New("service_name cannot be nil or blank.")
	}

	msg := client.createMessage(serviceName, message)
	fmt.Println(fmt.Sprintf("Sending %s: to service %s", message, serviceName))

	_, err := client.socket.SendMessage(msg)

	return err
}

func Start(brokerUrl string) *Client{
	client, err := newClient(brokerUrl)

	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR]: Client creation error %s", err.Error()))
		panic(err)
	}

	err = client.socket.Connect(brokerUrl)

	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR]: Client could not connect to broker. %s", err.Error()))
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Starting client id %s by connecting to %s", client.identity, client.brokerUrl))
	return client
}
