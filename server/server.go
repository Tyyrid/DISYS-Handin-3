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
	"fmt"
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
	LamportTimestamp                 int64
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
		LamportTimestamp: 1,
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

func (s *Server) SendMessage(ctx context.Context, msg *proto.ClientPublishMessage) (*proto.ServerPublishMessageOk, error) {
	log.Printf("Received a message from %d '%v' at time %d\n", msg.ClientId, msg, s.LamportTimestamp)

	//update LamportTime. compare servers time and the time for the msg we received
	if msg.LamportTimestamp > s.LamportTimestamp{
		s.LamportTimestamp = msg.LamportTimestamp + 1
	} else {
		s.LamportTimestamp+=1
	}

	log.Printf("This is the updated timestamp: %d", s.LamportTimestamp)

	if msg.Message == "quit" {
		for i, client := range s.clients {
			if msg.ClientId == int64(client.clientID) {
				//removes the client who quited - the same as pop. ... means we want to add to our array
				s.clients = append(s.clients[:i], s.clients[i+1:]...)
				//sends a message that we want to break the connection to the server
				client.chQuit <- 0
				s.SendToAllClients(fmt.Sprintf("%s left at %d", client.name, s.LamportTimestamp))
				break
			}
		}
	} else {
		s.SendToAllClients(msg.Message)
	}

	return &proto.ServerPublishMessageOk{
		Time:       time.Now().String(),
		ServerName: s.name,
	}, nil
}


func (s *Server) SendToAllClients(msg string) {
	//for each client send them a message
	for _, client := range s.clients {
		log.Printf("Sending the message to client %d\n", client.clientID)
		(*client.stream).Send(&proto.MessageStreamConnection {
			StreamMessage: msg,
			LamportTimestamp: s.LamportTimestamp,
		})
	}
}

func (s *Server) ConnectToServer(msg *proto.ClientConnectMessage, stream proto.TimeAsk_ConnectToServerServer) error {
	s.LamportTimestamp += 1
	log.Printf("%s connected to the server at time %d\n", msg.Name, s.LamportTimestamp)
	//log.Printf("%s connected to the server", msg.Name)
	clientStream := &ClientStream {
		name: msg.Name,
		clientID: int(msg.ClientId),
		stream: &stream,
		chQuit: make(chan int),
	}

	//saves the clients that are connected to the server
	s.clients = append(s.clients, clientStream)

	s.SendToAllClients(fmt.Sprintf("%s joined at %d", msg.Name, s.LamportTimestamp))

	log.Printf("There are %d clients connected", len(s.clients))

	//as long as there is no message, the channel will stay open
	<-clientStream.chQuit

	return nil
}



