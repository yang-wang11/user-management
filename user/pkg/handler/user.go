package handler

import (
	"context"
	"gorm.io/gorm"
	"user/pkg/pb"
	"user/pkg/repository"
)

type ReturnCode uint32

const (
	HTTPSuccess ReturnCode = 0
	HTTPFail    ReturnCode = 1
)

func (c ReturnCode) GetUint32() uint32 {
	return uint32(c)
}

type UserService struct {
	//service.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) UserLogin(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user := repository.NewUser(u.db)
	resp := new(pb.UserResponse)
	resp.ReturnCode = HTTPSuccess.GetUint32()
	user, err := user.ShowUserInfo(req)
	if err != nil {
		resp.ReturnCode = HTTPFail.GetUint32()
		return resp, err
	}

	resp.UserDetail = user.Translate()
	return resp, nil
}

func (u *UserService) UserRegister(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user := repository.NewUser(u.db)
	resp := new(pb.UserResponse)
	resp.ReturnCode = HTTPSuccess.GetUint32()
	if err := user.CreateUser(req); err != nil {
		resp.ReturnCode = HTTPFail.GetUint32()
		return resp, err
	}

	resp.UserDetail = user.Translate()
	return resp, nil
}
