package wsjtx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MagicNumber = 0xadbccbda
	Schema2     = 2
	Schema3     = 3
)

type MessageType uint32

const (
	MsgHeartbeat           MessageType = 0
	MsgStatus              MessageType = 1
	MsgDecode              MessageType = 2
	MsgClear               MessageType = 3
	MsgReply               MessageType = 4
	MsgQSOLogged           MessageType = 5
	MsgClose               MessageType = 6
	MsgReplay              MessageType = 7
	MsgHaltTx              MessageType = 8
	MsgFreeText            MessageType = 9
	MsgWSPRDecode          MessageType = 10
	MsgLocation            MessageType = 11
	MsgLoggedADIF          MessageType = 12
	MsgHighlightCallsign   MessageType = 13
	MsgSwitchConfiguration MessageType = 14
	MsgConfigure           MessageType = 15
	MsgAnnotationInfo      MessageType = 16
)

type Message interface {
	Type() MessageType
}

type BaseMessage struct {
	Id string
}

type HeartbeatMessage struct {
	BaseMessage
	MaxSchema uint32
	Version   string
	Revision  string
}

func (m HeartbeatMessage) Type() MessageType { return MsgHeartbeat }

type StatusMessage struct {
	BaseMessage
	DialFrequency      uint64
	Mode               string
	DXCall             string
	Report             string
	TxMode             string
	TxEnabled          bool
	Transmitting       bool
	Decoding           bool
	RxDF               uint32
	TxDF               uint32
	DECall             string
	DEGrid             string
	DXGrid             string
	TxWatchdog         bool
	SubMode            string
	FastMode           bool
	SpecialOpMode      uint8
	FrequencyTolerance uint32
	TRPeriod           uint32
	ConfigName         string
	TxMessage          string
}

func (m StatusMessage) Type() MessageType { return MsgStatus }

type DecodeMessage struct {
	BaseMessage
	New            bool
	Time           uint32 // Milliseconds since midnight
	SNR            int32
	DeltaTime      float64
	DeltaFrequency uint32
	Mode           string
	Message        string
	LowConfidence  bool
	OffAir         bool
}

func (m DecodeMessage) Type() MessageType { return MsgDecode }

// Decoder handles parsing of WSJT-X UDP messages
type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Decode() (Message, error) {
	var magic uint32
	if err := binary.Read(d.r, binary.BigEndian, &magic); err != nil {
		return nil, err
	}
	if magic != MagicNumber {
		return nil, fmt.Errorf("invalid magic number: %x", magic)
	}

	var schema uint32
	if err := binary.Read(d.r, binary.BigEndian, &schema); err != nil {
		return nil, err
	}

	var msgType uint32
	if err := binary.Read(d.r, binary.BigEndian, &msgType); err != nil {
		return nil, err
	}

	id, err := d.readString()
	if err != nil {
		return nil, err
	}

	base := BaseMessage{Id: id}

	switch MessageType(msgType) {
	case MsgHeartbeat:
		return d.decodeHeartbeat(base, schema)
	case MsgStatus:
		return d.decodeStatus(base, schema)
	case MsgDecode:
		return d.decodeDecode(base, schema)
	default:
		// For now, just return base message or error for unknown types
		// Or maybe skip the rest? We don't know the length, so we can't skip easily without reading.
		// But since we are reading from a packet, maybe we can just ignore the rest.
		return nil, fmt.Errorf("unsupported message type: %d", msgType)
	}
}

func (d *Decoder) readString() (string, error) {
	var length uint32
	if err := binary.Read(d.r, binary.BigEndian, &length); err != nil {
		return "", err
	}
	if length == 0xffffffff { // Null string
		return "", nil
	}
	if length == 0 {
		return "", nil
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (d *Decoder) readBool() (bool, error) {
	var b bool
	if err := binary.Read(d.r, binary.BigEndian, &b); err != nil {
		return false, err
	}
	return b, nil
}

func (d *Decoder) readQTime() (uint32, error) {
	var msecs uint32
	if err := binary.Read(d.r, binary.BigEndian, &msecs); err != nil {
		return 0, err
	}
	return msecs, nil
}

func (d *Decoder) decodeHeartbeat(base BaseMessage, schema uint32) (*HeartbeatMessage, error) {
	msg := &HeartbeatMessage{BaseMessage: base}
	var err error
	if err = binary.Read(d.r, binary.BigEndian, &msg.MaxSchema); err != nil {
		return nil, err
	}
	if msg.Version, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.Revision, err = d.readString(); err != nil {
		return nil, err
	}
	return msg, nil
}

func (d *Decoder) decodeStatus(base BaseMessage, schema uint32) (*StatusMessage, error) {
	msg := &StatusMessage{BaseMessage: base}
	var err error
	if err = binary.Read(d.r, binary.BigEndian, &msg.DialFrequency); err != nil {
		return nil, err
	}
	if msg.Mode, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.DXCall, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.Report, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.TxMode, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.TxEnabled, err = d.readBool(); err != nil {
		return nil, err
	}
	if msg.Transmitting, err = d.readBool(); err != nil {
		return nil, err
	}
	if msg.Decoding, err = d.readBool(); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.RxDF); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.TxDF); err != nil {
		return nil, err
	}
	if msg.DECall, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.DEGrid, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.DXGrid, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.TxWatchdog, err = d.readBool(); err != nil {
		return nil, err
	}
	if msg.SubMode, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.FastMode, err = d.readBool(); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.SpecialOpMode); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.FrequencyTolerance); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.TRPeriod); err != nil {
		return nil, err
	}
	if msg.ConfigName, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.TxMessage, err = d.readString(); err != nil {
		return nil, err
	}
	return msg, nil
}

func (d *Decoder) decodeDecode(base BaseMessage, schema uint32) (*DecodeMessage, error) {
	msg := &DecodeMessage{BaseMessage: base}
	var err error
	if msg.New, err = d.readBool(); err != nil {
		return nil, err
	}
	if msg.Time, err = d.readQTime(); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.SNR); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.DeltaTime); err != nil {
		return nil, err
	}
	if err = binary.Read(d.r, binary.BigEndian, &msg.DeltaFrequency); err != nil {
		return nil, err
	}
	if msg.Mode, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.Message, err = d.readString(); err != nil {
		return nil, err
	}
	if msg.LowConfidence, err = d.readBool(); err != nil {
		return nil, err
	}
	if msg.OffAir, err = d.readBool(); err != nil {
		return nil, err
	}
	return msg, nil
}

// ParsePacket parses a raw UDP packet
func ParsePacket(data []byte) (Message, error) {
	return NewDecoder(bytes.NewReader(data)).Decode()
}

// Encoder handles encoding of WSJT-X UDP messages
type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) writeString(s string) error {
	if s == "" {
		// Empty string is usually encoded as 0 length? Or empty?
		// QT spec says: "If the string is null, the length is 0xffffffff"
		// "If the string is empty, the length is 0"
		// Let's assume empty string for now.
		return binary.Write(e.w, binary.BigEndian, uint32(0))
	}
	length := uint32(len(s))
	if err := binary.Write(e.w, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := e.w.Write([]byte(s))
	return err
}

func (e *Encoder) writeBool(b bool) error {
	return binary.Write(e.w, binary.BigEndian, b)
}

func (e *Encoder) writeHeader(msgType MessageType, id string, schema uint32) error {
	if err := binary.Write(e.w, binary.BigEndian, uint32(MagicNumber)); err != nil {
		return err
	}
	if err := binary.Write(e.w, binary.BigEndian, schema); err != nil {
		return err
	}
	if err := binary.Write(e.w, binary.BigEndian, uint32(msgType)); err != nil {
		return err
	}
	return e.writeString(id)
}

func (e *Encoder) Encode(msg Message) error {
	switch m := msg.(type) {
	case *HeartbeatMessage:
		return e.encodeHeartbeat(m)
	case *ReplyMessage:
		return e.encodeReply(m)
	case *HaltTxMessage:
		return e.encodeHaltTx(m)
	default:
		return fmt.Errorf("unsupported message type for encoding: %T", msg)
	}
}

func (e *Encoder) encodeHeartbeat(m *HeartbeatMessage) error {
	if err := e.writeHeader(MsgHeartbeat, m.Id, Schema2); err != nil {
		return err
	}
	if err := binary.Write(e.w, binary.BigEndian, m.MaxSchema); err != nil {
		return err
	}
	if err := e.writeString(m.Version); err != nil {
		return err
	}
	return e.writeString(m.Revision)
}

type ReplyMessage struct {
	BaseMessage
	Time      uint32
	SNR       int32
	DeltaTime float64
	DeltaFreq uint32
	Mode      string
	Message   string
	LowConf   bool
	Modifiers uint8
}

func (m ReplyMessage) Type() MessageType { return MsgReply }

func (e *Encoder) encodeReply(m *ReplyMessage) error {
	if err := e.writeHeader(MsgReply, m.Id, Schema2); err != nil {
		return err
	}
	if err := binary.Write(e.w, binary.BigEndian, m.Time); err != nil {
		return err
	}
	if err := binary.Write(e.w, binary.BigEndian, m.SNR); err != nil {
		return err
	}
	if err := binary.Write(e.w, binary.BigEndian, m.DeltaTime); err != nil {
		return err
	}
	if err := binary.Write(e.w, binary.BigEndian, m.DeltaFreq); err != nil {
		return err
	}
	if err := e.writeString(m.Mode); err != nil {
		return err
	}
	if err := e.writeString(m.Message); err != nil {
		return err
	}
	if err := e.writeBool(m.LowConf); err != nil {
		return err
	}
	return binary.Write(e.w, binary.BigEndian, m.Modifiers)
}

type HaltTxMessage struct {
	BaseMessage
	AutoTxOnly bool
}

func (m HaltTxMessage) Type() MessageType { return MsgHaltTx }

func (e *Encoder) encodeHaltTx(m *HaltTxMessage) error {
	if err := e.writeHeader(MsgHaltTx, m.Id, Schema2); err != nil {
		return err
	}
	return e.writeBool(m.AutoTxOnly)
}
