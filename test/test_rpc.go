package test

import (
	"abc/lobo"
	"context"
	"errors"
	"sync"
)

type UserService interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
}

type UnimplementedUserService struct{}

func (UnimplementedUserService) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, errors.New("method not implemented")
}

func (UnimplementedUserService) GetRPCName() string {
	return "UserService"
}

func (UnimplementedUserService) GetMethodList() map[string]*lobo.MethodInfo {
	m := make(map[string]*lobo.MethodInfo)
	info := new(lobo.MethodInfo)
	info.Name = "Login"
	info.Req = sync.Pool{
		New: func() interface{} {
			return new(LoginRequest)
		},
	}
	info.Resp = sync.Pool{
		New: func() interface{} {
			return new(LoginResponse)
		},
	}
	m[info.Name] = info

	return m
}
