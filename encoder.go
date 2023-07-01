//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package flv

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/yutopp/go-flv/tag"
)

type Encoder struct {
	w           io.Writer
	header      *Header
	encodedOnce bool
	cacheBuffer bytes.Buffer
}

func NewEncoder(w io.Writer, flags Flags) (*Encoder, error) {
	header := &Header{
		Version:    1, // only supports 1 currently
		Flags:      flags,
		DataOffset: HeaderLength,
	}
	if err := EncodeFlvHeader(w, header); err != nil {
		return nil, err
	}

	return &Encoder{
		w:      w,
		header: header,
	}, nil
}

func (enc *Encoder) Header() *Header {
	return enc.header
}

func (enc *Encoder) Encode(flvTag *tag.FlvTag) error {
	var previousTagSize uint32
	if !enc.encodedOnce {
		goto tagSize
	}

body:
	enc.cacheBuffer.Reset()
	if err := tag.EncodeFlvTag(&enc.cacheBuffer, flvTag); err != nil {
		return err
	}
	previousTagSize = uint32(enc.cacheBuffer.Len())
	if _, err := io.CopyN(enc.w, &enc.cacheBuffer, int64(previousTagSize)); err != nil {
		return err
	}

tagSize:
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, previousTagSize)
	if _, err := enc.w.Write(buf); err != nil {
		return err
	}

	if !enc.encodedOnce {
		enc.encodedOnce = true
		goto body
	}

	return nil
}

func EncodeFlvHeader(w io.Writer, header *Header) error {
	buf := make([]byte, HeaderLength)

	copy(buf, HeaderSignature)

	buf[3] = header.Version

	buf[4] = 0
	if (header.Flags & FlagsAudio) != 0 {
		buf[4] |= 0x04 // 0b00000100
	}
	if (header.Flags & FlagsVideo) != 0 {
		buf[4] |= 0x01 // 0b00000001
	}

	binary.BigEndian.PutUint32(buf[5:9], HeaderLength)

	_, err := w.Write(buf)
	return err
}
