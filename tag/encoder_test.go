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
	"github.com/yutopp/go-amf0"
	"io/ioutil"
	"testing"
)

func TestEncodeFlvTagCommon(t *testing.T) {
	for _, tc := range flvTagTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := EncodeFlvTag(&buf, tc.Value.(*FlvTag))
			assert.Nil(t, err)
			assert.Equal(t, tc.Binary, buf.Bytes())
		})
	}
}

func TestEncodeAudioDataCommon(t *testing.T) {
	for _, tc := range audioDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := EncodeAudioData(&buf, tc.Value.(*AudioData))
			assert.Nil(t, err)
			assert.Equal(t, tc.Binary, buf.Bytes())
		})
	}
}

func TestEncodeVideoDataCommon(t *testing.T) {
	for _, tc := range videoDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := EncodeVideoData(&buf, tc.Value.(*VideoData))
			assert.Nil(t, err)
			assert.Equal(t, tc.Binary, buf.Bytes())
		})
	}
}

func BenchmarkEncodeFlvTagCommon(b *testing.B) {
	tag := &FlvTag{
		TagType:   TagTypeAudio,
		Timestamp: 10,
		StreamID:  0,
		Data: &AudioData{
			SoundFormat:   SoundFormatAAC,
			SoundRate:     SoundRate44kHz,
			SoundSize:     SoundSize16Bit,
			SoundType:     SoundTypeStereo,
			AACPacketType: AACPacketTypeSequenceHeader,
			Data:          []byte("test"),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeFlvTag(ioutil.Discard, tag)
	}
}

func BenchmarkEncodeAudioDataCommon(b *testing.B) {
	data := &AudioData{
		SoundFormat:   SoundFormatAAC,
		SoundRate:     SoundRate44kHz,
		SoundSize:     SoundSize16Bit,
		SoundType:     SoundTypeStereo,
		AACPacketType: AACPacketTypeSequenceHeader,
		Data:          []byte("test"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeAudioData(ioutil.Discard, data)
	}
}

func BenchmarkEncodeVideoDataCommon(b *testing.B) {
	data := &VideoData{
		FrameType:       FrameTypeKeyFrame,
		CodecID:         CodecIDAVC,
		AVCPacketType:   AVCPacketTypeSequenceHeader,
		CompositionTime: 0,
		Data:            []byte("test"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeVideoData(ioutil.Discard, data)
	}
}

func BenchmarkEncodeScriptDataCommon(b *testing.B) {
	data := &ScriptData{
		Objects: map[string]amf0.ECMAArray{
			"test": nil,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncodeScriptData(ioutil.Discard, data)
	}
}
