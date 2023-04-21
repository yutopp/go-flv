//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package tag

import (
	"io"

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
	Data      interface{} // *AudioData | *VideoData | *ScriptData
}

// Close
func (t *FlvTag) Close() {
	// TODO: wrap an error?
	switch data := t.Data.(type) {
	case *AudioData:
		data.Close()
	case *VideoData:
		data.Close()
	}
}

// ========================================
// Audio tags

type SoundFormat uint8

const (
	SoundFormatLinearPCMPlatformEndian SoundFormat = 0
	SoundFormatADPCM                   SoundFormat = 1
	SoundFormatMP3                     SoundFormat = 2
	SoundFormatLinearPCMLittleEndian   SoundFormat = 3
	SoundFormatNellymoser16kHzMono     SoundFormat = 4
	SoundFormatNellymoser8kHzMono      SoundFormat = 5
	SoundFormatNellymoser              SoundFormat = 6
	SoundFormatG711ALawLogarithmicPCM  SoundFormat = 7
	SoundFormatG711muLawLogarithmicPCM SoundFormat = 8
	SoundFormatReserved                SoundFormat = 9
	SoundFormatAAC                     SoundFormat = 10
	SoundFormatSpeex                   SoundFormat = 11
	SoundFormatMP3_8kHz                SoundFormat = 14
	SoundFormatDeviceSpecificSound     SoundFormat = 15
)

type SoundRate uint8

const (
	SoundRate5_5kHz SoundRate = 0
	SoundRate11kHz  SoundRate = 1
	SoundRate22kHz  SoundRate = 2
	SoundRate44kHz  SoundRate = 3
)

type SoundSize uint8

const (
	SoundSize8Bit  SoundSize = 0
	SoundSize16Bit SoundSize = 1
)

type SoundType uint8

const (
	SoundTypeMono   SoundType = 0
	SoundTypeStereo SoundType = 1
)

type AudioData struct {
	SoundFormat   SoundFormat
	SoundRate     SoundRate
	SoundSize     SoundSize
	SoundType     SoundType
	AACPacketType AACPacketType
	Data          io.Reader
}

func (d *AudioData) Read(buf []byte) (int, error) {
	return d.Data.Read(buf)
}

func (d *AudioData) Close() {
	_, _ = io.Copy(io.Discard, d.Data) //  // TODO: wrap an error?
}

type AACPacketType uint8

const (
	AACPacketTypeSequenceHeader AACPacketType = 0
	AACPacketTypeRaw            AACPacketType = 1
)

type AACAudioData struct {
	AACPacketType AACPacketType
	Data          io.Reader
}

// ========================================
// Video Tags

type FrameType uint8

const (
	FrameTypeKeyFrame              FrameType = 1
	FrameTypeInterFrame            FrameType = 2
	FrameTypeDisposableInterFrame  FrameType = 3
	FrameTypeGeneratedKeyFrame     FrameType = 4
	FrameTypeVideoInfoCommandFrame FrameType = 5
)

type CodecID uint8

const (
	CodecIDJPEG                   CodecID = 1
	CodecIDSorensonH263           CodecID = 2
	CodecIDScreenVideo            CodecID = 3
	CodecIDOn2VP6                 CodecID = 4
	CodecIDOn2VP6WithAlphaChannel CodecID = 5
	CodecIDScreenVideoVersion2    CodecID = 6
	CodecIDAVC                    CodecID = 7
)

type VideoData struct {
	FrameType       FrameType
	CodecID         CodecID
	AVCPacketType   AVCPacketType
	CompositionTime int32
	Data            io.Reader
}

func (d *VideoData) Read(buf []byte) (int, error) {
	return d.Data.Read(buf)
}

func (d *VideoData) Close() {
	_, _ = io.Copy(io.Discard, d.Data) //  // TODO: wrap an error?
}

type AVCPacketType uint8

const (
	AVCPacketTypeSequenceHeader AVCPacketType = 0
	AVCPacketTypeNALU           AVCPacketType = 1
	AVCPacketTypeEOS            AVCPacketType = 2
)

type AVCVideoPacket struct {
	AVCPacketType   AVCPacketType
	CompositionTime int32
	Data            io.Reader
}

// ========================================
// Data tags

type ScriptData struct {
	// all values are represented as subset of AMF0
	Objects map[string]amf0.ECMAArray
}
