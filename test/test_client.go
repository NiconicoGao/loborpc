package test

import (
	"abc/lobo"
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"
)

type UserServiceClient struct {
	lobo.RPCClient
	once sync.Once
	etcd *lobo.ServiceDiscovery
}

func NewUserServiceClient(endpoints []string) (*UserServiceClient, error) {
	etcd := lobo.NewServiceDiscovery(endpoints)

	return &UserServiceClient{
		etcd: etcd,
		once: sync.Once{},
	}, nil

}

func (s *UserServiceClient) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	var err error
	s.once.Do(
		func() {
			err = s.etcd.WatchService("node")
		},
	)

	if err != nil {
		return nil, err
	}

	addrList := s.etcd.GetServices()
	if len(addrList) == 0 {
		return nil, errors.New("no service found")
	}

	err = s.Init(s.ChooseService(addrList))
	if err != nil {
		return nil, err
	}

	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	out := new(lobo.RPCRequest)
	out.ServerName = "UserService"
	out.MethodName = "Login"
	out.Message = data

	outStream, err := proto.Marshal(out)
	if err != nil {
		return nil, err
	}

	buffer := s.Packet(outStream)
	_, err = s.Conn.Write(buffer)
	if err != nil {
		return nil, err
	}

	msg, err := s.UnPacket(s.Conn)
	if err != nil {
		return nil, err
	}

	back := new(lobo.RPCResponce)
	err = proto.Unmarshal(msg, back)
	if err != nil {
		return nil, err
	}

	if back.Code != 0 {
		return nil, fmt.Errorf("Error code %v, %v", back.Code, back.Msg)
	}

	resp := new(LoginResponse)
	err = proto.Unmarshal(back.Resp, resp)
	if err != nil {
		return nil, err
	}
	s.RPCClient.Close()

	return resp, nil
}

func (s *UserServiceClient) Close() {
	s.RPCClient.Close()
	s.etcd.Close()
}
