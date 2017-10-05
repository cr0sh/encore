package raknet

import (
	"bytes"
	"github.com/cr0sh/encore/util/binary"
	"io"
	"net"
	"time"
	"unsafe"
)

const (
	// WindowSize is default size of PacketWindow.
	WindowSize = 1024
)

// Options for QueueEncapsulated
const (
	STREAM_OPT_NONE      = iota
	STREAM_OPT_MSGIDX    = 1
	STREAM_OPT_ORDERCHAN = 2 // TODO
)

// PacketWindow is a sized pool for buffering/recovering misordered packet stream.
// NOTE: pool is a fixed-sized array of unsafe.Pointer, but it can be changed to slice in the future.
// PacketWindow.Init must be called once for initialization.
type PacketWindow struct {
	start, end uint64 // valid range: [start,end)
	pool       [WindowSize]unsafe.Pointer
	missing    map[uint64]struct{}
}

// Init initializes PacketWindow.
// Init returns the PacketWindow itself, so we can define
// initialized PacketWindow with new(PacketWindow).Init()
func (window *PacketWindow) Init(trackMissing bool) *PacketWindow {
	window.end = WindowSize
	if trackMissing {
		window.missing = make(map[uint64]struct{})
	}

	return window
}

// Put puts given packet(pointer) into pool if its order has
// a gap between last packet, or does nothing and returns itself.
//
// Put returns ordered list of packet(pointer)s to be processed
// if the packet can be put into the window, otherwise returns nil.
// The returned slice is valid for use only until the next call.
//
// TODO: make time complexity O(1) (currently it is O(WindowSize))
func (window *PacketWindow) Put(order uint64, ptr unsafe.Pointer) []unsafe.Pointer {

	if order == window.start {
		ptrs := []unsafe.Pointer{ptr}
		window.start++
		window.end++

		for window.pool[window.start%WindowSize] != nil {
			ptrs = append(ptrs, window.pool[window.start%WindowSize])

			if window.missing != nil {
				delete(window.missing, window.start)
			}
			window.pool[window.start%WindowSize] = nil
			window.start++
			window.end++
		}
		return ptrs
	}

	if order < window.start || order >= window.end {
		return nil
	}

	window.pool[order%WindowSize] = ptr

	if window.missing != nil {
		delete(window.missing, order)
		order--
		for order >= window.start {
			if _, ok := window.missing[order]; ok ||
				window.pool[order%WindowSize] != nil {
				break
			}
			window.missing[order] = struct{}{}
			order--
		}
	}

	return []unsafe.Pointer{}
}

// GetMissing returns missing sequence numbers for window.
// missing will be reset after call, so callers must process returned list with NACK.
//
// If window is initialized with trackMissing == false, GetMissing returns nil immediately.
func (window *PacketWindow) GetMissing() []uint64 {
	if window.missing == nil {
		return nil
	}
	m := make([]uint64, len(window.missing))
	i := 0
	for order, _ := range window.missing {
		m[i] = order
		i++
	}
	window.missing = make(map[uint64]struct{})
	return m
}

type splitPool struct {
	count   uint32
	packets [][]byte
}

func (sp *splitPool) put(idx uint32, b []byte) []byte {
	l := uint32(len(sp.packets))
	if sp.packets[idx] != nil || idx >= l {
		return nil
	}
	sp.packets[idx] = b
	sp.count++
	if sp.count == l {
		b_ := make([]byte, 0)
		for _, p := range sp.packets {
			b_ = append(b_, p...)
		}
		return b_
	}
	return nil
}

// Session is a set of values for handling single raknet session.
// Its main implementaion purpose is for servers, but also designed for client uses.
// Session.Init must be called once for initialization.
type Session struct {
	// Status:
	//  0: recently initialized, haven't sent/received OpenConnectionRequest1
	//  1: Sent/Received OpenConnectionRequest1
	//  2: Sent/Received OpenConnectionRequest2
	//  3: Handshake succeeded (datapackets available)
	Status int

	// ID is a Client's GUID.
	ID uint64

	StartTime time.Time

	// Address is a remote endpoint address.
	Conn *net.UDPConn
	MTU  int

	sendSplitID      uint16
	sendMessageIndex uint32
	sendQueue        []EncapsulatedPacket

	splitPools map[uint16]splitPool

	// EncapsulatedPacket reliability
	encapsulatedPacketWindow PacketWindow

	// DataPacket reliability
	ackPool, nackPool ACKMap
	recoveryPool      map[uint32][]EncapsulatedPacket
	recvSeq, sendSeq  uint32
	dataPacketWindow  PacketWindow
}

// Init initializes Session.
// Init returns the Session itself, so we can define
// Session with new(Session).Init()
func (sess *Session) Init(conn *net.UDPConn) *Session {
	sess.StartTime = time.Now()
	sess.Conn = conn

	sess.sendQueue = make([]EncapsulatedPacket, 0)

	sess.splitPools = make(map[uint16]splitPool)
	sess.encapsulatedPacketWindow.Init(false)

	sess.ackPool = make(ACKMap)
	sess.nackPool = make(ACKMap)
	sess.recoveryPool = make(map[uint32][]EncapsulatedPacket)
	sess.dataPacketWindow.Init(true)
	return sess
}

// Send copies reader stream to Conn.
func (sess *Session) Send(rd io.Reader) error {
	_, err := io.Copy(sess.Conn, rd)
	return err
}

// FlushSendQueue sends queued EncapsulatedPackets to Conn and resets sendQueue.
func (sess *Session) FlushSendQueue() error {
	err := sess.SendEncapsulatedPackets(sess.sendQueue)
	sess.sendQueue = make([]EncapsulatedPacket, 0)
	return err
}

func splitStream(rd io.Reader, mtu int) ([][]byte, error) {
	bs := make([][]byte, 0)
	for {
		b := make([]byte, mtu-34)
		n, err := io.ReadFull(rd, b)
		if err == io.ErrUnexpectedEOF {
			bs = append(bs, b[:n])
			break
		} else if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

// NOTE: encapsulateBytes has a side-effect that increments
// Session's sendSplitID and sendSplitindex.
func (sess *Session) encapsulateBytes(bs [][]byte, option uint32) []EncapsulatedPacket {
	eps := make([]EncapsulatedPacket, 0, len(bs))
	if len(bs) > 1 {
		for i, b := range bs {
			ep := EncapsulatedPacket{
				IsSplit:    true,
				SplitCount: uint32(len(bs)),
				SplitID:    sess.sendSplitID,
				SplitIndex: uint32(i),

				Payload: b,
			}
			sess.sendSplitID++

			if option&STREAM_OPT_MSGIDX != 0 {
				ep.Reliability = 2
				ep.MessageIndex = sess.sendMessageIndex
				sess.sendMessageIndex++
			}

			eps = append(eps, ep)
		}
	} else {
		ep := EncapsulatedPacket{
			Payload: bs[0],
		}

		if option&STREAM_OPT_ORDERCHAN != 0 {
			ep.Reliability = 2
			ep.MessageIndex = sess.sendMessageIndex
			sess.sendMessageIndex++
		}
		eps = append(eps, ep)
	}

	return eps
}

// QueueEncapsulatedStream queues given reader stream to be sent with EncapsulatedPacket
// until the stream ends.
func (sess *Session) QueueEncapsulatedStream(rd io.Reader, option uint32) error {
	bs, err := splitStream(rd, sess.MTU)
	if err != nil {
		return err
	}
	sess.sendQueue = append(sess.sendQueue, sess.encapsulateBytes(bs, option)...)
	return nil
}

// SendEncapsulatedStream directly sends given stream with EncapsulatedPacket.
func (sess *Session) SendEncapsulatedStream(rd io.Reader, option uint32) error {
	bs, err := splitStream(rd, sess.MTU)
	if err != nil {
		return err
	}

	return sess.SendEncapsulatedPackets(sess.encapsulateBytes(bs, option))
}

// SendEncapsulatedPackets sends given EncapsulatedPackets with
// appropriate number of DataPackets.
func (sess *Session) SendEncapsulatedPackets(eps []EncapsulatedPacket) error {
	length := 0
	start_idx := 0
	for idx := range eps {
		if length+eps[idx].Len()+4 < sess.MTU {
			length += eps[idx].Len()
		} else {
			dp := DataPacket{
				Seq:     binary.LTriad(sess.sendSeq),
				Packets: eps[start_idx:idx],
			}
			buf := new(bytes.Buffer)
			if err := binary.Marshal(dp, buf); err != nil {
				return err
			}

			sess.recoveryPool[sess.sendSeq] = dp.Packets
			sess.sendSeq++

			if err := sess.Send(buf); err != nil {
				return err
			}

			start_idx = idx
		}
	}

	dp := DataPacket{
		Seq:     binary.LTriad(sess.sendSeq),
		Packets: eps[start_idx:],
	}
	buf := new(bytes.Buffer)
	if err := binary.Marshal(dp, buf); err != nil {
		return err
	}

	sess.recoveryPool[sess.sendSeq] = dp.Packets
	sess.sendSeq++

	if err := sess.Send(buf); err != nil {
		return err
	}

	return nil
}

// SendACK packs ackPool into single ACK packet and sends to Conn.
func (sess *Session) SendACK() error {
	if len(sess.ackPool) == 0 {
		return nil
	}

	buf := new(bytes.Buffer)
	buf.WriteByte(0xc0)
	EncodeACK(sess.ackPool, buf)
	return sess.Send(buf)
}

// SendNACK packs nackPool into single NACK packet and sends to Conn.
func (sess *Session) SendNACK() error {
	if len(sess.nackPool) == 0 {
		return nil
	}

	buf := new(bytes.Buffer)
	buf.WriteByte(0xa0)
	EncodeACK(sess.nackPool, buf)
	return sess.Send(buf)
}

// HandleACK handles received ACK packet.
func (sess *Session) HandleACK(keys []uint32) {
	for _, k := range keys {
		delete(sess.recoveryPool, k)
	}
}

// HandleNACK handles received NACK packet.
func (sess *Session) HandleNACK(keys []uint32) error {
	for _, k := range keys {
		if eps, ok := sess.recoveryPool[k]; ok {
			delete(sess.recoveryPool, k)
			if err := sess.SendEncapsulatedPackets(eps); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sess *Session) putSplit(ep EncapsulatedPacket) []byte {
	if !ep.IsSplit {
		panic("putSplit only accepts split packets")
	}
	pool, ok := sess.splitPools[ep.SplitID]
	if !ok {
		pool = splitPool{packets: make([][]byte, ep.SplitCount)}
		sess.splitPools[ep.SplitID] = pool
	}

	return pool.put(ep.SplitIndex, ep.Payload)
}

// HandleDataPacket processes given DataPacket for session and
// returns list of payloads to be processed.
func (sess *Session) HandleDataPacket(dp DataPacket) [][]byte {
	sess.ackPool[uint32(dp.Seq)] = struct{}{}
	ptrs := sess.dataPacketWindow.Put(uint64(dp.Seq), unsafe.Pointer(&dp))
	if len(ptrs) == 0 {
		if ptrs != nil {
			missing := sess.dataPacketWindow.GetMissing()
			for _, m := range missing {
				sess.nackPool[uint32(m)] = struct{}{}
			}
		}
		return nil
	}

	bs := make([][]byte, 0)
	for _, ptr := range ptrs {
		dp := (*DataPacket)(ptr)
		for _, ep := range dp.Packets {
			var ptrs []unsafe.Pointer
			if ep.Reliability >= 2 && ep.Reliability != 5 {
				ptrs = sess.encapsulatedPacketWindow.Put(uint64(ep.MessageIndex), unsafe.Pointer(&ep))

			} else {
				ptrs = []unsafe.Pointer{unsafe.Pointer(&ep)}
			}
			for _, ptr := range ptrs {
				ep := (*EncapsulatedPacket)(ptr)
				if ep.IsSplit {
					if b := sess.putSplit(*ep); b != nil {
						bs = append(bs, b)
					}
				} else {
					bs = append(bs, ep.Payload)
				}
			}
		}
	}

	return bs
}
