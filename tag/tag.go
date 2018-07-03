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

// ========================================
// FLV tags

type TagType uint8

const (
	TagTypeAudio      TagType = 8
	TagTypeVideo      TagType = 9
	TagTypeScriptData TagType = 18
)

type FlvTag struct {
	TagType
	Timestamp uint32
	StreamID  uint32      // 24bit
	Data      interface{} // *AudioData | *VideoData
}

// ========================================
// Audio tags

type SoundFormat uint8

const (
	SoundFormatLinearPCMPlatformEndian SoundFormat = 0
	SoundFormatADPCM                               = 1
	SoundFormatMP3                                 = 2
	SoundFormatLinearPCMLittleEndian               = 3
	SoundFormatNellymoser16kHzMono                 = 4
	SoundFormatNellymoser8kHzMono                  = 5
	SoundFormatNellymoser                          = 6
	SoundFormatG711ALawLogarithmicPCM              = 7
	SoundFormatG711muLawLogarithmicPCM             = 8
	SoundFormatReserved                            = 9
	SoundFormatAAC                                 = 10
	SoundFormatSpeex                               = 11
	SoundFormatMP3_8kHz                            = 14
	SoundFormatDeviceSpecificSound                 = 15
)

type SoundRate uint8

const (
	SoundRate5_5kHz SoundRate = 0
	SoundRate11kHz            = 1
	SoundRate22kHz            = 2
	SoundRate44kHz            = 3
)

type SoundSize uint8

const (
	SoundSize8Bit  SoundSize = 0
	SoundSize16Bit           = 1
)

type SoundType uint8

const (
	SoundTypeMono   SoundType = 0
	SoundTypeStereo           = 1
)

type AudioData struct {
	SoundFormat   SoundFormat
	SoundRate     SoundRate
	SoundSize     SoundSize
	SoundType     SoundType
	AACPacketType AACPacketType
	Data          []byte
}

type AACPacketType uint8

const (
	AACPacketTypeSequenceHeader AACPacketType = 0
	AACPacketTypeRaw                          = 1
)

type AACAudioData struct {
	AACPacketType AACPacketType
	Data          []byte
}

// ========================================
// Video Tags

type FrameType uint8

const (
	FrameTypeKeyFrame              FrameType = 1
	FrameTypeInterFrame                      = 2
	FrameTypeDisposableInterFrame            = 3
	FrameTypeGeneratedKeyFrame               = 4
	FrameTypeVideoInfoCommandFrame           = 5
)

type CodecID uint8

const (
	CodecIDJPEG                   CodecID = 1
	CodecIDSorensonH263                   = 2
	CodecIDScreenVideo                    = 3
	CodecIDOn2VP6                         = 4
	CodecIDOn2VP6WithAlphaChannel         = 5
	CodecIDScreenVideoVersion2            = 6
	CodecIDAVC                            = 7
)

type VideoData struct {
	FrameType       FrameType
	CodecID         CodecID
	AVCPacketType   AVCPacketType
	CompositionTime int32
	Data            []byte
}

type AVCPacketType uint8

const (
	AVCPacketTypeSequenceHeader AVCPacketType = 0
	AVCPacketTypeNALU                         = 1
	AVCPacketTypeEOS                          = 2
)

type AVCVideoPacket struct {
	AVCPacketType   AVCPacketType
	CompositionTime int32
	Data            []byte
}

// ========================================
// Data tags

type ScriptData struct {
	// all values are represented as subset of AMF0
	Objects map[string]amf0.ECMAArray
}
