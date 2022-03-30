package main

import (
	"fmt"
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
	c, err := config.LoadConfig()

	h := db.Init(c.DBUrl)
	d := sender.Init(c.MailSvcAMQPUrl)

	j := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "jp-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	s := services.Server{
		H:       h,
		Jwt:     j,
		MailSvc: d,
	}

	fmt.Println("Listen on:", c.Port)

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
