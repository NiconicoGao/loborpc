package lobo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Allenxuxu/gev"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type RPCEvent struct {
	server *RPCServer
}

func (s *RPCEvent) OnConnect(c *gev.Connection) {
	fmt.Println(" OnConnect ： ", c.PeerAddr())
}

func (s *RPCEvent) BadRequest(code int32, message string) []byte {
	resp := new(RPCResponce)
	resp.Code = code
	resp.Msg = message
	resp.Resp = make([]byte, 0)
	out, err := proto.Marshal(resp)
	if err != nil {
		fmt.Printf("Marshal error\n")
	}
	return out
}

func (s *RPCEvent) SuccessRequest(data []byte) []byte {
	resp := new(RPCResponce)
	resp.Code = 0
	resp.Msg = ""
	resp.Resp = data
	out, err := proto.Marshal(resp)
	if err != nil {
		fmt.Printf("Marshal error\n")
	}
	return out
}

func (s *RPCEvent) OnMessage(c *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	var err error
	req := new(RPCRequest)

	if err = proto.Unmarshal(data, req); err != nil {
		return s.BadRequest(401, "Bad Request")
	}

	serverInfo := s.server.register[req.ServerName]
	if serverInfo == nil {
		return s.BadRequest(404, "Server Not Found")
	}

	methodInfo := serverInfo.Methods[req.MethodName]
	if methodInfo == nil {
		return s.BadRequest(404, "Method Not Found")
	}

	fmt.Printf("Get RPC Request of Server %v Method %v\n", req.ServerName, req.MethodName)

	userReq := methodInfo.Req.Get().(protoreflect.ProtoMessage)
	if err := proto.Unmarshal(req.Message, userReq); err != nil {
		return s.BadRequest(402, "Bad Request")
	}

	method := serverInfo.r.MethodByName(req.MethodName)
	userOut := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(userReq)})
	if len(userOut) != 2 {
		return s.BadRequest(404, "Method Not Found")
	}

	if !userOut[1].IsNil() {
		return s.BadRequest(200, err.Error())
	}

	userResp := userOut[0].Interface().(proto.Message)
	origin, err := proto.Marshal(userResp)
	if err != nil {
		return s.BadRequest(405, "Marshal Error")
	}

	return s.SuccessRequest(origin)

}

func (s *RPCEvent) OnClose(c *gev.Connection) {

}
