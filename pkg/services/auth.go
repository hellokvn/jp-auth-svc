package services

import (
	"context"
	"net/http"

	"github.com/hellokvn/jp-auth-svc/pkg/db"
	"github.com/hellokvn/jp-auth-svc/pkg/models"
	"github.com/hellokvn/jp-auth-svc/pkg/pb"
	"github.com/hellokvn/jp-auth-svc/pkg/sender"
	"github.com/hellokvn/jp-auth-svc/pkg/utils"
)

type Server struct {
	H       db.Handler
	MailSvc sender.Handler
	Jwt     utils.JwtWrapper
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var auth models.Auth

	if result := s.H.DB.Where(&models.Auth{Email: req.Email}).First(&auth); result.Error == nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	auth.Email = req.Email
	auth.Password = utils.HashPassword(req.Password)

	s.H.DB.Create(&auth)

	if result := s.H.DB.Create(&auth); result.Error == nil {
		return &pb.RegisterResponse{
			Status: http.StatusBadGateway,
			Error:  result.Error.Error(),
		}, nil
	}

	b := &sender.SendMailBody{
		Template: "register",
		To:       auth.Email,
	}

	s.MailSvc.SendMail(b)

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var auth models.Auth

	if result := s.H.DB.Where(&models.Auth{Email: req.Email}).First(&auth); result.Error != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, auth.Password)

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(auth)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var auth models.Auth

	if result := s.H.DB.Where(&models.Auth{Email: claims.Email}).First(&auth); result.Error != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		AuthId: auth.Id,
	}, nil
}
