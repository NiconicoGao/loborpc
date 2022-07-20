package lobo

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"sync"

	"github.com/Allenxuxu/gev"
)

type MethodInfo struct {
	Req  sync.Pool
	Resp sync.Pool
	Name string
}

type ServerImpl interface {
	GetMethodList() map[string]*MethodInfo
	GetRPCName() string
}

type ServerInfo struct {
	Methods map[string]*MethodInfo
	r       reflect.Value
}

type RPCServer struct {
	register map[string]*ServerInfo
	Logger   *Logger
}

func NewServer() *RPCServer {
	s := new(RPCServer)
	s.register = make(map[string]*ServerInfo)
	s.Logger = NewLogger()
	return s
}

func (rpc *RPCServer) Register(s ServerImpl) {
	info := new(ServerInfo)
	info.Methods = s.GetMethodList()
	info.r = reflect.ValueOf(s).Elem()
	rpc.register[s.GetRPCName()] = info
}

func (rpc *RPCServer) Serve(port int, endpoints []string) error {
	ser, err := NewServiceRegister(endpoints, fmt.Sprintf("node%v", port), fmt.Sprintf("127.0.0.1:%v", port), 10)
	if err != nil {
		panic(err)
	}
	go ser.ListenLeaseRespChan()
	handler := new(RPCEvent)
	handler.server = rpc
	s, err := gev.NewServer(handler,
		gev.Network("tcp"),
		gev.Address(":"+strconv.Itoa(port)),
		gev.NumLoops(runtime.NumCPU()),
		gev.CustomProtocol(&RPCProtocol{}),
	)
	if err != nil {
		return err
	}

	rpc.Logger.Info("Starting server at port %v\n", port)
	s.Start()
	return nil
}
