package main

import (
	"io"
	"bufio"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	proto "simpleGuide/grpc"
	"strconv"
)

type Client struct {
	id         				int
	LamportTimestamp 		int64
	stream 					*proto.TimeAsk_ConnectToServerClient
}

// go run . -name Hannah. Command to connect to server via a chosen name.
var (
	name = flag.String("name", "<name>", "Name of this participant")
	serverPort = flag.Int("sPort", 5454, "server port number (should match the port used for the server)")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()


	// Connect to the server
	serverConnection, _ := connectToServer()
	stream, err := serverConnection.ConnectToServer(context.Background(), &proto.ClientConnectMessage {
		Name: *name,
		ClientId: int64(os.Getpid()),
	})
	if err != nil {
		log.Fatalf("Connection failed")
	}

	// Create a client
	client := &Client {
		id:         		1,
		LamportTimestamp: 	0,
		stream: 			&stream,
	}

	go client.listenForMessages()

	// Wait for input in the client terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		log.Printf("Client wants to send a message: %s\n", input)

		// Increase the Lamport time and send message to Server
		client.LamportTimestamp += 1
		serverConnection.SendMessage(context.Background(), &proto.ClientPublishMessage {
			ClientId: 			int64(os.Getpid()),
			Message: 			input,
			LamportTimestamp: 	client.LamportTimestamp,
		})
	}
}

func connectToServer() (proto.TimeAskClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewTimeAskClient(conn), nil
}

func (c *Client) listenForMessages() {
	//while loop runs forever
	for {
		//if the client sent 'quit', this will close the connection
		msg, err := (*c.stream).Recv()
		if err == io.EOF {
			log.Fatalf("Closed connection to server")
		}
		if err != nil {
			log.Fatalf("There was some error: %v", err)
		}


		if msg.LamportTimestamp > c.LamportTimestamp {
			c.LamportTimestamp = msg.LamportTimestamp + 1
		} else {
			c.LamportTimestamp += 1
		}
		
		//"%v" print as a string
		log.Printf("Client has received message '%s' at Client Lamport time %d", msg.StreamMessage, c.LamportTimestamp)
	}
}