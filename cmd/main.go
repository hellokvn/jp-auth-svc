package main

import (
	"fmt"

	"github.com/hellokvn/jp-auth-svc/pkg/config"
)

func main() {
	config, err := config.LoadConfig()
	fmt.Println(config.Port)
	fmt.Println(err)
	// viper.SetConfigFile("./pkg/envs/.env")
	// viper.ReadInConfig()

	// port := viper.Get("PORT").(string)
	// dbUrl := viper.Get("DB_URL").(string)
	// secretKey := viper.Get("JWT_SECRET_KEY").(string)

	// h := db.Init(dbUrl)

	// jwt := utils.JwtWrapper{
	// 	SecretKey:       secretKey,
	// 	Issuer:          "jp-auth-svc",
	// 	ExpirationHours: 24 * 365,
	// }

	// lis, err := net.Listen("tcp", port)

	// if err != nil {
	// 	log.Fatalln("Failed to listing:", err)
	// }

	// s := services.Server{
	// 	H:   h,
	// 	Jwt: jwt,
	// }

	// grpcServer := grpc.NewServer()

	// pb.RegisterAuthServiceServer(grpcServer, &s)

	// if err := grpcServer.Serve(lis); err != nil {
	// 	log.Fatalln("Failed to serve:", err)
	// }
}
