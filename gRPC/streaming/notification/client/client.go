package main

import (
	"context"
	"io"
	"log"
	"notify/proto/notificationpb"

	"google.golang.org/grpc"
)

func main() {
	conn, _ := grpc.Dial("localhost:50052", grpc.WithInsecure())
	defer conn.Close()

	client := notificationpb.NewNotificationServiceClient(conn)

	stream, err := client.StreamNotifications(context.Background(), &notificationpb.NotificationRequest{User: "Ananda"})
	if err != nil {
		log.Fatal(err)
	}

	for {
		notif, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(" %s: %s", notif.Title, notif.Body)
	}
}
