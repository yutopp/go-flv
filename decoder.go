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
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/yutopp/go-flv/tag"
)

type Decoder struct {
	r           io.Reader
	header      *Header
	decodedOnce bool
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	header, err := DecodeFlvHeader(r)
	if err != nil {
		return nil, err
	}

	if header.DataOffset > HeaderLength {
		offset := header.DataOffset - HeaderLength
		if _, err := io.CopyN(io.Discard, r, int64(offset)); err != nil {
			return nil, err
		}
	}

	return &Decoder{
		r:      r,
		header: header,
	}, nil
}

func (dec *Decoder) Header() *Header {
	return dec.header
}

func (dec *Decoder) Decode(flvTag *tag.FlvTag) error {
	// read previous tag size
	previousTagSize, err := dec.decodeTagSize()
	if err != nil {
		return errors.Wrap(err, "Failed to decode tag size")
	}
	// first size must be 0
	if !dec.decodedOnce {
		if previousTagSize != 0 {
			return fmt.Errorf("initial tag size should be 0: Actual = %d", previousTagSize)
		}

		dec.decodedOnce = true
	}
	// decode tag
	if err := tag.DecodeFlvTag(dec.r, flvTag); err != nil {
		return err
	}
	return nil
}

func (dec *Decoder) decodeTagSize() (uint32, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadAtLeast(dec.r, buf, len(buf)); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(buf), nil
}

func DecodeFlvHeader(r io.Reader) (*Header, error) {
	buf := make([]byte, HeaderLength)
	if _, err := io.ReadAtLeast(r, buf, len(buf)); err != nil {
		return nil, err
	}

	signature := buf[0:3]
	if !bytes.Equal(signature, HeaderSignature) {
		return nil, fmt.Errorf("signature is not matched (FLV): %+v", signature)
	}

	version := buf[3]

	flags := buf[4]
	flagsAudio := (flags & 0b00000100) >> 2
	flagsVideo := (flags & 0b00000001)

	dataOffset := binary.BigEndian.Uint32(buf[5:9])

	header := &Header{
		Version:    version,
		DataOffset: dataOffset,
	}

	if flagsAudio != 0 {
		header.Flags |= FlagsAudio
	}
	if flagsVideo != 0 {
		header.Flags |= FlagsVideo
	}

	return header, nil
}
