package lobo

import (
	"encoding/binary"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
)

const RPCHeaderLen = 4

type RPCProtocol struct{}

func (d *RPCProtocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if buffer.VirtualLength() > RPCHeaderLen {
		buf := pbytes.GetLen(RPCHeaderLen)
		defer pbytes.Put(buf)
		_, _ = buffer.VirtualRead(buf)
		dataLen := binary.BigEndian.Uint32(buf)

		if buffer.VirtualLength() >= int(dataLen) {
			ret := make([]byte, dataLen)
			_, _ = buffer.VirtualRead(ret)

			buffer.VirtualFlush()
			return nil, ret
		} else {
			buffer.VirtualRevert()
		}
	}
	return nil, nil
}

func (d *RPCProtocol) Packet(c *gev.Connection, data interface{}) []byte {
	dd := data.([]byte)
	dataLen := len(dd)
	ret := make([]byte, RPCHeaderLen+dataLen)
	binary.BigEndian.PutUint32(ret, uint32(dataLen))
	copy(ret[4:], dd)
	return ret
}
