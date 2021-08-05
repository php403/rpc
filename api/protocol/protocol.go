package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = uint32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + uint32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)




// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen uint16
		packLen   uint32
	)
	buf := make([]byte,_rawHeaderSize)
	if _, err = io.ReadFull(rr,buf) ; err != nil {
		return
	}
	packLen = binary.BigEndian.Uint32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Uint16(buf[_headerOffset:_verOffset])
	p.Ver = uint32(binary.BigEndian.Uint16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Uint32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Uint32(buf[_seqOffset:])
	bodyLen = int(packLen - uint32(headerLen))

	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - uint32(headerLen)); bodyLen > 0 {
		body := make([]byte,bodyLen)
		_, err = io.ReadFull(rr,body)
		p.Body = body
	} else {
		p.Body = nil
	}
	return
}

// WriteTCP write a proto to TCP writer.
func (p *Proto) WriteTCP(wr *bufio.Writer) (err error) {
	var (
		packLen uint32
	)
	packLen = _rawHeaderSize + uint32(len(p.Body))
	buf := make([]byte,packLen)
	binary.BigEndian.PutUint32(buf[_packOffset:], packLen)
	binary.BigEndian.PutUint16(buf[_headerOffset:], uint16(_rawHeaderSize))
	binary.BigEndian.PutUint16(buf[_verOffset:], uint16(p.Ver))
	binary.BigEndian.PutUint32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutUint32(buf[_seqOffset:], p.Seq)
	binary.BigEndian.PutUint32(buf[_seqOffset:], p.Seq)

	if p.Body != nil {
		_, err = wr.Write(buf)
	}
	return
}


