package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	proto "simpleGuide/grpc"
	"strconv"
	"time"
)

type ClientStream struct {
	name string
	clientID int
	stream *proto.TimeAsk_ConnectToServerServer
	chQuit chan int
}

// Struct that will be used to represent the Server.
type Server struct {
	proto.UnimplementedTimeAskServer // Necessary
	name                             string
	port                             int
	clients                          []*ClientStream
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 5454, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
		clients: []*ClientStream{},
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterTimeAskServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

/*func (c *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	log.Printf("Client with ID %d asked for the time\n", in.ClientId)
	return &proto.TimeMessage{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil
}*/

func (c *Server) SendMessage(ctx context.Context, msg *proto.ClientPublishMessage) (*proto.ServerPublishMessageOk, error) {
	log.Printf("Received a message '%v'\n", msg)

	if msg.Message == "quit" {
		for i, client := range c.clients {
			if msg.ClientId == int64(client.clientID) {
				//removes the client who quited - the same as pop. ... means we want to add to our array
				c.clients = append(c.clients[:i], c.clients[i+1:]...)
				//sends a message that we want to break the connection to the server
				client.chQuit <- 0
				break
			}
		}
	} else {
		//for each client send them a message
		for _, client := range c.clients {
			log.Printf("Sending the message to client %d\n", client.clientID)
			(*client.stream).Send(&proto.MessageStreamConnection {
				StreamMessage: msg.Message,
			})
		}
	}


	return &proto.ServerPublishMessageOk{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil
}

func (c *Server) ConnectToServer(msg *proto.ClientConnectMessage, stream proto.TimeAsk_ConnectToServerServer) error {
	//log.Printf("%s connected to the server at time %d\n", msg.Name, ???) //sæt tiden på
	log.Printf("%s connected to the server", msg.Name)
	clientStream := &ClientStream {
		name: msg.Name,
		clientID: int(msg.ClientId),
		stream: &stream,
		chQuit: make(chan int),
	}

	//saves the clients that are connected to the server
	c.clients = append(c.clients, clientStream)

	log.Printf("There are %d clients connected", len(c.clients))

	<-clientStream.chQuit

	return nil
}



