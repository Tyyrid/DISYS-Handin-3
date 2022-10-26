package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	proto "DISYS-HANDIN-3/grpc"
	"strconv"
	"time"
)

type Server struct {
	proto.UnimplementedTimeAskServer // Necessary
	name                             string
	port                             int

	mutex          sync.Mutex // used to lock the server to avoid race conditions.
}

var serverName = flag.String("name", "default", "Senders name") // set with "-name <name>" in terminal
var port = flag.Int("port", 8000, "server port number")


func main() {
	// Get the port from the command line when the server is run
	flag.Parse()
	fmt.Println(".:server is starting:.")

	// Create a server struct
	server := &Server{
		name: *serverName,
		port: *port,
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
}

func startServer(server *Server) {
	log.Printf("Server %s: Attempts to create listener on port %d\n", *serverName, *port)

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	//listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Server %s: Failed to listen on port %d: %v", *serverName, *port, err)
		//If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}
	log.Printf("Started server %s at port: %d\n", *serverName, *port)

	// Register the grpc server and serve its listener
	proto.RegisterTimeAskServer(grpcServer, server)

	log.Printf("Server %s: Listening at %v\n", *serverName, listener.Addr())

	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener: %v", serveError)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
	
}