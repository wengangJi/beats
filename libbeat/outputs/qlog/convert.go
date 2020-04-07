package qlog

import (
	"encoding/binary"
	"time"
)

// NewQlog func
func NewQlog(server, topic string, level int, idc, thread, relay string, ts time.Time, file string, line int, msg []byte) *Qlog {
	sec := ts.Unix()
	usec := ts.Nanosecond() / 1000
	return &Qlog{
		20,
		1,
		server,
		topic,
		uint32(level),
		idc,
		thread,
		relay,
		uint32(sec),
		uint32(usec),
		file,
		uint32(line),
		msg,
		make([]byte, 0, 1024),
	}
}

// Pack func
func (q *Qlog) Pack() {
	if len(q.bin) != 0 {
		return
	}
	q.appendInt(0)
	q.appendByte(q.version)
	q.appendByte(q.unicode)
	q.appendString(q.server)
	q.appendString(q.topic)
	q.appendUint32(q.level)
	q.appendString(q.idc)
	q.appendString(q.thread)
	q.appendString(q.relay)
	q.appendUint32(q.sec)
	q.appendUint32(q.usec)
	q.appendString(q.file)
	q.appendUint32(q.line)
	q.appendByteSlice(q.msg)
	q.setLength()
}

// GetBin func
func (q *Qlog) GetBin() []byte {
	return q.bin
}

func uint32ToBytes(i uint32) []byte {
	bi := make([]byte, 4)
	binary.BigEndian.PutUint32(bi, i)
	return bi
}

func intToBytes(i int) []byte {
	return uint32ToBytes(uint32(i))
}

func (q *Qlog) setLength() {
	copy(q.bin, intToBytes(len(q.bin)-4))
}

func (q *Qlog) appendUint32(i uint32) {
	q.bin = append(q.bin, uint32ToBytes(i)...)
}

func (q *Qlog) appendInt(i int) {
	q.appendUint32(uint32(i))
}

func (q *Qlog) appendByte(b byte) {
	q.bin = append(q.bin, b)
}

func (q *Qlog) appendRune(r rune) {
	q.appendByte(byte(r))
}

func (q *Qlog) appendString(s string) {
	sl := len(s)
	q.appendInt(sl)
	for _, r := range s {
		q.appendRune(r)
	}
}

func (q *Qlog) appendByteSlice(bs []byte) {
	q.appendInt(len(bs))
	for _, b := range bs {
		q.bin = append(q.bin, b)
	}
}
