package main

import (
	pb "abc/test"
	"context"
)

type UserServiceImpl struct {
	pb.UnimplementedUserService
}

func (UserServiceImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp := new(pb.LoginResponse)
	resp.Base = new(pb.BaseResp)
	resp.Success = (string(req.GetName()) == "123" && string(req.GetPwd()) == "abc")
	return resp, nil
}
