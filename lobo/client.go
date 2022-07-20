package lobo

import (
	"encoding/binary"
	"io"
	"math/rand"
	"net"
	"time"
)

type RPCClient struct {
	Addr   string
	Conn   net.Conn
	Logger *Logger
}

func NewRPCClient(addr string) (*RPCClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &RPCClient{
		Addr:   addr,
		Conn:   conn,
		Logger: NewLogger(),
	}, nil
}

func (s *RPCClient) Init(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	s.Addr = addr
	s.Conn = conn
	s.Logger = NewLogger()
	return nil
}

func (s *RPCClient) Close() {
	s.Conn.Close()
}

func (*RPCClient) Packet(data []byte) []byte {
	buffer := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buffer, uint32(len(data)))
	copy(buffer[4:], data)
	return buffer
}

func (*RPCClient) UnPacket(c net.Conn) ([]byte, error) {
	var header = make([]byte, 4)

	_, err := io.ReadFull(c, header)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header)
	contentByte := make([]byte, length)
	_, e := io.ReadFull(c, contentByte)
	if e != nil {
		return nil, e
	}

	return contentByte, nil
}

func (s *RPCClient) ChooseService(l []string) string {
	if len(l) == 0 {
		return ""
	}

	rand.Seed(time.Now().Unix())
	str := l[rand.Int()%len(l)]
	s.Logger.Info("Choosing %v server", str)
	return str
}
