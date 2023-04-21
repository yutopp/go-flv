//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package flv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yutopp/go-flv/tag"
)

func TestDecodeSkipBrokenTag(t *testing.T) {
	bin := []byte{
		0x46, 0x4c, 0x56, // FLV
		0x01,                   // version (only 1 is supported)
		0x05,                   // flags (audio + video)
		0x00, 0x00, 0x00, 0x09, // header size (9)
		// 0-pad tag
		0x00, 0x00, 0x00, 0x00, // previous tag size
		// script data (broken)
		0x12,
		0x00, 0x00, 0x05, // 5Bytes
		0x00, 0x00, 0x00, 0x00, // timestamp
		0x00, 0x00, 0x00, // stream id
		0xff, 0xff, 0xff, 0xff, 0xff, // invalid data
		0x00, 0x00, 0x00, 0x10, // script data is 16Bytes
		// video data
		0x09,
		0x00, 0x00, 0x05, // 5Bytes
		0x00, 0x00, 0x00, 0x00, // timestamp
		0x00, 0x00, 0x00, // stream id
		0x01, 0x02, 0x03, 0x04, 0x05,
		0x00, 0x00, 0x00, 0x10, // video data is 16Bytes
	}
	buf := bytes.NewReader(bin)

	dec, err := NewDecoder(buf)
	require.Nil(t, err)

	var flvTag tag.FlvTag

	// script data is broken, thus skipped
	err = dec.Decode(&flvTag)
	require.NotNil(t, err)

	//
	err = dec.Decode(&flvTag)
	require.Nil(t, err)
	require.Equal(t, tag.TagTypeVideo, flvTag.TagType)
}
