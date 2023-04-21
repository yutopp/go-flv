//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package tag

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yutopp/go-amf0"
)

func TestEncodeFlvTagCommon(t *testing.T) {
	for _, tc := range flvTagTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			switch data := tc.Value.(*FlvTag).Data.(type) {
			case *AudioData:
				data.Data = bytes.NewReader(tc.Payload)
			case *VideoData:
				data.Data = bytes.NewReader(tc.Payload)
			}

			var buf bytes.Buffer
			err := EncodeFlvTag(&buf, tc.Value.(*FlvTag))
			require.Nil(t, err)
			require.Equal(t, tc.Binary, buf.Bytes())
		})
	}
}

func TestEncodeAudioDataCommon(t *testing.T) {
	for _, tc := range audioDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			tc.Value.(*AudioData).Data = bytes.NewReader(tc.Payload) // Restore reader state

			var buf bytes.Buffer
			err := EncodeAudioData(&buf, tc.Value.(*AudioData))
			require.Nil(t, err)
			require.Equal(t, tc.Binary, buf.Bytes())
		})
	}
}

func TestEncodeVideoDataCommon(t *testing.T) {
	for _, tc := range videoDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			tc.Value.(*VideoData).Data = bytes.NewReader(tc.Payload) // Restore reader state

			var buf bytes.Buffer
			err := EncodeVideoData(&buf, tc.Value.(*VideoData))
			require.Nil(t, err)
			require.Equal(t, tc.Binary, buf.Bytes())
		})
	}
}

func BenchmarkEncodeFlvTagCommon(b *testing.B) {
	payload := []byte("test")
	audioData := &AudioData{
		SoundFormat:   SoundFormatAAC,
		SoundRate:     SoundRate44kHz,
		SoundSize:     SoundSize16Bit,
		SoundType:     SoundTypeStereo,
		AACPacketType: AACPacketTypeSequenceHeader,
		Data:          nil,
	}
	tag := &FlvTag{
		TagType:   TagTypeAudio,
		Timestamp: 10,
		StreamID:  0,
		Data:      audioData,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		audioData.Data = bytes.NewReader(payload)
		_ = EncodeFlvTag(io.Discard, tag)
	}
}

func BenchmarkEncodeAudioDataCommon(b *testing.B) {
	payload := []byte("test")
	data := &AudioData{
		SoundFormat:   SoundFormatAAC,
		SoundRate:     SoundRate44kHz,
		SoundSize:     SoundSize16Bit,
		SoundType:     SoundTypeStereo,
		AACPacketType: AACPacketTypeSequenceHeader,
		Data:          nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Data = bytes.NewReader(payload)
		_ = EncodeAudioData(io.Discard, data)
	}
}

func BenchmarkEncodeVideoDataCommon(b *testing.B) {
	payload := []byte("test")
	data := &VideoData{
		FrameType:       FrameTypeKeyFrame,
		CodecID:         CodecIDAVC,
		AVCPacketType:   AVCPacketTypeSequenceHeader,
		CompositionTime: 0,
		Data:            nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Data = bytes.NewReader(payload)
		_ = EncodeVideoData(io.Discard, data)
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
		_ = EncodeScriptData(io.Discard, data)
	}
}
