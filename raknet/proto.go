package raknet

import (
	"github.com/cr0sh/encore/util/binary"
)

// Packet ID: 0x00
type ConnectedPing struct {
	SendPingTime int64
}

func (*ConnectedPing) ID() byte {
	return 0x00
}

// Packet ID: 0x01
type UnconnectedPing struct {
	PingID     uint64
	OfflineMsg offlineMessageDataID
}

func (*UnconnectedPing) ID() byte {
	return 0x01
}

// Packet ID: 0x03
type ConnectedPong struct {
	SendPingTime int64
	SendPongTime int64
}

func (*ConnectedPong) ID() byte {
	return 0x03
}

// Packet ID: 0x1c
type UnconnectedPong struct {
	PingID     uint64
	ServerID   uint64
	OfflineMsg offlineMessageDataID
	ServerName binary.FixedMCString
}

func (*UnconnectedPong) ID() byte {
	return 0x1c
}

// Packet ID: 0x05
type OpenConnectionRequest1 struct {
	OfflineMsg   offlineMessageDataID
	ProtoVersion byte
}

func (*OpenConnectionRequest1) ID() byte {
	return 0x05
}

// Packet ID: 0x06
type OpenConnectionReply1 struct {
	OfflineMsg offlineMessageDataID
	ServerGUID uint64
	Security   bool
	MTU        uint16
}

func (*OpenConnectionReply1) ID() byte {
	return 0x06
}

// Packet ID: 0x07
type OpenConnectionRequest2 struct {
	OfflineMsg offlineMessageDataID
	RemoteAddr IPAddr
	MTU        uint16
	ClientGUID uint64
}

func (*OpenConnectionRequest2) ID() byte {
	return 0x07
}

// Packet ID: 0x08
type OpenConnectionReply2 struct {
	OfflineMsg offlineMessageDataID
	ServerGUID uint64
	ClientAddr IPAddr
	MTU        uint16
	Security   bool
}

func (*OpenConnectionReply2) ID() byte {
	return 0x08
}

// Packet ID: 0x09
type ConnectionRequest struct {
	ClientGUID   uint64
	SendPingTime int64
	Security     bool
}

func (*ConnectionRequest) ID() byte {
	return 0x09
}

// Packet ID: 0x10
type ServerHandshake struct {
	SystemAddr   IPAddr
	SystemIdx    uint16
	SystemAddrs  systemAddresses
	SendPingTime int64
	SendPongTime int64
}

func (*ServerHandshake) ID() byte {
	return 0x10
}

// Packet ID: 0x13
type ClientHandshake struct {
	ClientAddr   IPAddr
	SystemAddrs  systemAddresses
	SendPingTime int64
	SendPongTime int64
}

func (*ClientHandshake) ID() byte {
	return 0x13
}

// Packet ID: 0x15
type ClientDisconnect struct{}

func (*ClientDisconnect) ID() byte {
	return 0x15
}
