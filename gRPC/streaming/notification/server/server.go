package main

import (
	"fmt"
	"log"
	"net"
	"notify/proto/notificationpb"
	"time"

	"google.golang.org/grpc"
)

type notificationServer struct {
	notificationpb.UnimplementedNotificationServiceServer
}

func (s *notificationServer) StreamNotifications(req *notificationpb.NotificationRequest, stream notificationpb.NotificationService_StreamNotificationsServer) error {
	user := req.User
	log.Printf("Sending notifications to user: %s", user)

	for i := 1; i <= 5; i++ {
		notif := &notificationpb.Notification{
			Title: fmt.Sprintf("Alert #%d", i),
			Body:  fmt.Sprintf("Hello %s! This is notification #%d", user, i),
		}
		if err := stream.Send(notif); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func main() {
	lis, _ := net.Listen("tcp", ":50052")
	grpcServer := grpc.NewServer()
	notificationpb.RegisterNotificationServiceServer(grpcServer, &notificationServer{})
	log.Println("Notification Server running on :50052")
	log.Fatal(grpcServer.Serve(lis))
}
