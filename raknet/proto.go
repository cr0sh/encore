package raknet

import (
	"github.com/cr0sh/encore/util/binary"
)

// Packet ID: 0x00
type ConnectedPing struct {
	SendPingTime uint64
}

// Packet ID: 0x01
type UnconnectedPing struct {
	PingID     uint64
	OfflineMsg offlineMessageDataID
}

// Packet ID: 0x03
type ConnectedPong struct {
	SendPingTime uint64
	SendPongTime uint64
}

// Packet ID: 0x1c
type UnconnectedPong struct {
	PingID     uint64
	ServerID   uint64
	OfflineMsg offlineMessageDataID
	ServerName binary.FixedMCString
}

// Packet ID: 0x05
type OpenConnectionRequest1 struct {
	OfflineMsg   offlineMessageDataID
	ProtoVersion byte
}

// Packet ID: 0x06
type OpenConnectionReply1 struct {
	OfflineMsg offlineMessageDataID
	ServerGUID uint64
	Security   bool
	MTU        uint16
}

// Packet ID: 0x07
type OpenConnectionRequest2 struct {
	OfflineMsg offlineMessageDataID
	RemoteAddr IPAddr
	MTU        uint16
	ClientGUID uint64
}

// Packet ID: 0x08
type OpenConnectionRePly2 struct {
	OfflineMsg offlineMessageDataID
	ServerGUID uint64
	ClientAddr IPAddr
	MTU        uint16
	Security   bool
}

// Packet ID: 0x09
type ConnectRequest struct {
	ClientGUID   uint64
	SendPingTime uint64
	Security     bool
}

// Packet ID: 0x10
type ServerHandshake struct {
	SystemAddr   IPAddr
	SystemIdx    uint16
	SystemAddrs  systemAddresses
	SendPingTime uint64
	SendPongTime uint64
}

// Packet ID: 0x13
type ClientHandshake struct {
	ClientAddr   IPAddr
	SystemAddrs  systemAddresses
	SendPingTime uint64
	SendPongTime uint64
}

// Packet ID: 0x15
type ClientDisconnect struct{}
