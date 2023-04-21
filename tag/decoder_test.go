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
)

func assertEqualAudioData(t *testing.T, expected *AudioData, payload []byte, actual *AudioData) {
	require.Equal(t, expected.SoundFormat, actual.SoundFormat)
	require.Equal(t, expected.SoundRate, actual.SoundRate)
	require.Equal(t, expected.SoundSize, actual.SoundSize)
	require.Equal(t, expected.SoundType, actual.SoundType)
	require.Equal(t, expected.AACPacketType, actual.AACPacketType)

	actualPayload, err := io.ReadAll(actual.Data)
	require.Nil(t, err)
	require.Equal(t, payload, actualPayload)
}

func assertEqualVideoData(t *testing.T, expected *VideoData, payload []byte, actual *VideoData) {
	require.Equal(t, expected.FrameType, actual.FrameType)
	require.Equal(t, expected.CodecID, actual.CodecID)
	require.Equal(t, expected.AVCPacketType, actual.AVCPacketType)
	require.Equal(t, expected.CompositionTime, actual.CompositionTime)

	actualPayload, err := io.ReadAll(actual.Data)
	require.Nil(t, err)
	require.Equal(t, payload, actualPayload)
}

func TestDecodeFlvTagCommon(t *testing.T) {
	for _, tc := range flvTagTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewReader(tc.Binary)

			var flvTag FlvTag
			err := DecodeFlvTag(r, &flvTag)
			require.Nil(t, err)

			require.Equal(t, tc.Value.(*FlvTag).TagType, flvTag.TagType)
			require.Equal(t, tc.Value.(*FlvTag).Timestamp, flvTag.Timestamp)
			require.Equal(t, tc.Value.(*FlvTag).StreamID, flvTag.StreamID)

			switch data := flvTag.Data.(type) {
			case *AudioData:
				tcData := tc.Value.(*FlvTag).Data.(*AudioData)
				assertEqualAudioData(t, tcData, tc.Payload, data)
			case *VideoData:
				tcData := tc.Value.(*FlvTag).Data.(*VideoData)
				assertEqualVideoData(t, tcData, tc.Payload, data)
			default:
				require.Equal(t, tc.Value.(*FlvTag).Data, flvTag.Data)
			}

			require.Equal(t, 0, r.Len())
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
			require.Nil(t, err)
			assertEqualAudioData(t, tc.Value.(*AudioData), tc.Payload, &audioData)

			require.Equal(t, 0, r.Len())
		})
	}
}

func TestDecodeEmptyAudio(t *testing.T) {
	r := bytes.NewReader([]byte{})

	var audioData AudioData
	err := DecodeAudioData(r, &audioData)
	require.Equal(t, io.EOF, err)
}

func TestDecodeBrokenAudio(t *testing.T) {
	r := bytes.NewReader([]byte{0xa0}) // AAC requires at least 2Bytes

	var audioData AudioData
	err := DecodeAudioData(r, &audioData)
	require.Equal(t, io.ErrUnexpectedEOF, err)
}

func TestDecodeVideoDataCommon(t *testing.T) {
	for _, tc := range videoDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewBuffer(tc.Binary)

			var videoData VideoData
			err := DecodeVideoData(r, &videoData)
			require.Nil(t, err)
			assertEqualVideoData(t, tc.Value.(*VideoData), tc.Payload, &videoData)

			require.Equal(t, 0, r.Len())
		})
	}
}

func TestDecodeEmptyVideo(t *testing.T) {
	r := bytes.NewReader([]byte{})

	var videoData VideoData
	err := DecodeVideoData(r, &videoData)
	require.Equal(t, io.EOF, err)
}

func TestDecodeBrokenVideo(t *testing.T) {
	r := bytes.NewReader([]byte{0x07}) // AVC requires at least 2Bytes

	var videoData VideoData
	err := DecodeVideoData(r, &videoData)
	require.Equal(t, io.ErrUnexpectedEOF, err)
}

func TestDecodeScriptDataCommon(t *testing.T) {
	for _, tc := range scriptDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewBuffer(tc.Binary)

			var scriptData ScriptData
			err := DecodeScriptData(r, &scriptData)
			require.Nil(t, err)
			require.Equal(t, tc.Value, &scriptData)

			require.Equal(t, 0, r.Len())
		})
	}
}

func TestDecodeEmptyScriptData(t *testing.T) {
	r := bytes.NewReader([]byte{})

	var scriptData ScriptData
	err := DecodeScriptData(r, &scriptData)
	require.Equal(t, nil, err)
}

func TestDecodeBrokenScriptData(t *testing.T) {
	r := bytes.NewReader([]byte{0x01})

	var scriptData ScriptData
	err := DecodeScriptData(r, &scriptData)
	require.EqualError(t, err, "failed to decode key: unexpected EOF")
}

func TestDecodeScriptDataPartial(t *testing.T) {
	bin := []byte{0x00} // Invalid data
	r := bytes.NewBuffer(bin)

	var scriptData ScriptData
	err := DecodeScriptData(r, &scriptData)
	require.NotNil(t, err)
}

func TestDecodeSequentialTag(t *testing.T) {
	bin := []byte{
		// 1
		0x08, 0x00, 0x00, 0x06, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x00, 0xaf, 0x00, 0x74, 0x65, 0x73, 0x74,
		// 2
		0x08, 0x00, 0x00, 0x06, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x00, 0xaf, 0x00, 0x74, 0x65, 0x73, 0x74,
	}

	t.Run("Failed because tag is not closed", func(t *testing.T) {
		buf := bytes.NewReader(bin)

		var v FlvTag
		// 1
		err := DecodeFlvTag(buf, &v)
		require.Nil(t, err)

		// 2
		err = DecodeFlvTag(buf, &v)
		require.NotNil(t, err)
	})

	t.Run("Failed because tag is not closed", func(t *testing.T) {
		buf := bytes.NewReader(bin)

		var v FlvTag
		// 1
		err := DecodeFlvTag(buf, &v)
		require.Nil(t, err)
		v.Close()

		// 2
		err = DecodeFlvTag(buf, &v)
		require.Nil(t, err)
		v.Close()
	})
}

func BenchmarkDecodeFlvTagCommon(b *testing.B) {
	bin := []byte{
		0x08, 0x00, 0x00, 0x06, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x00, 0xaf, 0x00, 0x74, 0x65, 0x73, 0x74,
	}
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
