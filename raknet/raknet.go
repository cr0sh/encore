// package raknet implements RakNet protocol for servers/clients.
// Its main usage and implementation goal is for servers,
// but it's designed to be used for both.
//
// Many protocol informations are from:
//  https://github.com/Nukkit/Nukkit
//  https://github.com/NiclasOlofsson/MiNET
package raknet

import (
	"bytes"
	"errors"
	"github.com/cr0sh/encore/util/binary"
	"net"
	"strconv"
)

const OFFLINE_MESSAGE_DATA_ID = "\x00\xff\xff\x00\xfe\xfe\xfe\xfe\xfd\xfd\xfd\xfd\x12\x34\x56\x78"

type offlineMessageDataID struct{}

func (offlineMessageDataID) MarshalStream(buf *bytes.Buffer) (err error) {
	_, err = buf.WriteString(OFFLINE_MESSAGE_DATA_ID)
	return
}

func (*offlineMessageDataID) UnmarshalStream(buf *bytes.Buffer) (err error) {
	buf.Next(len(OFFLINE_MESSAGE_DATA_ID))
	return nil
}

// IPAddr represents a single UDP endpoint address in raknet.
type IPAddr net.UDPAddr

// MarshalStream implements Stream Marshaler interface.
func (a IPAddr) MarshalStream(buf *bytes.Buffer) error {
	v4ip := a.IP.To4() // Currently only supports IPv4
	b := make([]byte, 7)
	b[0] = 4
	copy(b[1:5], v4ip[:4])
	binary.BigEndian.PutUint16(b[5:7], uint16(a.Port))
	_, err := buf.Write(b)
	return err
}

// UnmarshalStream implements Strream Unmarshaler interface.
func (a *IPAddr) UnmarshalStream(buf *bytes.Buffer) (err error) {
	b := make([]byte, 7)
	if _, err = buf.Read(b); err != nil {
		return
	}
	if b[0] != 4 {
		return errors.New("IPAddr only supports IPv4, v" + strconv.Itoa(int(b[0])) + " given")
	}

	a.IP = net.IPv4(^b[1], ^b[2], ^b[3], ^b[4])
	a.Port = int(binary.BigEndian.Uint16(b[5:7]))

	return nil
}

const systemAddressesReady = "\x04\x80\xff\xff\xfe\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00" +
	"\x04\xff\xff\xff\xff\x00\x00"

type systemAddresses struct{}

func (systemAddresses) MarshalStream(buf *bytes.Buffer) (err error) {
	_, err = buf.WriteString(systemAddressesReady)
	return
}

func (*systemAddresses) UnmarshalStream(buf *bytes.Buffer) error {
	buf.Next(len(systemAddressesReady))
	return nil
}
