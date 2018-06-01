//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package tag

import (
	"encoding/binary"
	"fmt"
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

	copy(ui32[1:], buf[1:4]) // 24bit
	dataSize := binary.BigEndian.Uint32(ui32)

	var timestamp uint32
	copy(ui32[1:], buf[4:7]) // 24bit
	timestampBase := binary.BigEndian.Uint32(ui32)
	timestampExt := buf[7] // 8bit

	timestamp = timestampBase // TODO: fix
	_ = timestampExt          // TODO: fix

	copy(ui32[1:], buf[8:11])
	streamID := binary.BigEndian.Uint32(ui32)

	lr := io.LimitReader(r, int64(dataSize))
	var data interface{}
	var err error
	switch tagType {
	case TagTypeAudio:
		data, err = DecodeAudioData(lr)
	case TagTypeVideo:
		data, err = DecodeVideoData(lr)
	case TagTypeScriptData:
		// TODO: implement
		_, err = io.CopyN(ioutil.Discard, lr, int64(dataSize))
	default:
		err = fmt.Errorf("Unsupported tag type: %+v", tagType)
	}
	if err != nil {
		return err
	}

	flvTag.TagType = tagType
	flvTag.Timestamp = timestamp
	flvTag.StreamID = streamID
	flvTag.Data = data

	return nil
}

func DecodeAudioData(r io.Reader) (*AudioData, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(r, buf, 1); err != nil {
		return nil, err
	}

	soundFormat := SoundFormat(buf[0] & 0xf0 >> 4) // 0b11110000
	soundRate := SoundRate(buf[0] & 0x0c >> 2)     // 0b00001100
	soundSize := SoundSize(buf[0] & 0x02 >> 1)     // 0b00000010
	soundType := SoundType(buf[0] & 0x01)          // 0b00000001

	audioData := &AudioData{
		SoundFormat: soundFormat,
		SoundRate:   soundRate,
		SoundSize:   soundSize,
		SoundType:   soundType,
	}
	if soundFormat == SoundFormatAAC {
		aac, err := DecodeAACAudioData(r)
		if err != nil {
			return nil, err
		}

		audioData.AACPacketType = aac.AACPacketType
		audioData.Data = aac.Data
	} else {
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		audioData.Data = data
	}

	return audioData, nil
}

func DecodeAACAudioData(r io.Reader) (*AACAudioData, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(r, buf, 1); err != nil {
		return nil, err
	}

	aacPacketType := AACPacketType(buf[0])
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &AACAudioData{
		AACPacketType: aacPacketType,
		Data:          data,
	}, nil
}

func DecodeVideoData(r io.Reader) (*VideoData, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadAtLeast(r, buf, 1); err != nil {
		return nil, err
	}

	frameType := FrameType(buf[0] & 0xf0 >> 4) // 0b11110000
	codecID := CodecID(buf[0] & 0x0f)          // 0b00001111

	videoData := &VideoData{
		FrameType: frameType,
		CodecID:   codecID,
	}
	if codecID == CodecIDAVC {
		avc, err := DecodeAVCVideoPacket(r)
		if err != nil {
			return nil, err
		}
		videoData.AVCPacketType = avc.AVCPacketType
		videoData.CompositionTime = avc.CompositionTime
		videoData.Data = avc.Data
	} else {
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		videoData.Data = data
	}

	return videoData, nil
}

func DecodeAVCVideoPacket(r io.Reader) (*AVCVideoPacket, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadAtLeast(r, buf, 4); err != nil {
		return nil, err
	}

	avcPacketType := AVCPacketType(buf[0])
	ctBin := make([]byte, 4)
	copy(ctBin[0:3], buf[1:4])
	compositionTime := int32(binary.BigEndian.Uint32(ctBin)) >> 8 // Signed Interger 24 bits. TODO: check
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &AVCVideoPacket{
		AVCPacketType:   avcPacketType,
		CompositionTime: compositionTime,
		Data:            data,
	}, nil
}
