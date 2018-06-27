//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package tag

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/yutopp/go-amf0"
	"io"
	"io/ioutil"
)

func DecodeFlvTag(r io.Reader, flvTag *FlvTag) error {
	ui32 := make([]byte, 4)
	buf := make([]byte, 11)
	if _, err := io.ReadAtLeast(r, buf, 1); err != nil {
		return err
	}

	tagType := TagType(buf[0])

	copy(ui32[1:], buf[1:4]) // 24bits
	dataSize := binary.BigEndian.Uint32(ui32)

	copy(ui32[1:], buf[4:7]) // lower 24bits
	ui32[0] = buf[7]         // upper  8bits
	timestamp := binary.BigEndian.Uint32(ui32)

	copy(ui32[1:], buf[8:11])
	ui32[0] = 0 // clear upper 8bits (not used)
	streamID := binary.BigEndian.Uint32(ui32)

	*flvTag = FlvTag{
		TagType:   tagType,
		Timestamp: timestamp,
		StreamID:  streamID,
	}

	lr := io.LimitReader(r, int64(dataSize))
	switch tagType {
	case TagTypeAudio:
		var v AudioData
		if err := DecodeAudioData(lr, &v); err != nil {
			return err
		}
		flvTag.Data = &v

	case TagTypeVideo:
		var v VideoData
		if err := DecodeVideoData(lr, &v); err != nil {
			return err
		}
		flvTag.Data = &v

	case TagTypeScriptData:
		var v ScriptData
		if err := DecodeScriptData(lr, &v); err != nil {
			return err
		}
		flvTag.Data = &v

	default:
		return fmt.Errorf("Unsupported tag type: %+v", tagType)
	}

	return nil
}

func DecodeAudioData(r io.Reader, audioData *AudioData) error {
	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(r, buf, 1); err != nil {
		return err
	}

	soundFormat := SoundFormat(buf[0] & 0xf0 >> 4) // 0b11110000
	soundRate := SoundRate(buf[0] & 0x0c >> 2)     // 0b00001100
	soundSize := SoundSize(buf[0] & 0x02 >> 1)     // 0b00000010
	soundType := SoundType(buf[0] & 0x01)          // 0b00000001

	*audioData = AudioData{
		SoundFormat: soundFormat,
		SoundRate:   soundRate,
		SoundSize:   soundSize,
		SoundType:   soundType,
	}

	if soundFormat == SoundFormatAAC {
		var aacAudioData AACAudioData
		if err := DecodeAACAudioData(r, &aacAudioData); err != nil {
			return err
		}

		audioData.AACPacketType = aacAudioData.AACPacketType
		audioData.Data = aacAudioData.Data
	} else {
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		audioData.Data = data
	}

	return nil
}

func DecodeAACAudioData(r io.Reader, aacAudioData *AACAudioData) error {
	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(r, buf, 1); err != nil {
		return err
	}

	aacPacketType := AACPacketType(buf[0])
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	*aacAudioData = AACAudioData{
		AACPacketType: aacPacketType,
		Data:          data,
	}

	return nil
}

func DecodeVideoData(r io.Reader, videoData *VideoData) error {
	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(r, buf, 1); err != nil {
		return err
	}

	frameType := FrameType(buf[0] & 0xf0 >> 4) // 0b11110000
	codecID := CodecID(buf[0] & 0x0f)          // 0b00001111

	*videoData = VideoData{
		FrameType: frameType,
		CodecID:   codecID,
	}

	if codecID == CodecIDAVC {
		var avcVideoPacket AVCVideoPacket
		if err := DecodeAVCVideoPacket(r, &avcVideoPacket); err != nil {
			return err
		}
		videoData.AVCPacketType = avcVideoPacket.AVCPacketType
		videoData.CompositionTime = avcVideoPacket.CompositionTime
		videoData.Data = avcVideoPacket.Data
	} else {
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		videoData.Data = data
	}

	return nil
}

func DecodeAVCVideoPacket(r io.Reader, avcVideoPacket *AVCVideoPacket) error {
	buf := make([]byte, 4)
	if _, err := io.ReadAtLeast(r, buf, 4); err != nil {
		return err
	}

	avcPacketType := AVCPacketType(buf[0])
	ctBin := make([]byte, 4)
	copy(ctBin[0:3], buf[1:4])
	compositionTime := int32(binary.BigEndian.Uint32(ctBin)) >> 8 // Signed Interger 24 bits. TODO: check
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	*avcVideoPacket = AVCVideoPacket{
		AVCPacketType:   avcPacketType,
		CompositionTime: compositionTime,
		Data:            data,
	}

	return nil
}

func DecodeScriptData(r io.Reader, data *ScriptData) error {
	r = io.MultiReader(bytes.NewReader([]byte{amf0.MarkerObject}), r) // to treat data as an object

	dec := amf0.NewDecoder(r)
	return dec.Decode(&data.Objects)
}
