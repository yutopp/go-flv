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
	"io"

	"github.com/yutopp/go-amf0"
)

func EncodeFlvTag(w io.Writer, flvTag *FlvTag) error {
	var dataBuf bytes.Buffer // TODO: check performance
	switch flvTag.TagType {
	case TagTypeAudio:
		ad, ok := flvTag.Data.(*AudioData)
		if !ok {
			return fmt.Errorf("Unexpected data is set: not *AudioData")
		}
		if err := EncodeAudioData(&dataBuf, ad); err != nil {
			return err
		}

	case TagTypeVideo:
		vd, ok := flvTag.Data.(*VideoData)
		if !ok {
			return fmt.Errorf("Unexpected data is set: not *VideoData")
		}
		if err := EncodeVideoData(&dataBuf, vd); err != nil {
			return err
		}

	case TagTypeScriptData:
		sd, ok := flvTag.Data.(*ScriptData)
		if !ok {
			return fmt.Errorf("Unexpected data is set: not *ScriptData")
		}
		if err := EncodeScriptData(&dataBuf, sd); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unsupported tag type: %+v", flvTag.TagType)
	}

	ui32 := make([]byte, 4)
	buf := make([]byte, 11)

	buf[0] = byte(flvTag.TagType)

	binary.BigEndian.PutUint32(ui32, uint32(dataBuf.Len()))
	copy(buf[1:4], ui32[1:]) // 24bits

	binary.BigEndian.PutUint32(ui32, flvTag.Timestamp)
	copy(buf[4:7], ui32[1:]) // lower 24bits
	buf[7] = ui32[0]         // upper  8bits

	binary.BigEndian.PutUint32(ui32, flvTag.StreamID)
	copy(buf[8:11], ui32[1:])

	if _, err := w.Write(buf); err != nil {
		return err
	}
	if _, err := io.CopyN(w, &dataBuf, int64(dataBuf.Len())); err != nil {
		return err
	}

	return nil
}

func EncodeAudioData(w io.Writer, audioData *AudioData) error {
	buf := make([]byte, 1)
	buf[0] |= byte(audioData.SoundFormat<<4) & 0xf0 // 0b11110000
	buf[0] |= byte(audioData.SoundRate<<2) & 0x0c   // 0b00001100
	buf[0] |= byte(audioData.SoundSize<<1) & 0x02   // 0b00000010
	buf[0] |= byte(audioData.SoundType) & 0x01      // 0b00000001

	if _, err := w.Write(buf); err != nil {
		return err
	}

	if audioData.SoundFormat == SoundFormatAAC {
		return EncodeAACAudioData(w, &AACAudioData{
			AACPacketType: audioData.AACPacketType,
			Data:          audioData.Data,
		})
	}

	if _, err := io.Copy(w, audioData.Data); err != nil {
		return err
	}

	return nil
}

func EncodeAACAudioData(w io.Writer, aacAudioData *AACAudioData) error {
	buf := make([]byte, 1)
	buf[0] = byte(aacAudioData.AACPacketType)
	if _, err := w.Write(buf); err != nil {
		return err
	}

	if _, err := io.Copy(w, aacAudioData.Data); err != nil {
		return err
	}

	return nil
}

func EncodeVideoData(w io.Writer, videoData *VideoData) error {
	buf := make([]byte, 1)
	buf[0] |= byte(videoData.FrameType<<4) & 0xf0 // 0b11110000
	buf[0] |= byte(videoData.CodecID) & 0x0f      // 0b00001111

	if _, err := w.Write(buf); err != nil {
		return err
	}

	if videoData.CodecID == CodecIDAVC {
		return EncodeAVCVideoPacket(w, &AVCVideoPacket{
			AVCPacketType:   videoData.AVCPacketType,
			CompositionTime: videoData.CompositionTime,
			Data:            videoData.Data,
		})
	}

	if _, err := io.Copy(w, videoData.Data); err != nil {
		return err
	}

	return nil
}

func EncodeAVCVideoPacket(w io.Writer, avcVideoPacket *AVCVideoPacket) error {
	buf := make([]byte, 4)
	buf[0] = byte(avcVideoPacket.AVCPacketType)

	ctBin := make([]byte, 4)
	binary.BigEndian.PutUint32(ctBin, uint32(avcVideoPacket.CompositionTime))
	copy(buf[1:4], ctBin[1:])

	if _, err := w.Write(buf); err != nil {
		return err
	}

	if _, err := io.Copy(w, avcVideoPacket.Data); err != nil {
		return err
	}

	return nil
}

func EncodeScriptData(w io.Writer, data *ScriptData) error {
	enc := amf0.NewEncoder(w)

	for key, value := range data.Objects {
		if err := enc.Encode(key); err != nil {
			return err
		}

		if err := enc.Encode(value); err != nil {
			return err
		}
	}

	return nil
}
