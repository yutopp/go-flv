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

			bin := make([]byte, len(tc.Binary))
			copy(bin, tc.Binary) // copy ownership

			buf := bytes.NewBuffer(bin)
			var flvTag FlvTag
			err := DecodeFlvTag(buf, &flvTag)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, &flvTag)
		})
	}
}

func TestDecodeAudioDataCommon(t *testing.T) {
	for _, tc := range audioDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			bin := make([]byte, len(tc.Binary))
			copy(bin, tc.Binary) // copy ownership

			buf := bytes.NewBuffer(bin)
			audioData, err := DecodeAudioData(buf)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, audioData)
		})
	}
}

func TestDecodeVideoDataCommon(t *testing.T) {
	for _, tc := range videoDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			bin := make([]byte, len(tc.Binary))
			copy(bin, tc.Binary) // copy ownership

			buf := bytes.NewBuffer(bin)
			videoData, err := DecodeVideoData(buf)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, videoData)
		})
	}
}

func TestDecodeScriptDataCommon(t *testing.T) {
	for _, tc := range scriptDataTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			bin := make([]byte, len(tc.Binary))
			copy(bin, tc.Binary) // copy ownership

			buf := bytes.NewBuffer(bin)
			var scriptData ScriptData
			err := DecodeScriptData(buf, &scriptData)
			assert.Nil(t, err)
			assert.Equal(t, tc.Value, &scriptData)
		})
	}
}
