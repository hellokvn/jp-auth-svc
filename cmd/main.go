package main

import (
	"log"
	"net"

	"github.com/hellokvn/jp-auth-svc/pkg/config"
	"github.com/hellokvn/jp-auth-svc/pkg/db"
	"github.com/hellokvn/jp-auth-svc/pkg/pb"
	"github.com/hellokvn/jp-auth-svc/pkg/sender"
	"github.com/hellokvn/jp-auth-svc/pkg/services"
	"github.com/hellokvn/jp-auth-svc/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	config, err := config.LoadConfig()

	h := db.Init(config.DBUrl)
	d := sender.Init(config.MailSvcAMQPUrl)

	j := utils.JwtWrapper{
		SecretKey:       config.JWTSecretKey,
		Issuer:          "jp-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", config.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	s := services.Server{
		H:       h,
		Jwt:     j,
		MailSvc: d,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
