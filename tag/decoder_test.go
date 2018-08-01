//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package tag

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeFlvTagCommon(t *testing.T) {
	for _, tc := range flvTagTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewReader(tc.Binary)

			var flvTag FlvTag
			err := DecodeFlvTag(r, &flvTag)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, &flvTag)

			assert.Equal(t, 0, r.Len())
		})
	}
}

func TestDecodeAudioDataCommon(t *testing.T) {
	for _, tc := range audioDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewReader(tc.Binary)

			var audioData AudioData
			err := DecodeAudioData(r, &audioData)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, &audioData)

			assert.Equal(t, 0, r.Len())
		})
	}
}

func TestDecodeVideoDataCommon(t *testing.T) {
	for _, tc := range videoDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewBuffer(tc.Binary)

			var videoData VideoData
			err := DecodeVideoData(r, &videoData)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, &videoData)

			assert.Equal(t, 0, r.Len())
		})
	}
}

func TestDecodeScriptDataCommon(t *testing.T) {
	for _, tc := range scriptDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewBuffer(tc.Binary)

			var scriptData ScriptData
			err := DecodeScriptData(r, &scriptData)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, &scriptData)

			assert.Equal(t, 0, r.Len())
		})
	}
}

func TestDecodeScriptDataPartial(t *testing.T) {
	bin := []byte{0x00} // Invalid data
	r := bytes.NewBuffer(bin)

	var scriptData ScriptData
	err := DecodeScriptData(r, &scriptData)
	assert.NotNil(t, err)
}

func BenchmarkDecodeFlvTagCommon(b *testing.B) {
	bin := []byte{0x08, 0x00, 0x00, 0x06, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x00, 0xaf, 0x00, 0x74, 0x65, 0x73, 0x74}
	buf := bytes.NewReader(bin)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v FlvTag
		_ = DecodeFlvTag(buf, &v)
	}
}

func BenchmarkDecodeAudioDataCommon(b *testing.B) {
	bin := []byte{0xaf, 0x00, 0x74, 0x65, 0x73, 0x74}
	buf := bytes.NewReader(bin)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v AudioData
		_ = DecodeAudioData(buf, &v)
	}
}

func BenchmarkDecodeVideoDataCommon(b *testing.B) {
	bin := []byte{0x17, 0x00, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74}
	buf := bytes.NewReader(bin)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v VideoData
		_ = DecodeVideoData(buf, &v)
	}
}

func BenchmarkDecodeScriptDataCommon(b *testing.B) {
	bin := []byte{0x02, 0x00, 0x04, 0x74, 0x65, 0x73, 0x74, 0x05, 0x00, 0x00, 0x09}
	buf := bytes.NewReader(bin)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v ScriptData
		_ = DecodeScriptData(buf, &v)
	}
}
