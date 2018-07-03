//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package tag

import (
	"github.com/yutopp/go-amf0"
)

type testCase struct {
	Name   string
	Value  interface{}
	Binary []byte
}

var flvTagTestCases = []testCase{
	testCase{
		Name: "AudioData Tag",
		Value: &FlvTag{
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
		},
		Binary: []byte{
			// AudioTag 8
			0x08,
			// DataSize 6
			0x00, 0x00, 0x06,
			// Timestamp 10
			0x00, 0x00, 0x0a,
			// Extended timestamp 0
			0x00,
			// StreamID 0
			0x00, 0x00, 0x00,
			// Audio Data
			0xaf, 0x00, 0x74, 0x65, 0x73, 0x74,
		},
	},
	testCase{
		Name: "VideoData Tag",
		Value: &FlvTag{
			TagType:   TagTypeVideo,
			Timestamp: 10,
			StreamID:  0,
			Data: &VideoData{
				FrameType:       FrameTypeKeyFrame,
				CodecID:         CodecIDAVC,
				AVCPacketType:   AVCPacketTypeSequenceHeader,
				CompositionTime: 0,
				Data:            []byte("test"),
			},
		},
		Binary: []byte{
			// VideoTag 9
			0x09,
			// DataSize 9
			0x00, 0x00, 0x09,
			// Timestamp 10
			0x00, 0x00, 0x0a,
			// Extended timestamp 0
			0x00,
			// StreamID 0
			0x00, 0x00, 0x00,
			// Video Data
			0x17, 0x00, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74,
		},
	},
	testCase{
		Name: "ScriptData Tag",
		Value: &FlvTag{
			TagType:   TagTypeScriptData,
			Timestamp: 10,
			StreamID:  0,
			Data: &ScriptData{
				Objects: map[string]interface{}{
					"test": amf0.ECMAArray{},
				},
			},
		},
		Binary: []byte{
			// ScriptDataTag 18
			0x12,
			// DataSize 15
			0x00, 0x00, 0x0f,
			// Timestamp 10
			0x00, 0x00, 0x0a,
			// Extended timestamp 0
			0x00,
			// StreamID 0
			0x00, 0x00, 0x00,
			// Script Data
			0x02, 0x00, 0x04, 0x74, 0x65, 0x73, 0x74,
			0x08, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x09,
		},
	},
	testCase{
		Name: "Extended timestamp (boundary)",
		Value: &FlvTag{
			TagType:   TagTypeVideo,
			Timestamp: 0xffffff,
			StreamID:  0,
			Data: &VideoData{
				FrameType:       FrameTypeKeyFrame,
				CodecID:         CodecIDAVC,
				AVCPacketType:   AVCPacketTypeSequenceHeader,
				CompositionTime: 0,
				Data:            []byte("test"),
			},
		},
		Binary: []byte{
			// VideoTag 9
			0x09,
			// DataSize 9
			0x00, 0x00, 0x09,
			// Timestamp 16777215
			0xff, 0xff, 0xff,
			// Extended timestamp 0
			0x00,
			// StreamID 0
			0x00, 0x00, 0x00,
			// Video Data
			0x17, 0x00, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74,
		},
	},
	testCase{
		Name: "Extended timestamp",
		Value: &FlvTag{
			TagType:   TagTypeVideo,
			Timestamp: 0xf0ffffff,
			StreamID:  0,
			Data: &VideoData{
				FrameType:       FrameTypeKeyFrame,
				CodecID:         CodecIDAVC,
				AVCPacketType:   AVCPacketTypeSequenceHeader,
				CompositionTime: 0,
				Data:            []byte("test"),
			},
		},
		Binary: []byte{
			// VideoTag 9
			0x09,
			// DataSize 9
			0x00, 0x00, 0x09,
			// Timestamp 16777215
			0xff, 0xff, 0xff,
			// Extended 250
			0xf0,
			// StreamID 0
			0x00, 0x00, 0x00,
			// Video Data
			0x17, 0x00, 0x00, 0x00, 0x00, 0x74, 0x65, 0x73, 0x74,
		},
	},
}

var audioDataTestCases = []testCase{
	testCase{
		Name: "AudioData(AAC, sequence header)",
		Value: &AudioData{
			SoundFormat:   SoundFormatAAC,
			SoundRate:     SoundRate44kHz,
			SoundSize:     SoundSize16Bit,
			SoundType:     SoundTypeStereo,
			AACPacketType: AACPacketTypeSequenceHeader,
			Data:          []byte("test"),
		},
		Binary: []byte{
			// 0xaf: 0b10101111
			//         1010     = SoundFormat 10(AAC)
			//             11   = SoundRate 3(Always 3 when AAC)
			//               1  = SoundSize 1(snd16Bit)
			//                1 = SoundType 1(Always 1 when AAC)
			0xaf,
			// 0x00 = AACPacketType 0(Sequence Header)
			0x00,
			// "test" = AudioSpecificConfig (!DUMMY DATA!)
			0x74, 0x65, 0x73, 0x74,
		},
	},
	testCase{
		Name: "AudioData(AAC, RAW)",
		Value: &AudioData{
			SoundFormat:   SoundFormatAAC,
			SoundRate:     SoundRate44kHz,
			SoundSize:     SoundSize16Bit,
			SoundType:     SoundTypeStereo,
			AACPacketType: AACPacketTypeRaw,
			Data:          []byte("test"),
		},
		Binary: []byte{
			// 0xaf: 0b10101111
			//         1010     = SoundFormat 10(AAC)
			//             11   = SoundRate 3(Always 3 when AAC)
			//               1  = SoundSize 1(snd16Bit)
			//                1 = SoundType 1(Always 1 when AAC)
			0xaf,
			// 0x01 = AACPacketType 1(Raw)
			0x01,
			// "test" = Frame data (!DUMMY DATA!)
			0x74, 0x65, 0x73, 0x74,
		},
	},
	testCase{
		Name: "AudioData(Except AAC)",
		Value: &AudioData{
			SoundFormat: SoundFormatSpeex,
			SoundRate:   SoundRate44kHz,
			SoundSize:   SoundSize16Bit,
			SoundType:   SoundTypeStereo,
			Data:        []byte("test"),
		},
		Binary: []byte{
			// 0xaf: 0b11101111
			//         1011     = SoundFormat 11(Speex)
			//             11   = SoundRate 3(Always 3 when AAC)
			//               1  = SoundSize 1(snd16Bit)
			//                1 = SoundType 1(Always 1 when AAC)
			0xbf,
			// "test" = Sound data (!DUMMY DATA!)
			0x74, 0x65, 0x73, 0x74,
		},
	},
}

var videoDataTestCases = []testCase{
	testCase{
		Name: "VideoData(AVC, sequence header)",
		Value: &VideoData{
			FrameType:       FrameTypeKeyFrame,
			CodecID:         CodecIDAVC,
			AVCPacketType:   AVCPacketTypeSequenceHeader,
			CompositionTime: 0,
			Data:            []byte("test"),
		},
		Binary: []byte{
			// 0x17: 0b00010111
			//         0001     = FrameType 1(Keyframe)
			//             0111 = CodecID 7(AVC)
			0x17,
			// 0x00 = AVCPacketType 0(Sequence Header)
			0x00,
			// 0x00 0x00 0x00 = CompositionTime 0(24bit, BigEndian)
			0x00, 0x00, 0x00,
			// "test" = VideoDecodeConfigurationRecord (!DUMMY DATA!)
			0x74, 0x65, 0x73, 0x74,
		},
	},
	testCase{
		Name: "VideoData(AVC, NALU)",
		Value: &VideoData{
			FrameType:       FrameTypeInterFrame,
			CodecID:         CodecIDAVC,
			AVCPacketType:   AVCPacketTypeNALU,
			CompositionTime: -256,
			Data:            []byte("test"),
		},
		Binary: []byte{
			// 0x27: 0b00020111
			//         0002     = FrameType 2(Interframe)
			//             0111 = CodecID 7(AVC)
			0x27,
			// 0x01 = AVCPacketType 1(NALU)
			0x01,
			// 0xff 0xff 0x00 = CompositionTime -256(24bit, BigEndian)
			0xff, 0xff, 0x00,
			// "test" = NALUs (!DUMMY DATA!)
			0x74, 0x65, 0x73, 0x74,
		},
	},
	testCase{
		Name: "VideoData(Expect AVC)",
		Value: &VideoData{
			FrameType: FrameTypeKeyFrame,
			CodecID:   CodecIDScreenVideo,
			Data:      []byte("test"),
		},
		Binary: []byte{
			// 0x17: 0b00010003
			//         0001     = FrameType 1(KeyInterframe)
			//             0003 = CodecID 3(ScreenVideo)
			0x13,
			// "test" = VideoPacket (!DUMMY DATA!)
			0x74, 0x65, 0x73, 0x74,
		},
	},
}

var scriptDataTestCases = []testCase{
	testCase{
		Name: "ScriptData",
		Value: &ScriptData{
			Objects: map[string]interface{}{
				"test": amf0.ECMAArray{},
			},
		},
		Binary: []byte{
			// AMF0 string marker
			0x02,
			// AMF0 string length: 10
			0x00, 0x04,
			// AMF0 string: test
			0x74, 0x65, 0x73, 0x74,
			// AMF0 ECMA Array marker
			0x08,
			// AMF0 ECMA Array length
			0x00, 0x00, 0x00, 0x00,
			// AMF0 object-property end
			0x00, 0x00, 0x09,
		},
	},
}
