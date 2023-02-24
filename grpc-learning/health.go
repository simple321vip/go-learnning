package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"time"
)

func GrpcHealthCheck() {
	s := grpc.NewServer()
	healthServer := health.NewServer() // healthServer 算是service更好一些
	healthpb.RegisterHealthServer(s, healthServer)

	ls, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		log.Println("can not to listen port 8080")
	}
	healthServer.SetServingStatus("grpc.health.v1.Health", healthpb.HealthCheckResponse_UNKNOWN)
	go func() {
		err = s.Serve(ls)
		if err != nil {
			log.Println("Server started error")
		}
	}()

	go func() {
		conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Println("can not connect grpc server")
		}

		cli := healthpb.NewHealthClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		c, err := cli.Watch(ctx, &healthpb.HealthCheckRequest{
			Service: "grpc.health.v1.Health",
		})

		if err != nil {
			log.Println("can not watch grpc server")
		}

		resp, err := c.Recv()

		if err != nil {
			log.Println("can not read response message")
		} else {
			log.Println(resp.Status)
		}
	}()

	// 保证主线程最后结束
	time.Sleep(2 * time.Second)
}
