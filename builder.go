package websocket

import (
	"github.com/ebar-go/websocket/epoll"
	cmap "github.com/orcaman/concurrent-map"
	"log"
)

type ServerBuilder struct {
	workerNumber int
	maxConnectionNumber int
	connectCallback Callback
	disconnectCallback Callback
}

func (builder *ServerBuilder) SetWorkerNumber(workerNumber int) *ServerBuilder {
	if workerNumber > 0 {
		builder.workerNumber = workerNumber
	}

	return builder
}

func (builder *ServerBuilder) SetMaxConnectionNumber(maxConnectionNumber int) *ServerBuilder{
	if maxConnectionNumber > 0 {
		builder.maxConnectionNumber = maxConnectionNumber
	}
	return builder
}

func (builder *ServerBuilder) SetConnectCallback(connectCallback Callback) *ServerBuilder{
	builder.connectCallback = connectCallback
	return builder
}

func (builder *ServerBuilder) SetDisconnectCallback(disconnectCallback Callback) *ServerBuilder{
	builder.disconnectCallback = disconnectCallback
	return builder
}

func (builder *ServerBuilder) Build() Server {
	e, err := epoll.Create()
	if err != nil {
		log.Fatalf("unable to create epoll:%v\n", err)
	}

	return &server{
		engine:      newEngine(),
		connections: cmap.New(),
		epoller:     e,
		connectCallback: builder.connectCallback,
		disconnectCallback: builder.disconnectCallback,
		workers:     newWorkerPool(builder.workerNumber, builder.maxConnectionNumber),
	}
}

func NewBuilder() *ServerBuilder {
	return &ServerBuilder{
		workerNumber:        50,
		maxConnectionNumber: 100000,
		connectCallback:      nil,
		disconnectCallback:   nil,
	}
}

